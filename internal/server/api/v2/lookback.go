package v2

import (
	"abs/internal/server/repository/course"
	"abs/internal/server/repository/material"
	"abs/internal/server/rules/validator"
	"abs/pkg/app"
	"abs/pkg/enums"
	"fmt"
	"github.com/gin-gonic/gin"
)

/**
 * 获取直播回放链接接口
 */
func GetLookBack(c *gin.Context) {
	var (
		err error
		req validator.LookBackRuleV2
	)

	// 参数校验
	AppId := app.GetAppId(c)
	if err = app.ParseQueryRequest(c, &req); err != nil {
		return
	}
	// req.AliveId = c.Query("alive_id")
	// req.Client, err = strconv.Atoi(c.Query("client"))
	if req.Client == 0 {
		// 默认公众号
		req.Client = 1
	}
	// if AppId == "" || req.AliveId == "" {
	// 	app.FailWithMessage("内容已被删除", enums.Code_Db_Not_Find, c)
	// 	return
	// }

	//获取直播数据
	aliveRep := course.AliveInfo{AppId: AppId, AliveId: req.AliveId}
	aliveInfo, err := aliveRep.GetAliveInfo()
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("GetAliveInfo错误:%s", err.Error()), enums.ERROR, c)
		return
	}

	//获取直播状态
	aliveState := aliveRep.GetAliveLookBackStates(aliveInfo)

	//获取直播回放链接
	lookBackRep := material.LookBack{AppId: AppId, AliveId: req.AliveId}
	data := lookBackRep.GetLookBackUrl(aliveInfo, aliveState, req.Client)

	app.OkWithData(data, c)
}
