package service

import (
	"testing"

	"github.com/kexirong/coredns-admin/config"
)

func TestGetAllEtcdItems(t *testing.T) {
	var err error
	err = InitEtcdClient(*config.Get())
	if err != nil {
		t.Fatal(err)
	}
	ex, err := GetAllEtcdItems("/coredns")
	if err != nil {
		t.Fatal(err)
	}
	if ex == nil {
		t.Fatal("ex  is nil")
	}

}
