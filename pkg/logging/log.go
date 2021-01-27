package logging

import (
	"fmt"
	"runtime/debug"

	"go.uber.org/zap"
)

// Info output logs at info level
func Info(v interface{}) {
	switch v.(type) {
	case string:
		ILogger.Info(v.(string))
	case map[string]interface{}:
		ILogger.Info("Map", zap.Any("info", v))
	default:
		ILogger.Info("Info!!!", zap.Any("Data", v))
	}
}

// Warn output logs at warn level
func Warn(param interface{}) {
	switch param.(type) {
	case string:
		ILogger.Warn(param.(string), zap.Any("error", param))
	case error:
		ILogger.Warn(param.(error).Error(), zap.Error(param.(error)), zap.String("stack", string(debug.Stack())))
	default:
		ILogger.Warn("Warn!!!", zap.Any("warn", param), zap.Stack("stack"))
	}
}

// Error output logs at error level
func Error(param interface{}) {
	switch param.(type) {
	case string:
		ELogger.Error(param.(string))
		EsLogger.Error(param.(string),
			zap.String("error", param.(string)),
			zap.String("type", "error"),
			zap.String("module_name", "alive_server_go"),
			zap.String("method", "-"),
			zap.String("target_url", "-"),
			zap.String("request", "-"),
			zap.String("stack", string(debug.Stack())),
		)
	case error:
		// ELogger.Error(param.(error).Error(), zap.Error(param.(error)), zap.String("stack", string(debug.Stack())))
		ELogger.Error(fmt.Sprintf("Error: %s\n stack: %s\n", param.(error).Error(), string(debug.Stack())))
		EsLogger.Error(param.(error).Error(),
			zap.Any("error", param.(error)),
			zap.String("type", "error"),
			zap.String("module_name", "alive_server_go"),
			zap.String("method", "-"),
			zap.String("target_url", "-"),
			zap.String("request", "-"),
			zap.String("stack", string(debug.Stack())),
		)
	default:
		ELogger.Error("Error!!!", zap.Any("error", param), zap.String("stack", string(debug.Stack())))
	}
}

// 一般日志插入ES日志文件
func LogToEs(msg string, data interface{}) {
	EsLogger.Error(msg,
		zap.Any("info", data),
		zap.String("type", "info"),
		zap.String("module_name", "alive_server_go"),
		zap.String("method", "-"),
		zap.String("target_url", "-"),
		zap.String("request", "-"),
	)
}