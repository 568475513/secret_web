package groups

import (
	//第三方包
	"github.com/gin-gonic/gin"

	//内部包
	"abs/internal/server/api/v2"
)

func AliveAppRouter(Router *gin.RouterGroup) {
	// app写在abs_go的路由都加上v2
	apiRouterV2 := Router.Group("_alive/api/v2")
	apiRouterV2.GET("get.alive.list.by.date", v2.GetSubscribeAliveListByDate)
}
