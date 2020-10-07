package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/router"
	"github.com/kexirong/coredns-admin/service"
	"golang.org/x/crypto/ssh/terminal"
)

var adduser bool
var confPath string

func main() {
	flag.BoolVar(&adduser, "adduser", false, "add user")
	flag.StringVar(&confPath, "C", "config.yaml", "config file path")
	flag.Parse()

	config.Set(confPath)
	conf := config.Get()
	err := service.EtcdInitClient(conf)
	if err != nil {
		panic(err)
	}
	if adduser {
		addUser(conf.UserEtcdPath)
	}
	err = router.Router.Run(fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	panic(err)

}

func addUser(prefixPath string) {
	var username, password, confirmPassword string
username:
	fmt.Print("Enter Username：")
	fmt.Scanln(&username)
	if username == "" {
		fmt.Println("Username can not be empty")
		goto username
	}
password:
	fmt.Print("Enter Password: ")
	password = getpassword()
	if len(password) < 6 {
		fmt.Println("\nPassword must be at least 6 characters")
		goto password
	}
	fmt.Print("\nConfirm Password: ")
	confirmPassword = getpassword()
	if confirmPassword != password {
		fmt.Println("\nPassword and confirm password doesn't match")
		goto password
	}

	path := prefixPath + username
	secret := service.MakeSecret(password)
	err := service.EtcdPutKv(path, secret)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nuser %s created\n", username)
	os.Exit(0)
}

func getpassword() string {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(bytePassword))
}
