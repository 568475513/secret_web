package logging

import (
	"abs/pkg/conf"
	"github.com/gin-gonic/gin"
	"runtime/debug"

	"go.uber.org/zap"
)

// Info output logs at info level
func Info(v interface{}) {
	switch v.(type) {
	case string:
		BLogger.Info(v.(string),
		)
	case map[string]interface{}:
		BLogger.Info("Map",
			zap.Any("Data", v),
		)
	default:
		BLogger.Info("Info!!!",
			zap.Any("Data", v),
		)
	}
}

// Warn output logs at warn level
func Warn(param interface{}) {
	switch param.(type) {
	case string:
		BLogger.Warn(param.(string),
		)
	case error:
		BLogger.Warn(param.(error).Error(),
		)
	default:
		BLogger.Warn("Warn!!!",
			zap.Any("Data", param),
			zap.Stack("stack"),
		)
	}
}

// Error output logs at error level
func Error(param interface{}) {
	switch param.(type) {
	case string:
		BLogger.Error(param.(string),
			zap.String("stack", string(debug.Stack())),
		)
	case error:
		BLogger.Error(param.(error).Error(),
			zap.String("stack", string(debug.Stack())),
		)
	default:
		BLogger.Error("Error!!!",
			zap.Any("error", param.(error)),
			zap.String("stack", string(debug.Stack())),
		)
	}
}

// Info With Ctx output logs at info level
func InfoWithCtx(v interface{}, ctx *gin.Context) {
	switch v.(type) {
	case string:
		BLogger.Info(v.(string),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
		)
	case map[string]interface{}:
		BLogger.Info("Map",
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.Any("info", v),
		)
	default:
		BLogger.Info("Info!!!",
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.Any("Data", v),
		)
	}
}

// Warn With Ctx output logs at warn level
func WarnWithCtx(param interface{}, ctx *gin.Context) {
	switch param.(type) {
	case string:
		BLogger.Warn(param.(string),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
		)
	case error:
		BLogger.Warn(param.(error).Error(),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
		)
	default:
		BLogger.Warn("Warn!!!",
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.Any("warn", param),
			zap.Stack("stack"),
		)
	}
}

// Error With Ctx output logs at error level
func ErrorWithCtx(param interface{}, ctx *gin.Context) {
	switch param.(type) {
	case string:
		BLogger.Error(param.(string),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.String("stack", string(debug.Stack())),
		)
	case error:
		BLogger.Error(param.(error).Error(),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.String("stack", string(debug.Stack())),
		)
	default:
		BLogger.Error("Error!!!",
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.Any("error", param.(error)),
			zap.String("stack", string(debug.Stack())),
		)
	}
}

// 一般日志插入ES日志文件
//Deprecated
//func LogToEs(msg string, data interface{}) {
//	BLogger.Error(msg,
//		zap.Any("info", data),
//		zap.String("type", "info"),
//		zap.String("module_name", "alive_server_go"),
//		zap.String("method", "-"),
//		zap.String("target_url", "-"),
//		zap.String("request", "-"),
//	)
//}