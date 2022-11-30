package service

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/kexirong/coredns-admin/config"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdcv3 "go.etcd.io/etcd/client/v3"

	"go.etcd.io/etcd/pkg/transport"
)

var etcdClient *etcdcv3.Client

//EtcdInitClient Must first call this function
func EtcdInitClient(conf *config.Config) (err error) {
	etcd := conf.Etcd
	var tlsConfig *tls.Config
	if len(etcd.TLS) == 3 {
		tlsInfo := transport.TLSInfo{
			CertFile:      etcd.TLS[0],
			KeyFile:       etcd.TLS[1],
			TrustedCAFile: etcd.TLS[2],
		}
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	etcdCfg := etcdcv3.Config{
		Endpoints: etcd.Endpoint,
		TLS:       tlsConfig,
		Username:  etcd.Username,
		Password:  etcd.Password,
	}

	etcdClient, err = etcdcv3.New(etcdCfg)
	if err != nil {
		return err
	}

	if etcd.Timeout > 0 {
		etcdTimeout = time.Duration(etcd.Timeout) * time.Second
	}

	return nil
}

var etcdTimeout = 5 * time.Second

//EtcdGetKvs get  etcd keys with path as prifix
func EtcdGetKvs(path string) ([]*mvccpb.KeyValue, error) {
	if etcdClient == nil {
		return nil, errors.New("etcd Client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	defer cancel()
	r, err := etcdClient.Get(ctx, path, etcdcv3.WithPrefix())

	if err != nil {
		return nil, err
	}

	return r.Kvs, nil
}

func EtcdPutKv(key, value string) (err error) {
	if etcdClient == nil {
		return errors.New("etcd Client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	defer cancel()
	_, err = etcdClient.Put(ctx, key, value)

	return err
}

func EtcdPutKvs(kvs map[string]string, del bool) (err error) {
	if etcdClient == nil {
		return errors.New("etcd Client is not initialized")
	}
	if len(kvs) == 0 {
		return
	}
	kvc := etcdcv3.NewKV(etcdClient)
	var ops []etcdcv3.Op
	for k, v := range kvs {
		if del {
			idx := strings.LastIndex(k, "/")
			ops = append(ops, etcdcv3.OpDelete(k[:idx]))
		}
		ops = append(ops, etcdcv3.OpPut(k, v))
	}

	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	defer cancel()
	_, err = kvc.Txn(ctx).If().
		Then(ops...).
		Commit()

	return err
}

func EtcdDelete(key string) (err error) {
	if etcdClient == nil {
		return errors.New("etcd Client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	defer cancel()
	_, err = etcdClient.Delete(ctx, key)

	return err
}

func EtcdGet(key string) (value []byte, err error) {
	if etcdClient == nil {
		return nil, errors.New("etcd Client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	r, err := etcdClient.Get(ctx, key)
	cancel()

	if r.Count != 1 || err != nil {
		return nil, err
	}

	return r.Kvs[0].Value, err
}
