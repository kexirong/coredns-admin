package main

import (
	"fmt"
	"os"
	"strings"

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

	var cmd = "run"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "run":
		err = router.Router.Run(fmt.Sprintf("%s:%s", conf.Host, conf.Port))
		panic(err)
	case "createuser":
		if len(os.Args) != 4 {
			fmt.Println("usage: createuser username password")
		} else {
			path := conf.UserEtcdPath
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}
			path += os.Args[2]
			secret := service.MakeSecret(os.Args[3])
			err := service.EtcdPutKv(path, secret)
			if err != nil {
				panic(err)
			}
			fmt.Printf("user %s created\n", os.Args[2])
		}
	}

}
