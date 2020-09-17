package config

const (
	defaultEndpoint    = "http://10.1.1.224:2379"
	defaultPathPrefix  = "/coredns"
	defaultEtcdTimeout = 5
	defaultHost        = ""
	defaultPort        = "8080"
)

func loadDefaultConfig() *Config {
	conf := &Config{
		Host: defaultHost,
		Port: defaultPort,
	}

	conf.Etcd.Endpoint = []string{defaultEndpoint}
	conf.Etcd.PathPrefix = defaultPathPrefix
	conf.Etcd.Timeout = defaultEtcdTimeout

	return conf
}
