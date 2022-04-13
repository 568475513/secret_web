package groups

import (
	"github.com/gin-gonic/gin"

	"abs/internal/server/api/secret"
)

// 隐私项目web接口
func SecretBaseRouter(Router *gin.RouterGroup) {
	//
	apiRouterV2 := Router.Group("api/secret")
	{
		// 用户登录注册接口
		apiRouterV2.POST("login/1.0.0", secret.UserLogin)
		// 获取用户信息接口
		apiRouterV2.POST("user_prevent_info", secret.UserPreventInfo)
		// 用户拦截信息接口
		apiRouterV2.POST("domain_prevent_add", secret.DomainPrevent)
		// 用户数据存储缓存
		apiRouterV2.GET("prevent_script_week", secret.WeekUserDataScript)
		// 每日推送用户数据
		apiRouterV2.GET("prevent_script_day", secret.DayUserDataScript)
	}
}
