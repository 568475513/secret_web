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
	//获取已订阅且正在直播中的直播
	apiRouterV2.GET("get.living.alive.list", v2.GetLivingAliveList)
	//获取已经订阅的未开始直播
	apiRouterV2.GET("get.unstart.alive.list", v2.GetSubscribeUnStartAliveList)
	//获取主业务接口
	apiRouterV2.GET("base_info", v2.GetBaseInfo)
	// 获取次级业务接口
	apiRouterV2.GET("secondary_info", v2.GetSecondaryInfo)
	// 获取直播回放链接接口
	apiRouterV2.GET("get_lookback_url", v2.GetLookBack)
	// 获取课件使用记录接口
	apiRouterV2.GET("get_courseware_records", v2.GetCourseWareRecords)
	// 获取课件列表数据接口
	apiRouterV2.GET("get_courseware_info", v2.GetCourseWareInfo)
}
