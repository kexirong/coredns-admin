package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
	etcdcv3 "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

var client *etcdcv3.Client

//InitEtcdClient Must first call this function
func InitEtcdClient(conf *config.Config) (err error) {
	etcd := conf.Etcd
	var tlsConfig *tls.Config
	if len(etcd.TLS) == 3 {
		tlsInfo := transport.TLSInfo{
			CertFile:      etcd.TLS[0],
			KeyFile:       etcd.TLS[1],
			TrustedCAFile: etcd.TLS[2],
		}
		tlsConfig, err = tlsInfo.ClientConfig()
	}

	if err != nil {
		log.Fatal(err)
	}
	etcdCfg := etcdcv3.Config{
		Endpoints: etcd.Endpoint,
		TLS:       tlsConfig,
		Username:  etcd.Username,
		Password:  etcd.Password,
	}

	client, err = etcdcv3.New(etcdCfg)
	if err != nil {
		return err
	}

	return nil
}

const etcdTimeout = 5 * time.Second

//GetAllEtcdItems get all etcd items with path as prifix
func GetAllEtcdItems(path string) (ex []*model.Etcd, err error) {
	if client == nil {
		return nil, errors.New("Client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	r, err := client.Get(ctx, path, etcdcv3.WithPrefix())
	cancel()

	if err != nil {
		return nil, err
	}

	for _, n := range r.Kvs {

		e := new(model.Etcd)
		e.Key = string(n.Key)

		if err := json.Unmarshal(n.Value, e); err != nil {
			return nil, fmt.Errorf("%s: %s", n.Key, err.Error())
		}

		ex = append(ex, e)

	}
	return ex, nil
}

//GetAllEtcdKeys get all etcd keys with path as prifix
func GetAllEtcdKeys(path string) (keys []string, err error) {
	if client == nil {
		return nil, errors.New("Client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	r, err := client.Get(ctx, path, etcdcv3.WithPrefix(), etcdcv3.WithKeysOnly())
	cancel()

	if err != nil {
		return nil, err
	}

	for _, n := range r.Kvs {
		keys = append(keys, string(n.Key))
	}
	return
}
