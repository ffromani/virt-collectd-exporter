package nameconv

import (
	"reflect"
	"testing"

	"collectd.org/api"
	"github.com/prometheus/client_golang/prometheus"
)

var valueLists = []api.ValueList{
	api.ValueList{
		Identifier: api.Identifier{
			Plugin: "cpu",
			Type:   "cpu",
		},
		DSNames: []string{"value"},
		Values:  []api.Value{api.Derive(0)},
	},
	api.ValueList{
		Identifier: api.Identifier{
			Plugin: "df",
			Type:   "df",
		},
		DSNames: []string{"used", "free"},
		Values:  []api.Value{api.Gauge(0), api.Gauge(1)},
	},
	api.ValueList{
		Identifier: api.Identifier{
			Plugin: "df",
			Type:   "df",
		},
		DSNames: []string{"used", "free"},
		Values:  []api.Value{api.Gauge(0), api.Gauge(1)},
	},
	api.ValueList{
		Identifier: api.Identifier{
			Plugin: "cpu",
			Type:   "percent",
		},
		DSNames: []string{"value"},
		Values:  []api.Value{api.Gauge(0)},
	},
	api.ValueList{
		Identifier: api.Identifier{
			Plugin: "interface",
			Type:   "if_octets",
		},
		DSNames: []string{"rx", "tx"},
		Values:  []api.Value{api.Counter(0), api.Counter(1)},
	},
	api.ValueList{
		Identifier: api.Identifier{
			Plugin: "interface",
			Type:   "if_octets",
		},
		DSNames: []string{"rx", "tx"},
		Values:  []api.Value{api.Counter(0), api.Counter(1)},
	},
	api.ValueList{
		Identifier: api.Identifier{
			Plugin: "docker",
			Type:   "cpu.percent",
		},
		DSNames: []string{"value"},
		Values:  []api.Value{api.Gauge(0)},
	},
}

func TestNameConverterName(t *testing.T) {
	cases := []struct {
		vl       api.ValueList
		index    int
		prefix   string
		expected string
	}{
		{api.ValueList{
			Identifier: api.Identifier{
				Plugin: "cpu",
				Type:   "cpu",
			},
			DSNames: []string{"value"},
			Values:  []api.Value{api.Derive(0)},
		}, 0, "test", "test_cpu_total"},
		{api.ValueList{
			Identifier: api.Identifier{
				Plugin: "df",
				Type:   "df",
			},
			DSNames: []string{"used", "free"},
			Values:  []api.Value{api.Gauge(0), api.Gauge(1)},
		}, 0, "test", "test_df_used"},
		{api.ValueList{
			Identifier: api.Identifier{
				Plugin: "df",
				Type:   "df",
			},
			DSNames: []string{"used", "free"},
			Values:  []api.Value{api.Gauge(0), api.Gauge(1)},
		}, 1, "test", "test_df_free"},
		{api.ValueList{
			Identifier: api.Identifier{
				Plugin: "cpu",
				Type:   "percent",
			},
			DSNames: []string{"value"},
			Values:  []api.Value{api.Gauge(0)},
		}, 0, "test", "test_cpu_percent"},
		{api.ValueList{
			Identifier: api.Identifier{
				Plugin: "interface",
				Type:   "if_octets",
			},
			DSNames: []string{"rx", "tx"},
			Values:  []api.Value{api.Counter(0), api.Counter(1)},
		}, 0, "test", "test_interface_if_octets_rx_total"},
		{api.ValueList{
			Identifier: api.Identifier{
				Plugin: "interface",
				Type:   "if_octets",
			},
			DSNames: []string{"rx", "tx"},
			Values:  []api.Value{api.Counter(0), api.Counter(1)},
		}, 1, "test", "test_interface_if_octets_tx_total"},
		{api.ValueList{
			Identifier: api.Identifier{
				Plugin: "docker",
				Type:   "cpu.percent",
			},
			DSNames: []string{"value"},
			Values:  []api.Value{api.Gauge(0)},
		}, 0, "collectd", "collectd_docker_cpu_percent"},
	}

	for _, c := range cases {
		nc, _ := NewNameConverter(c.prefix, c.prefix)
		got, err := nc.Name(c.vl, c.index)
		if err != nil {
			t.Errorf("%s", err)
		}
		if got != c.expected {
			t.Errorf("newName(%v): got %q, expected %q", c.vl, got, c.expected)
		}
	}
}

func TestNameConverterLabels(t *testing.T) {
	cases := []struct {
		vl       api.ValueList
		expected prometheus.Labels
	}{
		{api.ValueList{
			Identifier: api.Identifier{
				Host:           "example.com",
				Plugin:         "cpu",
				PluginInstance: "0",
				Type:           "cpu",
				TypeInstance:   "user",
			},
		}, prometheus.Labels{
			"cpu":      "0",
			"type":     "user",
			"instance": "example.com",
		}},
		{api.ValueList{
			Identifier: api.Identifier{
				Host:         "example.com",
				Plugin:       "df",
				Type:         "df_complex",
				TypeInstance: "used",
			},
		}, prometheus.Labels{
			"df":       "used",
			"instance": "example.com",
		}},
		{api.ValueList{
			Identifier: api.Identifier{
				Host:   "example.com",
				Plugin: "load",
				Type:   "load",
			},
		}, prometheus.Labels{
			"instance": "example.com",
		}},
	}

	nc, _ := NewNameConverter("collectd", "collectd")
	for _, c := range cases {
		got, err := nc.Labels(c.vl)
		if err != nil {
			t.Errorf("%s", err)
		}
		if !reflect.DeepEqual(got, c.expected) {
			t.Errorf("newLabels(%v): got %v, expected %v", c.vl, got, c.expected)
		}
	}
}
