package service

import (
	//内部包
	"abs/pkg/app"
	"abs/pkg/enums"
	"abs/pkg/logging"
	"errors"

	//系统标准包
	"fmt"
	"os"
	"time"
)

const (
	subscribeTimeOut         = 1000                              //超时设置，单位ms
	showMultipleSubscribeUrl = "api/subscribe/show_multiple_ids" //查询多个资源的订阅关系
)

type MultipleSubscribeResponse struct {
	app.Response
	Data MultipleSubscribeData `json:"data"`
}

type MultipleSubscribeData struct {
	Id []string `json:"id"`
}

//查询多个资源的订阅关系
func GetMultipleSubscribe(appId string, universalUnionId string, resourceIds []string) ([]string, error) {
	var result MultipleSubscribeResponse
	request := Post(fmt.Sprintf("%s%s?app_id=%s", os.Getenv("LB_PF_SUBSCRIBE_IN"), showMultipleSubscribeUrl, appId))
	request.SetParams(map[string]interface{}{
		"universal_union_id": universalUnionId,
		"resource_ids":       resourceIds,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(subscribeTimeOut * time.Millisecond)
	err := request.ToJSON(&result)
	if err != nil {
		logging.Info(err)
		return []string{}, err
	}
	if result.Code != enums.SUCCESS {
		return nil, errors.New(fmt.Sprintf("请求接口：%s错误:%s", showMultipleSubscribeUrl, result.Msg))
	}
	return result.Data.Id, nil
}