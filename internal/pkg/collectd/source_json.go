package collectd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"collectd.org/api"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type httpJSONCollector struct {
	address    string
	sink       dataSink
	router     *mux.Router
	lastUpdate prometheus.Gauge
}

func (j *httpJSONCollector) Configure(conf Config) error {
	j.router = mux.NewRouter().StrictSlash(true)
	name := "CollectdJSONPost"
	j.router.
		Methods("POST").
		Path(conf.CollectdJSONURLPath).
		Name(name).
		Handler(Logger(http.HandlerFunc(j.handlePost), name))

	j.lastUpdate = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "collectd_update_json_timestamp_seconds",
			Help: "Unix timestamp of the last received collectd metrics JSON push in seconds.",
		},
	)

	return nil
}

func (j *httpJSONCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- j.lastUpdate.Desc()
}

func (j *httpJSONCollector) Run(ctx context.Context) {
	log.Printf("CollectD starting listener: HTTP JSON")
	log.Fatal(http.ListenAndServe(j.address, j.router))
}

func (j *httpJSONCollector) handlePost(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var valueLists []*api.ValueList
	if err := json.Unmarshal(data, &valueLists); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, vl := range valueLists {
		j.sink.Write(r.Context(), vl)
	}
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
