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
		GetLogger().Info(v.(string))
	case map[string]interface{}:
		GetLogger().Info("Map",
			zap.Any("Data", v),
		)
	default:
		GetLogger().Info("Info!!!",
			zap.Any("Data", v),
		)
	}
}

// Warn output logs at warn level
func Warn(param interface{}) {
	switch param.(type) {
	case string:
		GetLogger().Warn(param.(string))
	case error:
		GetLogger().Warn(param.(error).Error())
	default:
		GetLogger().Warn("Warn!!!",
			zap.Any("Data", param),
			zap.Stack("stack"),
		)
	}
}

// Error output logs at error level
func Error(param interface{}) {
	switch param.(type) {
	case string:
		GetLogger().Error(param.(string),
			zap.String("stack", string(debug.Stack())),
		)
	case error:
		GetLogger().Error(param.(error).Error(),
			zap.String("stack", string(debug.Stack())),
		)
	default:
		GetLogger().Error("Error!!!",
			zap.Any("error", param.(error)),
			zap.String("stack", string(debug.Stack())),
		)
	}
}

// Info With Ctx output logs at info level
func InfoWithCtx(v interface{}, ctx *gin.Context) {
	switch v.(type) {
	case string:
		GetLogger().Info(v.(string),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
		)
	case map[string]interface{}:
		GetLogger().Info("Map",
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.Any("info", v),
		)
	default:
		GetLogger().Info("Info!!!",
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.Any("Data", v),
		)
	}
}

// Warn With Ctx output logs at warn level
func WarnWithCtx(param interface{}, ctx *gin.Context) {
	switch param.(type) {
	case string:
		GetLogger().Warn(param.(string),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
		)
	case error:
		GetLogger().Warn(param.(error).Error(),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
		)
	default:
		GetLogger().Warn("Warn!!!",
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
		GetLogger().Error(param.(string),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.String("stack", string(debug.Stack())),
		)
	case error:
		GetLogger().Error(param.(error).Error(),
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.String("stack", string(debug.Stack())),
		)
	default:
		GetLogger().Error("Error!!!",
			zap.String("requestId", ctx.GetString(conf.AbsRequestId)),
			zap.Any("error", param.(error)),
			zap.String("stack", string(debug.Stack())),
		)
	}
}

// 一般日志插入ES日志文件
//Deprecated
//func LogToEs(msg string, data interface{}) {
//	GetLogger().Error(msg,
//		zap.Any("info", data),
//		zap.String("type", "info"),
//		zap.String("module_name", "alive_server_go"),
//		zap.String("method", "-"),
//		zap.String("target_url", "-"),
//		zap.String("request", "-"),
//	)
//}
