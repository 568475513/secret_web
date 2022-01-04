package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"abs/pkg/app"
	e "abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
)

type AvailableService struct {
	AppId  string
	UserId string
}

// 专栏权益参数
type ProductAvailable struct {
	ProductId      string // 必填
	AgentType      int
	NeedExpire     bool
	RelatedUserIds []string
	CreatedAt      string
	GenerateType   string
	ContentAppId   string
}

// 用户权益参数
type ResourceAvailable struct {
	ResourceType string // 必填
	ResourceId   string // 必填
	AgentType    int
	NeedExpire   bool
	ContentAppId string
	VersionType  int
}

// 用户权益参数
type ECourseAvailable struct {
	ResourceId string `json:"resource_id"`
}

const (
	// 专栏是否可用
	cmdIsProductAvailable = "/isProductAvailable"
	// 资源是否可用
	cmdIsResourceAvailable = "/isResourceAvailable"

	cmdIsECourseAvailable = "/xe.course.business.e_course.content.access.get/1.0.0"

	// 权益超时设置ms[time.Millisecond]
	availableTimeout = 1000
)

// 用户权益请求
func (ava *AvailableService) IsResourceAvailable(params ResourceAvailable) (expireAt string, available bool, err error) {
	// 企业微信全部免费！
	if !util.IsQyApp(params.VersionType) && params.AgentType == e.AGENT_TYPE_WW {
		available = true
		return
	}
	// 发起请求
	request := Post(fmt.Sprintf("%sceopenclose%s", os.Getenv("LB_CT_COPENCLOSE_IN"), cmdIsResourceAvailable))
	request.SetParams(map[string]interface{}{
		"appId":          ava.AppId,
		"userId":         ava.UserId,
		"resourceType":   params.ResourceType,
		"resourceId":     params.ResourceId,
		"needExpire":     params.NeedExpire,
		"content_app_id": params.ContentAppId,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(availableTimeout * time.Millisecond)
	result, err := request.ToMap()
	if err != nil {
		logging.Error(fmt.Sprintf("权益IsResourceAvailable，Http获取错误：%s", err.Error()))
		return
	}

	// 权益返回适配处理
	// now := time.Now()
	if result["code"].(float64) != 0 {
		// logging.Info("[资源]用户：" + ava.UserId + "未购买" + params.ResourceId + "_" + now.String())
	} else {
		// logging.Info("[资源]用户：" + ava.UserId + "已购买" + params.ResourceId + "_" + now.String())
		data := result["data"].(map[string]interface{})
		switch data["resource"].(type) {
		case string:
			expireAt = data["resource"].(string)
			available = true
		case bool:
			available = data["resource"].(bool)
		}
	}
	return
}

// 鹅课程权益请求
func (ava *AvailableService) IsECourseAvailable(params ECourseAvailable) (data interface{}, err error) {
	// 发起请求
	url := fmt.Sprintf(os.Getenv("LB_PF_COURSEBUSINESS_IN") + cmdIsECourseAvailable)
	request := Post(url)
	logging.Info(fmt.Sprintf("权益IsECourseAvailable，Url:" + url))
	// 写死请求数据 模拟
	//ava.AppId = "appm00slggh2325"
	//ava.UserId = "u_61b6f28beb92a_WPxIVGI6lS"
	//params.ResourceId = "v_61c97abb60b2567868b46b00"

	requestParams := map[string]interface{}{
		"app_id":   ava.AppId,
		"user_id":  ava.UserId,
		"buz_data": params,
	}

	request.SetParams(requestParams)
	logging.Info(fmt.Sprintf("权益IsECourseAvailable，params:"))
	logging.Info(requestParams)

	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(availableTimeout * time.Millisecond)
	result, err := request.ToMap()
	if err != nil {
		logging.Info(fmt.Sprintf("权益IsECourseAvailable，Http获取错误：%s", err.Error()))
		return
	}
	logging.Info(fmt.Sprintf("权益IsECourseAvailable，Http获取"))
	logging.Info(request)
	logging.Info(result)
	// 权益返回适配处理
	data = result["data"].(map[string]interface{})
	return
}

// 专栏权益请求
func (ava *AvailableService) IsProductAvailable(params ProductAvailable) (expireAt string, available bool) {
	// 参数判断
	if params.ProductId == "" || ava.AppId == "" || ava.UserId == "" {
		return
	}
	// 企业微信全部免费！
	if params.AgentType == e.AGENT_TYPE_WW {
		available = true
		return
	}

	// 发起请求
	request := Post(fmt.Sprintf("%sceopenclose%s", os.Getenv("LB_CT_COPENCLOSE_IN"), cmdIsProductAvailable))
	request.SetParams(map[string]interface{}{
		"appId":          ava.AppId,
		"userId":         ava.UserId,
		"relatedUserIds": params.RelatedUserIds,
		"createdAt":      params.CreatedAt,
		"generate_type":  params.GenerateType,
		"content_app_id": params.ContentAppId,
		"productId":      params.ProductId,
		"needExpire":     params.NeedExpire,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(availableTimeout * time.Millisecond)
	result, err := request.ToMap()
	if err != nil {
		logging.Error(fmt.Sprintf("权益IsProductAvailable，Http获取错误：%s", err.Error()))
		return
	}

	// 处理返回结果
	// now := time.Now()
	if result["code"].(float64) != 0 {
		// logging.Info("[专栏]用户：" + ava.UserId + "未购买" + params.ProductId + "_" + now.String())
	} else {
		// logging.Info("[专栏]用户：" + ava.UserId + "已购买" + params.ProductId + "_" + now.String())
		data := result["data"].(map[string]interface{})
		switch data["resource"].(type) {
		case string:
			expireAt = data["resource"].(string)
			available = true
		case bool:
			available = data["resource"].(bool)
		}
	}
	return
}

// 判断内部课程和加密课程方法
func (ava *AvailableService) IsResourceAccess(resourceId string, filterFree bool, needExpire int) (bool, error) {
	var responseMap app.Response
	request := Get(fmt.Sprintf("%sxe.user.permission.check", os.Getenv("LB_PF_RIGHTS_IN")))
	params := map[string]string{
		"region": ava.AppId,
		"sub":    ava.UserId,
		"obj":    resourceId,
		"expire": strconv.Itoa(needExpire),
		"app_id": ava.AppId,
	}
	if filterFree {
		params["source"] = "all"
	}
	request.SetParams(params)
	request.SetTimeout(availableTimeout * time.Millisecond)
	err := request.ToJSON(&responseMap)
	if err != nil || responseMap.Code != e.SUCCESS {
		logging.Error(err)
		return false, err
	}
	if responseMap.Code != e.SUCCESS {
		return false, errors.New(fmt.Sprintf("请求权益【内部课程和加密】信息错误：%s", responseMap.Msg))
	}
	data, ok := responseMap.Data.(map[string]interface{})
	if !ok {
		logging.Error(fmt.Sprintf("请求权益断言错误：%+v", responseMap))
	}
	AuthState, _ := data["auth_state"].(float64)
	return AuthState == 1, nil
}
