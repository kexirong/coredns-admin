package config

import (
	"testing"
)

func TestGet(t *testing.T) {

}

func TestLoadDefaultConfig(t *testing.T) {
	loadDefaultConfig(conf)

	if conf.Listen != defaultListen {
		t.Fatalf("conf.Listen != defaultListen")
	}
	etcd := conf.Etcd
	if etcd.PathPrefix != defaultPathPrefix {
		t.Fatalf("etcd.PathPrefix != defaultPathPrefix")
	}
	if etcd.Endpoint[0] != defaultEndpoint {
		t.Fatalf("etcd.Endpoint[0] != defaultEndpoint")
	}
	if etcd.Timeout != defaultEtcdTimeout {
		t.Fatalf("etcd.Timeout != defaultEtcdTimeout ")
	}
}

func TestLoadConfig(t *testing.T) {
	loadConfig("../config.yaml")

}
