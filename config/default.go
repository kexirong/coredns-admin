package config

import "strings"

const (
	defaultEndpoint   = "http://10.1.1.224:2379"
	defaultPathPrefix = "/coredns"

	defaultEtcdTimeout = 5

	defaultListen = ":8080"
)

func loadDefaultConfig(conf *Config) {

	if conf.Listen == "" {
		conf.Listen = defaultListen
	}

	if len(conf.Etcd.Endpoint) == 0 {
		conf.Etcd.Endpoint = []string{defaultEndpoint}
	}
	if conf.Etcd.PathPrefix == "" {
		conf.Etcd.PathPrefix = defaultPathPrefix
	}
	if !strings.HasSuffix(conf.Etcd.PathPrefix, "/") {
		conf.Etcd.PathPrefix += "/"
	}
	if conf.Etcd.Timeout == 0 {
		conf.Etcd.Timeout = defaultEtcdTimeout
	}

}
