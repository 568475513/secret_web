package middleware

import (
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 处理全局的参数
func ReqParamHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 转化网关参数【用到的再设置】
		c.Set("app_id", c.GetHeader("XE_X_APP_ID"))
		c.Set("user_id", c.GetHeader("XE_X_USE_ID"))
		// c.Set("wx_app_id", c.GetHeader("XE_X_WX_APP_ID"))
		// c.Set("user_account_type", c.GetHeader("XE_X_USER_ACCOUNT_TYPE"))
		c.Set("app_version", c.GetHeader("XE_X_APP_VERSION"))
		// c.Set("force_collection", c.GetHeader("XE_X_FORCE_COLLECTION"))
		c.Set("buz_uri", c.GetHeader("XE_X_BUZ_URI"))
		c.Set("client_ip", c.GetHeader("XE_X_CLIENT_IP"))
		// c.Set("buz_referer", c.GetHeader("XE_X_BUZ_REFERER"))
		c.Set("agent", c.GetHeader("XE_X_AGENT"))
		// c.Set("is_manager", c.GetHeader("XE_X_IS_MANAGER"))
		client, _ := strconv.Atoi(c.GetHeader("XE_X_CLIENT"))
		agentType, _ := strconv.Atoi(c.GetHeader("XE_X_AGENT_TYPE"))
		c.Set("client", client)
		c.Set("agent_type", agentType)
		// c.Set("agent_version", c.GetHeader("XE_X_AGENT_VERSION"))

		//保存当前RequestId
		SetRequestId(c)

		// 暂时不这么用
		// 设置全局参数
		// 处理请求
		c.Next()
	}
}

//随机生成RequestId
func generateRequestId() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b[0:])
}