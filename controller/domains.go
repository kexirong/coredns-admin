package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/service"
)

func GetDomains(c *gin.Context) {
	var conf = config.Get()
	path := conf.Etcd.PathPrefix
	var deep uint8 = 2
	if value, err := strconv.ParseUint(c.Query("deep"), 10, 8); err == nil {
		deep = uint8(value)
	}
	ds, err := service.Domains(path, deep)

	if err != nil {
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": ds,
	})
}
