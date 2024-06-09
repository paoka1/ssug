package base

import (
	"flag"
	"github.com/sirupsen/logrus"
	"ssug/internal/utils"
	"strconv"
)

var Debug bool

type Para struct {
	Host     string
	Port     int
	Key      string
	InitLen  int
	TTL      int64
	HostPort string
}

// ParsePara 解析命令行参数
func ParsePara() Para {
	host := flag.String("host", "127.0.0.1", "监听地址")
	port := flag.Int("port", 8000, "监听端口")
	initLen := flag.Int("len", 3, "短链最小长度")
	key := flag.String("key", "key123456", "管理短链的密钥")
	ttl := flag.Int64("ttl", 60*60*24, "短链存活时长（秒）")
	d := flag.Bool("debug", false, "debug模式")

	flag.Parse()

	if *d {
		utils.Logger.SetLevel(logrus.DebugLevel)
		Debug = true
		utils.Logger.Debug("已开启debug模式")
	}
	return Para{
		Host:     *host,
		Port:     *port,
		Key:      *key,
		InitLen:  *initLen,
		TTL:      *ttl,
		HostPort: *host + ":" + strconv.Itoa(*port),
	}
}
