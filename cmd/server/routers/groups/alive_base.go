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
		// 获取主业务接口
		apiRouterV2.POST("login/1.0.0", secret.UserLogin)
		// 获取次级业务接口
		apiRouterV2.POST("user_prevent_info", secret.UserPreventInfo)
	}
}
