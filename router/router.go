package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/pub/e"
	"github.com/yiningv/nblog/service"
	"net/http"
)

var (
	srv *service.Service
)

func Routes(s *service.Service) *gin.Engine {
	srv = s
	r := gin.New()

	return r
}

func ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": e.OK})
}

func okData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": e.OK, "data": data})
}

func fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": e.ERR, "msg": msg})
}
