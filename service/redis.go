package service

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"time"

	redisV8 "github.com/go-redis/redis/v8"
	"github.com/kexirong/coredns-admin/config"
	"go.etcd.io/etcd/pkg/transport"
)

var redisClient redisV8.UniversalClient

func RedisInitClient(conf *config.Config) (err error) {
	var tlsConfig *tls.Config
	redis := conf.Redis
	if len(redis.TLS) == 3 {
		tlsInfo := transport.TLSInfo{
			CertFile:      redis.TLS[0],
			KeyFile:       redis.TLS[1],
			TrustedCAFile: redis.TLS[2],
		}
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	redisClient = redisV8.NewUniversalClient(&redisV8.UniversalOptions{
		Addrs:       redis.Addresses,
		Username:    redis.Username,
		Password:    redis.Password,
		DialTimeout: time.Second * time.Duration(redis.ConnectTimeout),
		ReadTimeout: time.Duration(redis.ReadTimeout) * time.Second,
		TLSConfig:   tlsConfig,
	})

	return nil
}

func RedisKeys(pattern string) ([]string, error) {
	if redisClient == nil {
		return nil, errors.New("etcd Client is not initialized")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return redisClient.Keys(ctx, pattern).Result()
}

func RedisHGetAll(pattern string) (map[string]string, error) {
	if redisClient == nil {
		return nil, errors.New("etcd Client is not initialized")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return redisClient.HGetAll(ctx, pattern).Result()
}
