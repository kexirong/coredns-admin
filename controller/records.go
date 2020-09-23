package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"
)

var conf = config.Get()

//GetRecords gin handle, get all records form etcd database
func GetRecords(c *gin.Context) {

	ex, err := service.GetAllEtcdItems(conf.Etcd.PathPrefix)

	if err != nil {
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	data := []*model.Record{}
	for _, e := range ex {
		r := e.ToRecord()
		if r == nil || "/"+r.Path != conf.Etcd.PathPrefix {
			continue
		}
		data = append(data, r)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": data,
	})
}
func PostRecords(c *gin.Context) {
	var rec model.Record
	err := c.ShouldBindJSON(&rec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if rec.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Empty value for Name field"})
		return
	}
	if rec.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Empty value for Content field"})
		return
	}
	if rec.Type.String() == `""` {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Type field invalid"})
		return
	}
	rec.Content = strings.Trim(rec.Content, " .")
	rec.Path = conf.Etcd.PathPrefix
	etcd, err := rec.ToEtcd()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	j, _ := json.Marshal(etcd)
	fmt.Println(string(j))
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
func PutRecords(c *gin.Context) {

	ex, err := service.GetAllEtcdItems(conf.Etcd.PathPrefix)

	if err != nil {
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	data := []*model.Record{}
	for _, e := range ex {
		r := e.ToRecord()
		if r == nil || "/"+r.Path != conf.Etcd.PathPrefix {
			continue
		}
		data = append(data, r)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": data,
	})
}
func DeleteRecords(c *gin.Context) {

	ex, err := service.GetAllEtcdItems(conf.Etcd.PathPrefix)

	if err != nil {
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	data := []*model.Record{}
	for _, e := range ex {
		r := e.ToRecord()
		if r == nil || "/"+r.Path != conf.Etcd.PathPrefix {
			continue
		}
		data = append(data, r)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": data,
	})
}
