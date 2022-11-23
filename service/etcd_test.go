package service

import (
	"testing"

	"github.com/kexirong/coredns-admin/config"
)

func TestGetEtcd(t *testing.T) {
	var err error
	config.LoadDefaultConfig()
	err = EtcdInitClient(config.Get())
	if err != nil {
		t.Fatal(err)
	}

	value, err := EtcdGet("/user/coredns/admin")
	if err != nil {
		t.Fatal(err)
	}
	if value == nil {
		t.Fatal("value  is nil")
	}

}
