package collectd

import (
	"collectd.org/api"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, src := range c.srcs {
		src.Describe(ch)
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.rw.RLock()
	values := make([]api.ValueList, 0, len(c.values))
	for _, vl := range c.values {
		values = append(values, vl)
	}
	c.rw.RUnlock()

	for _, vl := range values {
		for i := range vl.Values {
			m, err := c.conv.Convert(vl, i)
			if err != nil {
				log.Printf("%s", err) // TODO
				continue
			}

			ch <- m
		}
	}

}
