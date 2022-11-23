package etcd

import (
	"testing"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"
)

func TestEtcdGetItems(t *testing.T) {
	var err error
	err = service.EtcdInitClient(config.Get())
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

}

func TestDomains(t *testing.T) {
	config.LoadDefaultConfig()
	err := service.EtcdInitClient(config.Get())
	if err != nil {
		t.Error(err)
	}
	tree, err := domains("/coredns/", 2)
	if err != nil {
		t.Error(err)
	}
	var f func(d *model.Domain)
	f = func(d *model.Domain) {

		for _, v := range d.SubDomain {
			if len(v.SubDomain) > 0 {
				f(v)
			}
		}
	}
	f(tree)

}
