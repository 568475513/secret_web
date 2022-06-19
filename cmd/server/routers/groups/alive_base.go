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
		// 用户登录注册接口
		apiRouterV2.POST("login/2.0.0", secret.UserLoginV2)
		// 获取用户拦截信息接口
		apiRouterV2.POST("user_prevent_info", secret.UserPreventInfo)
		// 获取用户拦截信息列表接口
		apiRouterV2.POST("user_prevent_info_list/2.0.0", secret.UserPreventInfoList)
		// 获取用户拦截信息分类接口
		apiRouterV2.POST("user_prevent_info_classify/2.0.0", secret.UserPreventInfoClassify)
		// 获取用户拦截信息分类详情接口
		apiRouterV2.POST("user_prevent_info_classify_detail/2.0.0", secret.UserPreventInfoClassifyDetail)
		// 获取用户拦截分类开关接口
		apiRouterV2.POST("user_prevent_classify_switch/2.0.0", secret.UserPreventClassifySwitch)
		// 用户问题反馈接口
		apiRouterV2.POST("user_complain/2.0.0", secret.UserComplain)
		// 用户拦截信息接口
		apiRouterV2.POST("domain_prevent_add", secret.DomainPrevent)
		// 用户配置信息接口
		apiRouterV2.POST("get_user_config", secret.GetUserConfigList)
		// 获取支持拦截类型
		apiRouterV2.POST("get_prevent_list", secret.GetPreventList)
		// 用户数据存储缓存
		apiRouterV2.GET("prevent_script_week", secret.WeekUserDataScript)
		// 每日推送用户数据
		apiRouterV2.GET("prevent_script_day", secret.DayUserDataScript)
		// 每周日推送用户数据
		apiRouterV2.GET("prevent_script_week_push", secret.DayUserDataPushScript)
	}
}
