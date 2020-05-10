package service

import (
	"flag"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/pub/log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := flag.Set("conf", "../conf/conf.toml"); err != nil {
		panic(err)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.InitLogByConfig(conf.Conf.Zap)
	Init(conf.Conf)
	os.Exit(m.Run())
}
