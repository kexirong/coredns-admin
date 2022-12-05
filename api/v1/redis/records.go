package redis

import (
	"encoding/base64"
	"net/http"
	"strings"

	redisV8 "github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
)

//GetRecords gin handle, get all records form etcd database
func GetRecords(c echo.Context) error {
	param := c.Param("path")
	bp, _ := base64.RawURLEncoding.DecodeString(param)

	var conf = config.Get()
	prefix := conf.Redis.KeyPrefix

	prefix = strings.Join([]string{prefix, string(bp)}, ":")
	prefix = strings.Trim(prefix, ":")
	keys, err := RedisGetKeys(prefix)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"reason": err.Error(),
		})
	}
	rx, err := RedisGetItems(keys)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"reason": err.Error(),
		})
	}
	data := []*model.Record{}
	for _, r := range rx {
		r := r.ToRecords(conf.Redis.KeyPrefix)
		if r == nil {
			continue
		}
		data = append(data, r...)
	}
	return c.JSON(http.StatusOK, data)
}

func PostRecord(c echo.Context) error {
	var rec = new(model.Record)
	var conf = config.Get()
	err := c.Bind(rec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})
	}
	if rec.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "name field empty"})
	}
	if rec.Content == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "content field empty"})
	}
	if strings.HasPrefix(rec.Type.String(), "TYPE") {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Type field invalid"})
	}
	rItem, err := RedisItemFromRecord(rec)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"reason": err.Error()})
	}

	rHasItem, err := RedisGetValue(rec.MustKey(conf.Redis.KeyPrefix), rItem.Field)

	if err != nil && err != redisV8.Nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"reason": err.Error()})
	}

	if err == nil {
		MergeRedisItem(rItem, rHasItem)
	}
	r := rItem.ToRedis(rec.MustKey(conf.Redis.KeyPrefix))
	err = RedisSet(r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"reason": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"key": r.Key,
	})
}

func PostRecordSignature(c echo.Context) error {
	var rec model.Record

	err := c.Bind(&rec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"fingerprint": rec.Signature(),
	})
}

func DeleteRecord(c echo.Context) error {
	pk := c.Param("key")

	key, _ := base64.RawURLEncoding.DecodeString(pk)
	var conf = config.Get()

	fingerprint := c.Request().Header.Get("fingerprint")
	err := RedisDelItem(conf.Redis.KeyPrefix, string(key), fingerprint)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func PutRecord(c echo.Context) error {
	key, err := base64.RawURLEncoding.DecodeString(c.Param("key"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})
	}
	var rec = new(model.Record)

	err = c.Bind(rec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})
	}
	if rec.Key != string(key) {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "invalid key"})
	}
	if rec.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "name field empty"})
	}
	if rec.Content == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "content field empty"})
	}

	var conf = config.Get()

	fingerprint := c.Request().Header.Get("fingerprint")
	err = RedisUpdate(conf.Redis.KeyPrefix, fingerprint, rec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
