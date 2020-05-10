package main

import (
	"context"
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
	"time"
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
	service.Init(conf.Conf)
	go func() {
		service.LoadCache()
		service.SyncData()
	}()
	routes := router.Routes()
	server := &http.Server{
		Addr:    address,
		Handler: routes,
	}
	handleSignal(server)
	if err := server.ListenAndServe(); err != nil {
		log.Error(fmt.Sprintf("listen and serve failed: %v", err))
	}
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		s := <-c
		log.Info(fmt.Sprintf("get a signal %s", s.String()))
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); nil != err {
				log.Info(fmt.Sprintf("server shutdown failed: %v", err))
			}
			if err := service.Close(); nil != err {
				log.Info(fmt.Sprintf("service close failed: %v", err))
			}
			log.Info("nblog exited")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}()
}
