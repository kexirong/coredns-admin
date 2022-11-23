package main

import (
	"flag"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/controller"

	"github.com/kexirong/coredns-admin/api"
	"github.com/kexirong/coredns-admin/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var confPath string

func main() {
	flag.StringVar(&confPath, "C", "config.yaml", "config file path")
	flag.Parse()
	config.Set(confPath)
	conf := config.Get()
	err := service.EtcdInitClient(conf)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	// e.File("/", "./dist/index.html")
	e.Static("/js", "./dist/js")
	e.Static("/css", "./dist/css")
	e.Static("/fonts", "./dist/fonts")
	e.File("/favicon.ico", "./dist/favicon.ico")
	e.File("/*", "./dist/index.html")

	e.Use(middleware.CORS())
	e.POST("/login", controller.Login)
	g := e.Group("/api")
	api.RegistRoute(g)

	e.Start(conf.Listen)

}
