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
	var li = course.ListInfo{
		AppId:            req.AppId,
		UniversalUnionId: req.UniversalUnionId,
	}
	aliveList, err = li.GetALiveListByTime(req.StartTime, req.EndTime)
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("查询直播列表错误:%s", err.Error()), enums.ERROR, c)
		return
	}

	//筛出当前用户已订阅的直播
	aliveList = li.GetSubscribedALiveList(aliveList)

	app.OkWithData(aliveList, c)
}
