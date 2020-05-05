package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/router/admin"
)

func Routes(c *conf.Config) *gin.Engine {
	admin.Init(c)

	r := gin.New()

	return r
}
