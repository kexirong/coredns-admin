package redis

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"

	"github.com/labstack/echo/v4"
)

func GetDomains(c echo.Context) error {
	var conf = config.Get()
	prefix := conf.Redis.KeyPrefix
	var deep uint8 = 2
	if value, err := strconv.ParseUint(c.QueryParam("deep"), 10, 8); err == nil {
		deep = uint8(value)
	}
	ds, err := domains(prefix, deep)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"reason": err.Error(),
		})

	}

	return c.JSON(http.StatusOK, ds)
}

func domains(prefix string, deep uint8) (tree *model.Domain, err error) {
	tree = model.NewDomainTree()
	keys, err := RedisGetKeys(prefix)

	if err != nil {
		return nil, err
	}

	for _, key := range keys {

		key := strings.TrimPrefix(key, prefix)
		domain := strings.ReplaceAll(key, ":", ".")
		tree.AddSubDomain(domain, deep)
	}
	return tree, nil
}
