package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"
)

var conf = config.Get()

//GetRecords gin handle, get all records form etcd database
func GetRecords(c *gin.Context) {

	ex, err := service.GetAllEtcdItems(conf.Etcd.PathPrefix)
	log.Println("len(ex): ", len(ex))
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
		if r == nil {
			continue
		}
		data = append(data, r)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": data,
	})
}
