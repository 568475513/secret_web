//只要是C端获取直播相关列表的接口逻辑就放这里，不然弄死
package v2

import (
	//内部包
	"abs/internal/server/repository/course"
	"abs/internal/server/rules/validator"
	"abs/models/alive"
	"abs/pkg/app"
	"abs/pkg/enums"
	//系统标准包
	"fmt"

	//第三方包
	"github.com/gin-gonic/gin"
)

//根据时间获取用户已订阅直播列表
func GetSubscribeAliveListByDate(c *gin.Context) {
	var (
		err       error
		req       validator.GetSubscribeAliveListByDateV2
		aliveList []*alive.Alive
	)

	//校验，不通过就给爷爬
	err = app.ParseQueryRequest(c, &req)
	if err != nil {
		return
	}

	//按时间查询指定店铺下的直播id
	li := course.ListInfo{
		AppId:            req.AppId,
		UniversalUnionId: req.UniversalUnionId,
	}
	aliveList, err = li.GetALiveListByTime(req.StartTime, req.EndTime, []string{"*"})
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("查询直播列表错误:%s", err.Error()), enums.ERROR, c)
		return
	}

	//筛出当前用户已订阅的直播
	result := li.GetSubscribedALiveList(aliveList)

	app.OkWithData(result, c)
}

//根据时间获取用户已订阅直播数量
func GetSubscribeAliveNumByDate(c *gin.Context) {
	var (
		err       error
		req       validator.GetSubscribeAliveNumByDateV2
		aliveList []*alive.Alive
	)

	//校验请求参数
	err = app.ParseQueryRequest(c, &req)
	if err != nil {
		return
	}

	//按时间查询指定店铺下的直播id
	li := course.ListInfo{
		AppId:            req.AppId,
		UniversalUnionId: req.UniversalUnionId,
	}
	aliveList, err = li.GetALiveListByTime(req.StartTime, req.EndTime, []string{"id", "zb_start_at"})
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("查询直播列表错误:%s", err.Error()), enums.ERROR, c)
		return
	}

	//筛出当前用户已订阅的直播
	subscribedALiveList := li.GetSubscribedALiveList(aliveList)

	//计数
	result := li.CountAliveNum(subscribedALiveList)

	app.OkWithData(result, c)
}
