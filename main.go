package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/pub/log"
	"github.com/yiningv/nblog/router"
	"github.com/yiningv/nblog/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	conf.Conf.Zap.IsDev = conf.Conf.Server.RunMode == gin.DebugMode
	log.InitLogByConfig(conf.Conf.Zap)

	gin.SetMode(conf.Conf.Server.RunMode)
	address := conf.Conf.Server.Address
	log.Info(fmt.Sprintf("Listening and serving HTTP on %s\n", address))
	srv := service.New(conf.Conf)
	routes := router.Routes(srv)
	server := &http.Server{
		Addr:    address,
		Handler: routes,
	}
	handleSignal(server, srv)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server, srv *service.Service) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		log.Info(fmt.Sprintf("got signal [%s], exiting pipe now", s))
		if err := server.Close(); nil != err {
			log.Info(fmt.Sprintf("server close failed: %v", err))
		}

		if err := srv.Close(); nil != err {
			log.Info(fmt.Sprintf("service close failed: %v", err))
		}

		log.Info("nblog exited")
		os.Exit(0)
	}()
}
