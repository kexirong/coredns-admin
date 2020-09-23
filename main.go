package main

import (
	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/router"
	"github.com/kexirong/coredns-admin/service"
)

func main() {

	err := service.InitEtcdClient(config.Get())
	if err != nil {
		panic(err)
	}

	err = router.Router.Run(":8088")
	panic(err)
	// client
	// curl http://127.0.0.1:8080/JSONP?callback=x
}
