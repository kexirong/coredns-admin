package router

import (
	"github.com/kexirong/coredns-admin/controller"
	"github.com/kexirong/coredns-admin/middleware"
)

func init() {
	Router.Use(middleware.Cors())

	var v1 = Router.Group("/api/v1")
	v1.GET("/records", controller.GetRecords)
	v1.POST("/records", controller.PostRecords)
	v1.DELETE("/records/:key", controller.DeleteRecords)
}
