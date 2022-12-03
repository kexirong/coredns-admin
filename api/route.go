package api

import (
	etcd_v1 "github.com/kexirong/coredns-admin/api/v1/etcd"
	redis_v1 "github.com/kexirong/coredns-admin/api/v1/redis"
	"github.com/labstack/echo/v4"
)

func RegistRoute(g *echo.Group) {
	v1 := g.Group("/v1")
	etcd_v1.RegistRoute(v1.Group("/etcd"))
	redis_v1.RegistRoute(v1.Group("/redis"))
}
