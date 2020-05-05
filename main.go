package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/pub/log"
	"github.com/yiningv/nblog/router"
	"net/http"
)

func main() {
	flag.Parse()
	conf.ConfPath = "./conf/conf.toml"
	if err := conf.Init(); err != nil {
		panic(err)
	}
	conf.Conf.Zap.IsDev = conf.Conf.Server.RunMode == gin.DebugMode
	log.InitLogByConfig(conf.Conf.Zap)

	gin.SetMode(conf.Conf.Server.RunMode)
	address := conf.Conf.Server.Address
	log.Info(fmt.Sprintf("Listening and serving HTTP on %s\n", address))
	routes := router.Routes(conf.Conf)
	server := &http.Server{
		Addr:    address,
		Handler: routes,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
