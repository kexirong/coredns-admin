package etcd

import (
	"github.com/labstack/echo/v4"
)

func RegistRoute(g *echo.Group) {

	g.GET("/records", GetRecords)
	g.GET("/record/:path", GetRecords)
	g.POST("/record", PostRecord)
	g.DELETE("/record/:key", DeleteRecord)
	g.PUT("/record/:key", PutRecord)

	g.GET("/domains", GetDomains)
}
