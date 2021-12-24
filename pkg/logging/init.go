package logging

import (
	"abs/pkg/util"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"abs/pkg/conf"
)

// 日志实列
var (
	// 调用链路打印对象
	ZLogger  *zap.Logger
	// Job日志对象
	JLogger  *zap.Logger
)

// loggerPoolSize 预定义日志池大小
const loggerPoolSize = 100

// loggerPool 定义日志池
var loggerPool [loggerPoolSize]*zap.Logger

// 初始化调用日志
func InitLog() {
	for i := 0; i < loggerPoolSize; i++ {
		var pathFile string
		if conf.Env == "production" {
			//正式环境日志输出到指定目录 方便采集
			pathFile = fmt.Sprintf("%s/%s_%s_%d.log", os.Getenv("ES_LOG_PATCH"), os.Getenv("ES_LOG_NAME"), time.Now().Format(os.Getenv("TIMEFORMAT")), i)
		}else {
			//非正式环境日志输出到当前项目runtime目录 方便开发
			pathFile = util.GetRuntimeDir()+getLogFilePath("info")
		}

		writeSyncer := getLogWriter(pathFile, conf.ZapConf.MaxSize*2, conf.ZapConf.MaxBackups, 3)
		encoder := getJsonEncoder()
		core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
		loggerPool[i] = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	}
	fmt.Println(">>>初始化调用日志完成")
}

// 获取日志对象
func GetLogger() *zap.Logger {
	//获取随机数
	index := rand.Intn(loggerPoolSize)
	//通过随机数获取日志对象
	return loggerPool[index]
}

// 初始化Es日志
//Deprecated
func InitEs() {
	//pathFile := fmt.Sprintf("%s/%s_%s.log", os.Getenv("ES_LOG_PATCH"), os.Getenv("ES_LOG_NAME"), time.Now().Format(os.Getenv("TIMEFORMAT")))
	//writeSyncer := getLogWriter(pathFile, conf.ZapConf.MaxSize*2, conf.ZapConf.MaxBackups, 3)
	//encoder := getJsonEncoder()
	//core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	//EsLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	fmt.Println(">>>初始化Es日志完成")
}

// 初始化调用链路日志
// Deprecated
func InitZipkin() {
	pathFile := fmt.Sprintf("%s/%s_%s.log", os.Getenv("ZIPKIN_LOG_PATCH"), os.Getenv("ZIPKIN_LOG_NAME"), time.Now().Format(os.Getenv("TIMEFORMAT")))
	writeSyncerZ := getLogWriter(pathFile, conf.ZapConf.MaxSize*2, 5, 3)
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey: "msg",
	})
	ZLogger = zap.New(zapcore.NewCore(encoder, writeSyncerZ, zapcore.InfoLevel))
	fmt.Println(">>>初始化调用链路日志完成")
}

// 初始化Job打印日志对象
func InitJob() {
	pathFile := fmt.Sprintf("%s/job_abs_go_%s.log", os.Getenv("ES_LOG_PATCH"), time.Now().Format(os.Getenv("TIMEFORMAT")))
	writeSyncer := getLogWriter(pathFile, conf.ZapConf.MaxSize*2, conf.ZapConf.MaxBackups, 3)
	encoder := getJsonEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	JLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	fmt.Println(">>>初始化Job日志完成")
}

// getLogFileName get the save name of the log file
func getLogFilePath(category string) string {
	return fmt.Sprintf("logs/%s_%s_%s.%s",
		category,
		os.Getenv("LOGSAVENAME"),
		time.Now().Format(os.Getenv("TIMEFORMAT")),
		os.Getenv("LOGFILEEXT"),
	)
}

// 获取JsonEncoder
func getJsonEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 获取ConsoleEncoder
func getConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 获取LogWriter
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}