package groups

import (
	"github.com/gin-gonic/gin"

	"abs/internal/server/api/v2"
)

// 直播详情页基本接口
func AliveBaseRouter(Router *gin.RouterGroup) {
	// v2版接口
	apiRouterV2 := Router.Group("_alive/v2")
	{
		// 获取主业务接口
		apiRouterV2.GET("base_info", v2.GetBaseInfo)
		// 获取次级业务接口
		apiRouterV2.GET("secondary_info", v2.GetSecondaryInfo)
		// 获取直播回放链接接口
		apiRouterV2.GET("get_lookback_url", v2.GetLookBack)
		// 获取课件使用记录接口
		apiRouterV2.GET("get_courseware_records", v2.GetCourseWareRecords)
		// 获取课件列表数据接口
		apiRouterV2.GET("get_courseware_info", v2.GetCourseWareInfo)
		// 数据上报接口
		apiRouterV2.GET("data_reported", v2.DataReported)

		// 压测调试接口
		// apiRouterV2.GET("base_info_test", v2.GetBaseInfoTest)
	}
}
