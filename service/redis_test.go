package service

import (
	"testing"

	"github.com/kexirong/coredns-admin/config"
)

const testKey = "coredns:test"

func TestRedisHSet(t *testing.T) {
	var err error
	config.Set("../config.yaml")

	err = RedisInitClient(config.Get())
	if err != nil {
		t.Fatal(err)
	}
	var values = map[string]string{"A": "aaaa", "B": "bbbbb"}
	_, err = RedisHSet(testKey, values)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisHDel(t *testing.T) {
	var err error
	config.Set("../config.yaml")

	err = RedisInitClient(config.Get())
	if err != nil {
		t.Fatal(err)
	}
	fields := []string{"C", "B", "A"}
	for _, field := range fields {
		_, err = RedisHDel(testKey, field)
		if err != nil {
			t.Fatal(err)
		}
	}

}
