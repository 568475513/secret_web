package logging

import (
	"runtime/debug"

	"go.uber.org/zap"
)

// Info output logs at info level
func Info(v interface{}) {
	switch v.(type) {
	case string:
		BLogger.Info(v.(string))
	case map[string]interface{}:
		BLogger.Info("Map", zap.Any("info", v))
	default:
		BLogger.Info("Info!!!", zap.Any("Data", v))
	}
}

// Warn output logs at warn level
func Warn(param interface{}) {
	switch param.(type) {
	case string:
		BLogger.Warn(param.(string), zap.Any("error", param))
	case error:
		BLogger.Warn(param.(error).Error(), zap.Error(param.(error)), zap.Stack("stack"))
	default:
		BLogger.Warn("Warn!!!", zap.Any("warn", param), zap.Stack("stack"))
	}
}

// Error output logs at error level
func Error(param interface{}) {
	switch param.(type) {
	case string:
		BLogger.Error(param.(string), zap.String("stack", string(debug.Stack())))
	case error:
		BLogger.Error(param.(error).Error(),
			zap.Any("error", param.(error)),
			zap.String("stack", string(debug.Stack())),
		)
	default:
		BLogger.Error("未识别的 ERROR 类型",
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