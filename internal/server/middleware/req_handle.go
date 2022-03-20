package middleware

import (
	"abs/pkg/conf"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 处理全局的参数
func ReqParamHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		//注入RequestId，未来网关直接支持RequestId的话更佳
		c.Set(conf.AbsRequestId, GenerateRequestId())

		// 暂时不这么用
		// 设置全局参数
		// 处理请求
		c.Next()
	}
}

//GenerateRequestId 随机生成RequestId
func GenerateRequestId() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b[0:])
}
