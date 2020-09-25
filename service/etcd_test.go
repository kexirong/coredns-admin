package service

import (
	"fmt"
	"testing"

	"github.com/kexirong/coredns-admin/config"
)

func TestGetEtcdItems(t *testing.T) {
	var err error
	err = EtcdInitClient(config.Get())
	if err != nil {
		t.Fatal(err)
	}
	ex, err := EtcdGetItems("/coredns")
	if err != nil {
		t.Fatal(err)
	}
	if ex == nil {
		t.Fatal("ex  is nil")
	}
	for _, n := range ex {
		fmt.Printf("%+v\n", n)
	}
	t.Fatal()

}
