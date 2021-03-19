package service

import (
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

const (
	// 批量获取营期信息
	termBatchInfo = "v1/term/batch_info"
	// 超时设置ms
	timeoutTerm = 1000
)

// 营期信息发送请求
// func (t *CampService) GetCampTermInfoOld(methods string, route string, data map[string]interface{}) (result map[string]interface{}, err error) {
// 	bT := time.Now()
// 	done := make(chan bool)
// 	go func(data map[string]interface{}) {
// 		defer close(done)
// 		result, err = t.br.Request(methods, fmt.Sprintf("%s%s", os.Getenv("LB_SP_TRA_IN"), route), data)
// 	}(data)

// 	select {
// 	// 监听发送请求是否超时, 如果超时，则记录数据.
// 	case <-time.After(timeoutTerm * time.Millisecond):
// 		err = errors.New("-- GetCampTermInfo Timeout --")
// 		t.br.RecordFailData(data)
// 		return
// 	// 数据在规定时间内已经请求业务侧.
// 	case <-done:
// 		eT := time.Since(bT)
// 		fmt.Printf("service GetCampTermInfo Run time: %s And Res: %+v\n", eT, result)
// 		return
// 	}
// }

// 请求直播营期信息
func (t *CampService) GetCampTermInfo(ids, fields []string) ([]*business.PayProducts, error) {
	var result CampResponse
	request := Post(fmt.Sprintf("%s%s", os.Getenv("LB_SP_TRA_IN"), termBatchInfo))
	request.SetParams(map[string]interface{}{
		"select_fields": fields,
		"ids": ids,
		"app_id": t.AppId,
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
