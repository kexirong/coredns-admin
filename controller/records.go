package controller

import (
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	data := []model.Record{}
	for _, e := range ex {
		record := model.Record{}
		switch e.HostType() {
		case model.TypeA:

		}
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": ex,
	})
}
