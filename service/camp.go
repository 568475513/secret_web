package service

import (
	"abs/pkg/logging"
	"errors"
	"fmt"
	"os"
	"time"

	"abs/models/business"
	"abs/pkg/app"
	"abs/pkg/enums"
)

// 训练营营期服务
type CampService struct {
	// br BaseRequest
	AppId string
}

// 序列化结构体返回
type CampResponse struct {
	app.Response
	Data CampListTerms `json:"data"`
}

// 向下序列化结构体
type CampListTerms struct {
	Terms []business.PayProducts `json:"terms"`
}

// CampResponseV2 序列化结构体返回v2
type CampResponseV2 struct {
	app.Response
	Data CampListTermsV2 `json:"data"`
}

// CampListTermsV2 向下序列化结构体v2
type CampListTermsV2 struct {
	Terms []map[string]interface{} `json:"terms"`
}

const (
	// 批量获取营期信息
	termBatchInfo = "v1/term/batch_info"
	// 超时设置ms
	timeoutTerm = 1000
)

// 请求直播营期信息
func (t *CampService) GetCampTermInfo(ids, fields []string) ([]*business.PayProducts, error) {
	var result CampResponse
	request := Post(fmt.Sprintf("%s%s", os.Getenv("LB_SP_TRA_IN"), termBatchInfo))
	request.SetParams(map[string]interface{}{
		"select_fields": fields,
		"ids":           ids,
		"app_id":        t.AppId,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(timeoutTerm * time.Millisecond)
	err := request.ToJSON(&result)
	trems := []*business.PayProducts{}
	if err != nil {
		return trems, err
	}
	if result.Code != enums.SUCCESS {
		return trems, errors.New(fmt.Sprintf("请求直播营期信息错误：%s", result.Msg))
	}

	// 转化为返回指针
	for _, v := range result.Data.Terms {
		// 过滤隐藏的
		if v.DisplayState == 0 && v.RecycleBinState == 0 {
			// var tmp business.PayProducts
			tmp := v
			trems = append(trems, &tmp)
		}
	}
	return trems, nil
}

// GetCampTermInfoV2 请求直播营期信息
// 上边的GetCampTermInfo应该是有问题，接口响应字段名和business.PayProducts对不上（现网没有问题？）
func (t *CampService) GetCampTermInfoV2(ids, fields []string) ([]map[string]interface{}, error) {
	var (
		result CampResponseV2
		terms  []map[string]interface{}
		err    error
	)

	request := Post(fmt.Sprintf("%s%s", os.Getenv("LB_SP_TRA_IN"), termBatchInfo))
	request.SetParams(map[string]interface{}{
		"select_fields": fields,
		"ids":           ids,
		"app_id":        t.AppId,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(timeoutTerm * time.Millisecond)
	err = request.ToJSON(&result)
	logging.Info(fmt.Sprintf("训练营数据terms:%s", result))

	if err != nil {
		return terms, err
	}
	if result.Code != enums.SUCCESS {
		return terms, errors.New(fmt.Sprintf("请求直播营期信息错误：%s", result.Msg))
	}

	return result.Data.Terms, nil
}
