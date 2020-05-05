package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/yiningv/nblog/pub/db"
	"github.com/yiningv/nblog/pub/log"
)

var (
	ConfPath string
	Conf     = &Config{}
)

type Config struct {
	App    *AppConfig
	Server *ServerConfig
	ORM    *db.Config
	Zap    *log.ZapConfig
}

type AppConfig struct {
	PageSize int
}

type ServerConfig struct {
	RunMode string
	Address string
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// 初始化配置文件
func Init() (err error) {
	_, err = toml.DecodeFile(ConfPath, Conf)
	return
}
