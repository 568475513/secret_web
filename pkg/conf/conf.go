package conf

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
func Init() {
	fmt.Println(">开始初始化配置...")
	// 获取环境变量
	env := getEnvMode()

	// 读取对应环境变量
	err := godotenv.Load(".env." + env)
	fmt.Printf("Load config env:%s\n", env)
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
func getEnvMode() string {
	env := os.Getenv("GOENVMODE")
	// 命令调整参数
	for _, v := range os.Args {
		if strings.Contains(v, "env=") {
			env = v[4:]
		}
		if strings.Contains(v, "port=") {
			port, _ := strconv.Atoi(v[5:])
			ServerConf.HttpPort = port
		}
	}
	// 没有环境变量默认本地配置
	if env == "" {
		env = "local"
	}
	return env
}

// ServerConf MapTo maps section to given struct.
func serverConfMapTo(cfg *ini.File) {
	serverSection, err := cfg.NewSection("server")
	if err != nil {
		log.Fatalf("New server section err: %v", err)
	}
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
