package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"abs/pkg/conf"
	"abs/pkg/util"
)

// 日志实列
var (
	// 调用链路打印对象
	ZLogger  *zap.Logger
	// 封装日志打印对象
	ILogger  *zap.Logger
	ELogger  *zap.Logger
	// Es日志对象
	EsLogger *zap.Logger
	// Job日志对象
	JLogger  *zap.Logger
)

// 初始化调用日志
func InitLog() {
	// 错误致命日志
	writeSyncerE := getLogWriter(
		util.GetRuntimeDir()+getLogFilePath("error"), conf.ZapConf.MaxSize*2, conf.ZapConf.MaxBackups, conf.ZapConf.MaxAge)
	// 普通输出日志
	writeSyncerI := getLogWriter(
		util.GetRuntimeDir()+getLogFilePath("info"), conf.ZapConf.MaxSize, conf.ZapConf.MaxBackups, 1)
	encoder := getConsoleEncoder()
	// 不用提前设置Level吧
	// var level = new(zapcore.Level)
	// err := level.UnmarshalText([]byte(os.Getenv("ZAP_LEVEL")))
	// if err != nil {
	// 	log.Fatalf(fmt.Sprintf("level.UnmarshalText InitLogger failed, err: %v\n", err))
	// }
	// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncerE), // 打印到控制台和文件
	ELogger = zap.New(zapcore.NewCore(encoder, writeSyncerE, zapcore.ErrorLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	ILogger = zap.New(zapcore.NewCore(encoder, writeSyncerI, zapcore.DebugLevel), zap.AddCaller())
	fmt.Println(">>>初始化调用日志完成")
}

// 初始化Es日志
func InitEs() {
	pathFile := fmt.Sprintf("%s/%s_%s.log", os.Getenv("ES_LOG_PATCH"), os.Getenv("ES_LOG_NAME"), time.Now().Format(os.Getenv("TIMEFORMAT")))
	writeSyncer := getLogWriter(pathFile, conf.ZapConf.MaxSize*2, conf.ZapConf.MaxBackups, 3)
	encoder := getJsonEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	EsLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	fmt.Println(">>>初始化Es日志完成")
}

// 初始化调用链路日志
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