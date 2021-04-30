package groups

import (
	//第三方包
	"github.com/gin-gonic/gin"

	//内部包
	"abs/internal/server/api/v2"
)

func AliveAppRouter(Router *gin.RouterGroup) {
	// app写在abs_go的路由都加上v2
	apiRouterV2 := Router.Group("_alive/api-v2")
	//根据直播开始时间获取用户已订阅的直播课程列表
	apiRouterV2.GET("get.alive.list.by.date", v2.GetSubscribeAliveListByDate)
	//根据直播开始间获取用户已订阅的直播课程数量
	apiRouterV2.GET("get.alive.num.by.date", v2.GetSubscribeAliveNumByDate)
}
