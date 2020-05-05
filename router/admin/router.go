package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/service"
	"net/http"
)

const (
	OK  = 1
	ERR = 0
)

var (
	srv *service.Service
)

func Init(c *conf.Config) {
	srv = service.New(c)
}

func ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": OK})
}

func okData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": OK, "data": data})
}

func fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": ERR, "msg": msg})
}
