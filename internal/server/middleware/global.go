package middleware

import (
	"abs/pkg/conf"
	"bytes"
	"github.com/gin-gonic/gin"
	"runtime"
	"strconv"
)

//GlobalContext 妥协但相对实用的设计 目前主要是为了暂存 *gin.Context 方便打日志的时候获取 requestId
var GlobalContext map[uint64]*gin.Context

//SetCurrentContext 保存当前上下文到GlobalContext
func SetCurrentContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gID := getGID()
		GlobalContext[gID] = ctx

		defer delete(GlobalContext, gID)
		ctx.Next()
	}
}

//GetCurrentContext 获取当前上下文
func GetCurrentContext() *gin.Context {
	ctx, ok := GlobalContext[getGID()]
	if ok {
		return ctx
	}else {
		return nil
	}
}

//SetRequestId To gin.Context
func SetRequestId(ctx *gin.Context)  {
	//注入RequestId，未来网关直接支持RequestId的话更佳
	ctx.Set(conf.AbsRequestId, generateRequestId())
}

// GetRequestId 获取当前requestId
func GetRequestId() string{
	ctx := GetCurrentContext()
	if ctx == nil {
		return ""
	}
	return ctx.GetString(conf.AbsRequestId)
}

//获取协程ID
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}