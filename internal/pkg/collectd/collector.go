package collectd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"collectd.org/api"
	"collectd.org/network"

	"github.com/fromanirh/virt-collectd-exporter/pkg/nameconv"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var UnknownSecurityLevel = errors.New("Unknown security level")

const Name = "virt_collectd_exporter"

type dataCollector interface {
	Configure(Config) error
	Run(context.Context)
	Describe(ch chan<- *prometheus.Desc)
}

type dataSink interface {
	Write(context.Context, *api.ValueList) error
}

type Collector struct {
	ch      chan api.ValueList
	values  map[string]api.ValueList
	rw      *sync.RWMutex
	srcs    []dataCollector
	address string
	router  *mux.Router
	conv    *nameconv.NameConverter
}

func NewCollector(binaryProtoAddress, httpJSONAddress string) *Collector {
	c := &Collector{
		ch:     make(chan api.ValueList, 0),
		values: make(map[string]api.ValueList),
		rw:     &sync.RWMutex{},
	}
	if binaryProtoAddress != "" {
		log.Printf("CollectD binary protocol endpoint: '%s'", binaryProtoAddress)
		bin := &binaryProtoCollector{
			address: binaryProtoAddress,
			srv: &network.Server{
				Addr:   binaryProtoAddress,
				Writer: c,
			},
		}
		c.srcs = append(c.srcs, bin)
	}
	if httpJSONAddress != "" {
		log.Printf("CollectD HTTP JSON protocol endpoint: '%s'", httpJSONAddress)
		hj := &httpJSONCollector{
			address: httpJSONAddress,
			sink:    c,
		}
		c.srcs = append(c.srcs, hj)
	}
	c.conv, _ = nameconv.NewNameConverter("virt", "virt")
	return c
}

func (c *Collector) Configure(conf Config) error {
	for _, src := range c.srcs {
		if err := src.Configure(conf); err != nil {
			return err
		}
	}

	c.address = conf.MetricsAddress
	c.router = mux.NewRouter().StrictSlash(true)
	name := "metrics"
	c.router.
		Methods("GET").
		Path(conf.MetricsURLPath).
		Name(name).
		Handler(Logger(prometheus.Handler(), name))

	return nil
}

func (c *Collector) Run(ctx context.Context) {
	for _, src := range c.srcs {
		go src.Run(ctx)
	}

	log.Printf("Prometheus endpoint: starting")
	go log.Fatal(http.ListenAndServe(c.address, c.router))

	c.processSamples()
}

func (c Collector) processSamples() {
	ticker := time.NewTicker(time.Minute).C
	for {
		select {
		case vl := <-c.ch:
			c.update(vl)

		case <-ticker:
			c.purge(time.Now())
		}
	}
}

func (c Collector) update(vl api.ValueList) {
	id := vl.Identifier.String()
	c.rw.Lock()
	c.values[id] = vl
	c.rw.Unlock()
}

func (c Collector) purge(now time.Time) {
	c.rw.Lock()
	for id, vl := range c.values {
		expiration := vl.Time.Add(vl.Interval)
		if expiration.Before(now) {
			delete(c.values, id)
		}
	}
	c.rw.Unlock()
}

func (c Collector) Write(_ context.Context, vl *api.ValueList) error {
	c.ch <- *vl
	return nil
}
