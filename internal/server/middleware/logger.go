package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"abs/pkg/app"
)

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		requestTime := time.Since(start)
		logger.Info("requestData",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("clientIp", c.ClientIP()),
			zap.String("userAgent", c.Request.UserAgent()),
			zap.String("requestTime", requestTime.String()),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					// nolint: errcheck
					c.Error(err.(error))
					c.Abort()
					return
				}

				logger.Error("[Recovery from panic]",
					zap.Any("error", err),
					zap.String("type", "panic"),
					zap.String("module_name", "alive_server_go"),
					zap.String("method", c.Request.Method),
					zap.String("target_url", c.Request.URL.Path),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)
				fmt.Printf("Error: [%s]\nstack: %s\n", err.(error).Error(), (debug.Stack()))
				app.FailWithMessage(fmt.Sprintf("[%s]\nstack: %s", err, string(debug.Stack())), http.StatusInternalServerError, c)
				// c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
