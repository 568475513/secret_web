package conf

import (
	"fmt"
	"log"
	"os"
	// "strconv"
	// "strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/joho/godotenv"
)

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type ZapSetting struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

// 服务配置
var ServerConf = &Server{}

// Zap配置
var ZapConf = &ZapSetting{}

// Init 初始化配置项
func Init(env string) {
	// 获取环境变量
	env = getEnvMode(env)

	// 读取对应环境变量
	err := godotenv.Load(".env." + env)
	fmt.Printf(">开始初始化配置[%s]...\n", env)
	if err != nil {
		log.Fatalf("Load config env err: %v", err)
	}

	cfg := ini.Empty()
	// 特殊配置须转化为结构体
	serverConfMapTo(cfg)
	zapConfMapTo(cfg)
	fmt.Println(">>>初始化配置完成")
}

// 获取系统环境GOENVMODE, 根据GOENVMODE加载不同的配置文件,GOENVMODE有以下几个值
// local: 开发人员本地环境配置.
// develop: 公司开发机环境配置.
// test: 测试环境配置.
// production: 生产环境配置.
func getEnvMode(env string) string {
	// 没有环境变量默认本地配置
	if env == "" {
		env = os.Getenv("GOENVMODE")
		if env == "" {
			env = "local"
		}
	}
	// 命令调整参数【已转cmd结构】
	// for _, v := range os.Args {
	// 	if strings.Contains(v, "env=") {
	// 		env = v[6:]
	// 	}
	// 	if strings.Contains(v, "port=") {
	// 		port, _ := strconv.Atoi(v[7:])
	// 		ServerConf.HttpPort = port
	// 	}
	// 	if strings.Contains(v, "queue=") {
	// 		MachineryConf.DefaultQueue = v[8:]
	// 	}
	// }
	return env
}

// ServerConf MapTo maps section to given struct.
func serverConfMapTo(cfg *ini.File) {
	serverSection, err := cfg.NewSection("server")
	if err != nil {
		log.Fatalf("New server section err: %v", err)
	}
	// 后退判断！
	if ServerConf.HttpPort == 0 {
		serverSection.NewKey("HttpPort", os.Getenv("POST"))
	}
	serverSection.NewKey("RunMode", os.Getenv("RUNMODE"))
	serverSection.NewKey("ReadTimeout", os.Getenv("READTIMEOUT"))
	serverSection.NewKey("WriteTimeout", os.Getenv("WRITETIMEOUT"))
	err = cfg.Section("server").MapTo(ServerConf)
	if err != nil {
		log.Fatalf("Cfg.MapTo server err: %v", err)
	}
}

// ZapConf MapTo maps section to given struct.
func zapConfMapTo(cfg *ini.File) {
	zapSection, err := cfg.NewSection("zap")
	if err != nil {
		log.Fatalf("New zap section err: %v", err)
	}
	zapSection.NewKey("MaxSize", os.Getenv("MAXSIZE"))
	zapSection.NewKey("MaxBackups", os.Getenv("MAXBACKUPS"))
	zapSection.NewKey("MaxAge", os.Getenv("MAXAGE"))
	err = cfg.Section("zap").MapTo(ZapConf)
	if err != nil {
		log.Fatalf("Cfg.MapTo zap err: %v", err)
	}
}
