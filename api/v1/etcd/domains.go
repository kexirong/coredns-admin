package etcd

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"

	"github.com/labstack/echo/v4"
)

func GetDomains(c echo.Context) error {
	var conf = config.Get()
	path := conf.Etcd.PathPrefix
	var deep uint8 = 2
	if value, err := strconv.ParseUint(c.QueryParam("deep"), 10, 8); err == nil {
		deep = uint8(value)
	}
	ds, err := domains(path, deep)

	if err != nil {
		log.Println("err: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"msg": err.Error(),
		})

	}

	return c.JSON(http.StatusOK, echo.Map{
		"msg":  "success",
		"data": ds,
	})
}

func domains(path string, deep uint8) (tree *model.Domain, err error) {
	tree = model.NewDomainTree()
	kvs, err := service.EtcdGetKvs(path)
	if err != nil {
		return nil, err
	}

	for _, kv := range kvs {
		if strings.Contains(string(kv.Key), "/arpa/in-addr") || strings.Contains(string(kv.Key), "/dns/ns") {
			continue
		}
		if err := json.Unmarshal(kv.Value, new(Etcd)); err != nil {
			continue
		}
		key := strings.TrimPrefix(string(kv.Key), path)
		domain := strings.ReplaceAll(key, "/", ".")
		tree.AddSubDomain(domain, deep)

	}
	return tree, nil
}
