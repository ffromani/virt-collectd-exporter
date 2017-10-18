package collectd

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	"collectd.org/api"
	"collectd.org/network"
	"github.com/prometheus/client_golang/prometheus"
)

type binaryProtoCollector struct {
	address    string
	srv        *network.Server
	lastUpdate prometheus.Gauge
}

func (b *binaryProtoCollector) Configure(conf Config) error {
	if conf.CollectdTypesDBPath != "" {
		file, err := os.Open(conf.CollectdTypesDBPath)
		if err != nil {
			return err
		}
		defer file.Close()

		typesDB, err := api.NewTypesDB(file)
		if err != nil {
			return err
		}
		b.srv.TypesDB = typesDB
		log.Printf("CollectD types.db: '%s'", conf.CollectdTypesDBPath)
	}

	if conf.CollectdAuthPath != "" {
		b.srv.PasswordLookup = network.NewAuthFile(conf.CollectdAuthPath)
		log.Printf("CollectD auth file: '%s'", conf.CollectdAuthPath)
	}

	switch strings.ToLower(conf.CollectdSecurityLevel) {
	case "", "none":
		b.srv.SecurityLevel = network.None
	case "sign":
		b.srv.SecurityLevel = network.Sign
	case "encrypt":
		b.srv.SecurityLevel = network.Encrypt
	default:
		return UnknownSecurityLevel
	}
	log.Printf("CollectD security level: '%s'", conf.CollectdSecurityLevel)

	addr, err := net.ResolveUDPAddr("udp", b.address)
	if err != nil {
		return err
	}

	if addr.IP != nil && addr.IP.IsMulticast() {
		b.srv.Conn, err = net.ListenMulticastUDP("udp", nil, addr)
	} else {
		b.srv.Conn, err = net.ListenUDP("udp", addr)
	}
	if err != nil {
		return err
	}
	if conf.CollectdRecvBufferSize >= 0 {
		if err = b.srv.Conn.SetReadBuffer(conf.CollectdRecvBufferSize); err != nil {
			return err
		}
	}

	b.lastUpdate = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "collectd_update_binary_timestamp_seconds",
			Help: "Unix timestamp of the last received collectd metrics binary push in seconds.",
		},
	)

	return nil
}

func (b *binaryProtoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- b.lastUpdate.Desc()
}

func (b *binaryProtoCollector) Run(ctx context.Context) {
	log.Printf("CollectD starting listener: binary protocol")
	log.Fatal(b.srv.ListenAndWrite(ctx))
}
