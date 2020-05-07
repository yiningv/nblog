package notion

import (
	"flag"
	"github.com/smartystreets/goconvey/convey"
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
	os.Exit(m.Run())
}

func TestGetSiteConfig(t *testing.T) {
	convey.Convey("GetSiteConfig", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			config, err := GetSiteConfig()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			ctx.Convey("Then err should not be empty.", func(ctx convey.C) {
				ctx.So(config, convey.ShouldNotBeEmpty)
			})
		})
	})
}

func TestGetSourceConfig(t *testing.T) {
	convey.Convey("GetSourceConfig", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			config, err := GetSourceConfig()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			ctx.Convey("Then err should not be empty.", func(ctx convey.C) {
				ctx.So(config, convey.ShouldNotBeEmpty)
			})
		})
	})
}
