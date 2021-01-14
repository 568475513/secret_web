package data

import (
	"encoding/json"
	"fmt"
	"time"

	"abs/models/business"
	// "abs/pkg/logging"
	// "abs/pkg/util"
	"abs/service"
)

type UserPurchaseData struct {
	AppId          string `json:"app_id"`
	UserId         string `json:"user_id"`
	RawUrl         string `json:"raw_url"`
	Url            string `json:"url"`
	Referer        string `json:"referer"`
	AppVersion     string `json:"app_version"`
	Agent          string `json:"agent"`
	Client         int8
	UserCollection bool   `json:"use_collection"`
	Ip             string `json:"ip"`
	ResourceType   string `json:"resource_type"`
	ResourceId     string `json:"resource_id"`
	ProductId      string `json:"product_id"`
	IsResourcePay  bool   `json:"is_resource_pay"`
	Params         string `json:"params"`
}

type LogData struct {
	*UserPurchaseData
	IsPay           int8   `json:"is_pay"`
	CreatedAt       string `json:"created_at"`
	UniversalOpenId string `json:"universal_open_id"`
}

func InsertUserPurchaseLog(field *UserPurchaseData) {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(createdAt)
	userService := service.UserService{AppId: field.AppId, UserId: field.UserId}
	userInfoResponse, _ := userService.RequestUserInfo()
	var isPlay int8 = 0
	if isExist, err := business.ExistPurchaseRecord(field.AppId, field.UserId); isExist && err == nil {
		isPlay = 1
	}
	logData := LogData{
		IsPay:            isPlay,
		CreatedAt:        createdAt,
		UniversalOpenId:  userInfoResponse.Data.UniversalOpenId,
		UserPurchaseData: field,
	}
	jsonLog, _ := json.Marshal(logData)
	b := string(jsonLog)
	fmt.Println(b)
	//logging.WriterPurchaseLog(string(jsonLog))
	// logging.InitBussinessLogger(1, 3, 1024*1024, func() string {
	// 	var prefix string = "user_purchase_log_"
	// 	return prefix + time.Now().Format("2006_01_02")
	// }, util.GetRuntimeDir())
	// logging.Logger.SuperInfo(string(jsonLog))
}
