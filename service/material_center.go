package service

import (
	"abs/pkg/app"
	"abs/pkg/enums"
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	materialReplaceUrl = "api/xe.shop.material.replace.get/2.0.0" //批量替换素材接口
	materialTimeOut    = 1000                                     //超时设置ms
)

type MaterialResponse struct {
	app.Response
	Data MaterialReplace `json:"data"`
}

type MaterialReplace struct {
	ShopID     string            `json:"shop_id"`
	FilterData map[string]string `json:"filter_data"`
}

//获取批量替换素材url
func WashingData(appId string, data []string) (MaterialReplace, error) {
	var result MaterialResponse
	request := Post(fmt.Sprintf("%s%s?app_id=%s", os.Getenv("LB_PF_MATERIALTRANSFERCENTER_IN"), materialReplaceUrl, appId))
	request.SetParams(map[string]interface{}{
		"shop_id":     appId,
		"filter_data": data,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(materialTimeOut * time.Millisecond)
	err := request.ToJSON(&result)
	if err != nil {
		return MaterialReplace{}, err
	}
	if result.Code != enums.SUCCESS {
		return MaterialReplace{}, errors.New(fmt.Sprintf("请求接口：%s错误:%s", materialReplaceUrl, result.Msg))
	}
	return result.Data, nil
}
