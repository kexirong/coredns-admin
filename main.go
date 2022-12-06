package main

import (
	"flag"
	"net/http"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/controller"

	"github.com/kexirong/coredns-admin/api"
	kcm "github.com/kexirong/coredns-admin/middleware"
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
	err = service.RedisInitClient(conf)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Static("/js", "./dist/js")
	e.Static("/css", "./dist/css")
	e.Static("/fonts", "./dist/fonts")
	e.File("/favicon.ico", "./dist/favicon.ico")
	e.File("/*", "./dist/index.html")

	e.Use(middleware.CORS())
	e.POST("/login", controller.Login)
	g := e.Group("/api", kcm.JWT())
	api.RegistRoute(g)

	e.Start(conf.Listen)

}
func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	if c.Response().Committed {
		return
	}

	he, ok := err.(*echo.HTTPError)
	if !ok {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(he.Code, echo.Map{"reason": he.Message})
	}
	if err != nil {
		c.Logger().Error(err)
	}
}
