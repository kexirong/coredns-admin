package router

import "github.com/gin-gonic/gin"

var Router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.Default()
	initRoute()
}
