package nameconv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"collectd.org/api"
	"github.com/prometheus/client_golang/prometheus"
)

type LabelItem struct {
	Label string `json:"label"`
	Ident string `json:"ident"`
}

type ConfMap struct {
	Source string                 `json:"source"`
	Prefix string                 `json:"prefix"`
	Name   string                 `json:"name"`
	Labels map[string][]LabelItem `json:"labels"`
}

type VLDesc struct {
	Host           string
	Plugin         string
	PluginInstance string
	Type           string
	TypeInstance   string
	DSName         string
	IsTotal        bool
}

func process(vl api.ValueList, index int) VLDesc {
	vldesc := VLDesc{
		Host:           vl.Host,
		Plugin:         vl.Plugin,
		PluginInstance: vl.PluginInstance,
		Type:           vl.Type,
		TypeInstance:   vl.TypeInstance,
	}
	if index != -1 {
		vldesc.DSName = vl.DSName(index)
		switch vl.Values[index].(type) {
		case api.Counter, api.Derive:
			vldesc.IsTotal = true
		}
	}
	return vldesc
}

type NameConverter struct {
	source string
	prefix string
	conf   *ConfMap
}

func NewNameConverter(source, prefix string) (*NameConverter, error) {
	return &NameConverter{
		source: source,
		prefix: prefix + "_",
	}, nil
}

func NewNameConverterWithJSON(data []byte) (*NameConverter, error) {
	var c ConfMap
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &NameConverter{
		source: c.Source,
		prefix: c.Prefix + "_",
		conf:   &c,
	}, nil

}

func (n *NameConverter) Describe(vl api.ValueList, index int) (*prometheus.Desc, error) {
	vldesc := process(vl, index)

	name, err := n.convertName(vldesc)
	if err != nil {
		return nil, err
	}
	labels, err := n.convertLabels(vldesc)
	if err != nil {
		return nil, err
	}

	return prometheus.NewDesc(
		name,
		fmt.Sprintf("%s: plugin '%s' type: '%s' dstype: '%T' dsname: '%s'",
			n.source, vl.Plugin, vl.Type, vl.Values[index], vl.DSName(index)),
		[]string{},
		labels), nil
}

func (n *NameConverter) Convert(vl api.ValueList, index int) (prometheus.Metric, error) {
	var prometheusType prometheus.ValueType
	var prometheusValue float64

	switch v := vl.Values[index].(type) {
	case api.Counter:
		prometheusType = prometheus.CounterValue
		prometheusValue = float64(v)
	case api.Derive:
		prometheusType = prometheus.CounterValue
		prometheusValue = float64(v)
	case api.Gauge:
		prometheusType = prometheus.GaugeValue
		prometheusValue = float64(v)
	default:
		return nil, fmt.Errorf("Unknown value type: %T", v)
	}

	desc, err := n.Describe(vl, index)
	if err != nil {
		return nil, err
	}
	return prometheus.NewConstMetric(desc, prometheusType, prometheusValue)
}

func (n *NameConverter) Name(vl api.ValueList, index int) (string, error) {
	return n.convertName(process(vl, index))
}

func (n *NameConverter) convertName(vldesc VLDesc) (string, error) {
	var name string
	var err error
	if n.conf != nil {
		name, err = n.userName(vldesc)
	} else {
		name, err = n.builtinName(vldesc)
	}
	if err != nil {
		return "", err
	}
	return strings.Replace(name, ".", "_", -1), nil
}

func (n *NameConverter) userName(vldesc VLDesc) (string, error) {
	buf := new(bytes.Buffer)
	t := template.Must(template.New("name").Parse(n.conf.Name))
	err := t.Execute(buf, vldesc)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (n *NameConverter) builtinName(vldesc VLDesc) (string, error) {
	name := ""
	if vldesc.Plugin != vldesc.Type {
		name += vldesc.Plugin + "_"
	}

	if n.prefix != "" {
		// edge case: duplicated prefix
		if !strings.HasPrefix(name, n.prefix) {
			name = n.prefix + name
		}
	}

	name += vldesc.Type
	if dsname := vldesc.DSName; dsname != "value" {
		name += "_" + dsname
	}

	if vldesc.IsTotal {
		// edge case: duplicated total
		if !strings.HasSuffix(name, "total") {
			name += "_total"
		}
	}
	return name, nil
}

func (n *NameConverter) Labels(vl api.ValueList) (prometheus.Labels, error) {
	return n.convertLabels(process(vl, -1))
}

func (n *NameConverter) convertLabels(vldesc VLDesc) (prometheus.Labels, error) {
	if n.conf != nil {
		return n.userLabels(vldesc)
	}
	return n.builtinLabels(vldesc)
}

func (n *NameConverter) builtinLabels(vldesc VLDesc) (prometheus.Labels, error) {
	labels := prometheus.Labels{}
	labels["instance"] = vldesc.Host
	if vldesc.PluginInstance != "" {
		labels[vldesc.Plugin] = vldesc.PluginInstance
	}
	if vldesc.TypeInstance != "" {
		if vldesc.PluginInstance != "" {
			labels["type"] = vldesc.TypeInstance
		} else {
			labels[vldesc.Plugin] = vldesc.TypeInstance
		}
	}
	return labels, nil
}

func (n *NameConverter) userLabels(vldesc VLDesc) (prometheus.Labels, error) {
	labels := prometheus.Labels{}
	items, ok := n.conf.Labels["*"]
	if !ok {
		return labels, errors.New("No defaults")
	}
	v := reflect.ValueOf(vldesc)
	for _, item := range items {
		labels[resolve(v, item.Label)] = resolve(v, item.Ident)
	}
	return labels, nil
}

func resolve(v reflect.Value, name string) string {
	if strings.HasPrefix(name, "$") {
		name = strings.TrimPrefix(name, "$")
		return v.FieldByName(name).String()
	}
	return name
}
