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
	aliveList = li.GetSubscribedALiveList(aliveList)

	//直播按日期分组下
	result := li.ALiveListGroupByTime(aliveList)

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

	//直播按日期分组下
	subscribedALiveListGroupByDate := li.ALiveListGroupByTime(subscribedALiveList)

	//计数
	result := li.CountAliveNum(subscribedALiveListGroupByDate)

	app.OkWithData(result, c)
}

//获取用户已订阅且正在直播中的直播列表
func GetSubscribeLivingAliveList(c *gin.Context) {
	var (
		err       error
		req       validator.GetSubscribeLivingAliveListV2
		aliveList []*alive.Alive
	)

	//校验请求参数
	err = app.ParseQueryRequest(c, &req)
	if err != nil {
		return
	}

	//根据app_id获取正在直播中的直播
	li := course.ListInfo{
		UniversalUnionId: req.UniversalUnionId,
	}
	aliveList, err = li.GetLivingAliveList(req.AppId, []string{"*"})

	//筛出当前用户已订阅的直播
	subscribedALiveList := li.GetSubscribedALiveList(aliveList)

	app.OkWithData(subscribedALiveList, c)
}
