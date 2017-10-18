package nameconv

import (
	"fmt"

	"collectd.org/api"
	"github.com/prometheus/client_golang/prometheus"
)

type NameConverter struct {
	source string
	prefix string
}

func NewNameConverter(source, prefix string) (*NameConverter, error) {
	return &NameConverter{
		source: source,
		prefix: prefix + "_",
	}, nil
}

func (n *NameConverter) Describe(vl api.ValueList, index int) (*prometheus.Desc, error) {
	return prometheus.NewDesc(
		n.Name(vl, index),
		fmt.Sprintf("%s: plugin '%s' type: '%s' dstype: '%T' dsname: '%s'",
			n.source, vl.Plugin, vl.Type, vl.Values[index], vl.DSName(index)),
		[]string{},
		n.Labels(vl)), nil
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

func (n *NameConverter) Name(vl api.ValueList, index int) string {
	name := n.prefix
	if vl.Plugin != vl.Type {
		name += vl.Plugin + "_"
	}
	name += vl.Type
	if dsname := vl.DSName(index); dsname != "value" {
		name += "_" + dsname
	}
	switch vl.Values[index].(type) {
	case api.Counter, api.Derive:
		name += "_total"
	}
	return name
}

func (n *NameConverter) Labels(vl api.ValueList) prometheus.Labels {
	labels := prometheus.Labels{}
	labels["instance"] = vl.Host
	if vl.PluginInstance != "" {
		labels[vl.Plugin] = vl.PluginInstance
	}
	if vl.TypeInstance != "" {
		if vl.PluginInstance != "" {
			labels["type"] = vl.TypeInstance
		} else {
			labels[vl.Plugin] = vl.TypeInstance
		}
	}
	return labels
}
