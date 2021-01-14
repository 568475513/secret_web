package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"abs/pkg/enums"
)

// 直播模板消息服务
type TemplateMsgService struct {
	// br    BaseRequest
	AppId string
}

const (
	// 场景id 108-直播预约
	sceneId = "108"
	// 查询单个业务场景开关
	switchGetState = "xe.msg.b_user.scene.switch.status/2.0.0"
	// 超时设置ms
	timeoutSwitchState = 500
)

// 发送请求查询短信预约总开关【场景单一，不设参数】
// func (t *TemplateMsgService) GetSwitchStateOld() (result map[string]interface{}, err error) {
// 	bT := time.Now()
// 	done := make(chan bool)
// 	data := map[string]string{"app_id": t.AppId, "scene_id": sceneId}
// 	go func(data map[string]string) {
// 		defer close(done)
// 		result, err = t.br.Request("POST", fmt.Sprintf("%s%s", os.Getenv("LB_CT_MSGOUT_IN"), switchGetState), data)
// 	}(data)

// 	select {
// 	// 监听发送请求是否超时, 如果超时，则记录数据.
// 	case <-time.After(timeoutSwitchState * time.Millisecond):
// 		err = errors.New("-- GetSwitchState Timeout --")
// 		t.br.RecordFailData(data)
// 		return
// 	// 数据在规定时间内已经请求业务侧.
// 	case <-done:
// 		eT := time.Since(bT)
// 		fmt.Printf("模板总开关请求: %s%s[GetSwitchState cos time] - %s\n", os.Getenv("LB_CT_BUZSERVER_IN"), switchGetState, eT)
// 		return
// 	}
// }

// 发送请求查询短信预约总开关
func (t *TemplateMsgService) GetSwitchState() (map[string]interface{}, error) {
	// var result map[string]interface{}
	request := Post(fmt.Sprintf("%s%s", os.Getenv("LB_CT_MSGOUT_IN"), switchGetState))
	request.SetParams(map[string]interface{}{
		"app_id":   t.AppId,
		"scene_id": sceneId,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(timeoutSwitchState * time.Millisecond)
	result, err := request.ToMap()
	if err != nil {
		return result, err
	}
	if v, ok := result["code"]; ok && int(v.(float64)) != enums.SUCCESS {
		return result, errors.New(fmt.Sprintf("请求短信预约总开关错误：%s", result["msg"].(string)))
	}
	return result["data"].(map[string]interface{}), nil
}
