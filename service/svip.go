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
	// 接口文档 http://doc.xiaoeknow.com/web/#/529
	svipBindResUrl = "xe.svip.resource.bind.relation/1.0.0" // 根据资源信息获取绑定的超级会员列表
	svipTimeOut    = 1000                                   //超时设置ms
)

type SvipService struct {
	AppId string
}

type SvipBindResResponse struct {
	app.Response
	Data []SvipBindResInfo `json:"data"`
}

type SvipBindResInfo struct {
	SvipID         string `json:"shop_id"`
	RightsType     string `json:"rights_type"`
	IsSelectShow   string `json:"is_select_show"`
	SvipName       string `json:"svip_name"`
	EffactiveRange string `json:"effactive_range"`
	IsDiscount     string `json:"is_discount"`
}

func (s *SvipService) GetSvipBindRes(resourceId string, resourceType int) ([]SvipBindResInfo, error) {
	var result SvipBindResResponse

	request := Post(fmt.Sprintf("%s%s", os.Getenv("LB_SP_SVIPSYSTEM_IN"), svipBindResUrl))

	request.SetParams(map[string]interface{}{
		"app_id":        s.AppId,
		"resource_id":   resourceId,
		"resource_type": resourceType,
	})

	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(svipTimeOut * time.Millisecond)

	err := request.ToJSON(&result)

	if err != nil {
		return []SvipBindResInfo{}, err
	}

	if result.Code != enums.SUCCESS {
		return []SvipBindResInfo{}, errors.New(fmt.Sprintf("请求接口：%s错误:%s", svipBindResUrl, result.Msg))
	}
	return result.Data, nil
}
