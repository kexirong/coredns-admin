package etcd

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.etcd.io/etcd/api/v3/mvccpb"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"
)

//GetRecords gin handle, get all records form etcd database
func GetRecords(c echo.Context) error {
	param := c.Param("path")
	bp, _ := base64.RawURLEncoding.DecodeString(param)

	var conf = config.Get()
	path := conf.Etcd.PathPrefix

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	path += strings.TrimPrefix(string(bp), "/")

	ex, err := EtcdGetItems(path)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"reason": err.Error(),
		})
	}
	data := []*model.Record{}
	for _, e := range ex {
		r := e.ToRecord(conf.Etcd.PathPrefix)
		if r == nil {
			continue
		}
		data = append(data, r)
	}
	return c.JSON(http.StatusOK, data)
}

func PostRecord(c echo.Context) error {
	var rec model.Record
	var conf = config.Get()
	err := c.Bind(&rec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})
	}
	if rec.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "Empty value for Name field"})

	}
	if rec.Content == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "Empty value for Content field"})

	}
	// if rec.Type.String() == `""` {
	// 	return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Type field invalid"})

	// }

	rec.Content = strings.Trim(rec.Content, " .")

	etcd, err := EtcdFromRecord(&rec, conf.Etcd.PathPrefix)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})

	}
	kvs, err := service.EtcdGetKvs(etcd.Key)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})

	}

	if len(kvs) > 0 {
		bp, abKeys := maxBasicPrefix(kvs)

		if len(abKeys) > 0 {
			abKvs := make(map[string]string)
			for _, k := range abKeys {
				bp = growBasicPrefix(bp)
				value := GetValueFromKVS(kvs, k)
				if value != nil {
					abKvs[k+"/"+bp] = string(value)
				}
			}

			err := service.EtcdPutKvs(abKvs, true)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})

			}
		}
		etcd.Key += "/" + growBasicPrefix(bp)
	}
	err = EtcdPutItem(etcd)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})

	}
	return c.JSON(http.StatusOK, echo.Map{
		"key": etcd.Key,
	})
}

func maxBasicPrefix(kvs []*mvccpb.KeyValue) (bPrefix string, abKeys []string) {
	// bPrefix = "#0"

	for _, kv := range kvs {
		kPart := strings.Split(string(kv.Key), "/")
		if !strings.HasPrefix(kPart[len(kPart)-1], "#") {
			abKeys = append(abKeys, string(kv.Key))
			continue
		}

		switch {
		case len(bPrefix) < len(kPart[len(kPart)-1]):

			bPrefix = kPart[len(kPart)-1]

		case len(bPrefix) == len(kPart[len(kPart)-1]):
			if strings.Compare(bPrefix, kPart[len(kPart)-1]) < 0 {
				bPrefix = kPart[len(kPart)-1]

			}
		}
	}
	return bPrefix, abKeys
}

func growBasicPrefix(bPrefix string) string {
	if !strings.HasPrefix(bPrefix, "#") {
		return "#1"
	}
	if 49 > bPrefix[len(bPrefix)-1] || bPrefix[len(bPrefix)-1] > 56 {
		return bPrefix + "1"
	}
	bbp := []byte(bPrefix)
	bbp[len(bbp)-1]++
	return string(bbp)
}

func DeleteRecord(c echo.Context) error {
	pk := c.Param("key")

	key, _ := base64.RawURLEncoding.DecodeString(pk)

	if string(key) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "Empty value for Key field"})

	}
	err := service.EtcdDelete(string(key))
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
	var rec model.Record

	err = c.Bind(&rec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})

	}
	if rec.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "Empty value for Name field"})

	}
	if rec.Content == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "Empty value for Content field"})

	}

	rec.Content = strings.Trim(rec.Content, " .")
	var conf = config.Get()
	etcd, err := EtcdFromRecord(&rec, conf.Etcd.PathPrefix)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})

	}

	etcd.Key = string(key)
	err = EtcdPutItem(etcd)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": err.Error()})

	}
	return c.NoContent(http.StatusNoContent)
}
