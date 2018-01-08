package collectd

import flag "github.com/spf13/pflag"

type Config struct {
	MetricsAddress         string
	CollectdBinaryAddress  string
	CollectdJSONAddress    string
	CollectdRecvBufferSize int
	CollectdAuthPath       string
	CollectdSecurityLevel  string
	CollectdTypesDBPath    string
	MetricsURLPath         string
	CollectdJSONURLPath    string
	DebugLog               bool
}

func ConfigFromCommandLine() Config {
	conf := Config{}
	flag.StringVar(&conf.CollectdBinaryAddress, "collectd-bin-address", ":25826", "Network address on which to accept collectd network protocol pushes.")
	flag.IntVar(&conf.CollectdRecvBufferSize, "collectd-bin-recv-buffer-size", 0, "Size of the buffer of the collectd network protocol receiver")
	flag.StringVar(&conf.CollectdJSONAddress, "collectd-json-address", ":8103", "Network address on which to accept collectd JSON pushes.")
	flag.StringVar(&conf.CollectdJSONURLPath, "collectd-json-url-path", "/collectd", "Collectd write_http URL path")
	flag.StringVar(&conf.CollectdAuthPath, "collectd-auth-path", "", "Path of the collectd auth file")
	flag.StringVar(&conf.CollectdSecurityLevel, "collectd-security-level", "None", "Security level for collectd inbound data (\"None\", \"Sign\" and \"Encrypt\").")
	flag.StringVar(&conf.CollectdTypesDBPath, "collectd-typesdb-path", "/usr/share/collectd/types.db", "Path to collectd types.db (needed for network protocol).")
	flag.StringVar(&conf.MetricsAddress, "metrics-address", ":9103", "Address on which to expose metrics.")
	flag.StringVar(&conf.MetricsURLPath, "metrics-url-path", "/metrics", "Prometheus metrics URL path.")
	flag.BoolVar(&conf.DebugLog, "debug-log", false, "Enable verbose debug log.")
	flag.Parse()
	return conf
}
