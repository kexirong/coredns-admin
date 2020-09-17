package router

import (
	"github.com/kexirong/coredns-admin/controller"
)

var v1 = router.Group("/v1")

func init() {
	v1.GET("/records", controller.GetRecords)
}
