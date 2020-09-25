package main

import (
	"fmt"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/router"
	"github.com/kexirong/coredns-admin/service"
)

func main() {
	conf := config.Get()
	err := service.EtcdInitClient(conf)
	if err != nil {
		panic(err)
	}

	err = router.Router.Run(fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	panic(err)
	// client
	// curl http://127.0.0.1:8080/JSONP?callback=x
}
