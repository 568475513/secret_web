package business

import (
	"time"
)

type Purchase struct {
	AppId          string    `json:"app_id"`
	ContentAppId   string    `json:"content_app_id"`
	UserId         string    `json:"user_id"`
	PaymentType    int16     `json:"payment_type"`
	ResourceType   int16     `json:"resource_type"`
	ProductId      string    `json:"product_id"`
	GenerateType   int16     `json:"generate_type"`
	ResourceId     string    `json:"resource_id"`
	ChannelId      string    `json:"channel_id"`
	ShareUserId    string    `json:"share_user_id"`
	ShareType      int16     `json:"share_type"`
	PurchaseName   string    `json:"purchase_name"`
	ImgUrl         string    `json:"img_url"`
	Price          int16     `json:"price"`
	OrderId        string    `json:"order_id"`
	Remark         string    `json:"remark"`
	WxAppType      int16     `json:"wx_app_type"`
	ExpireAt       time.Time `json:"expire_at"`
	IsMemberRemind int8      `json:"is_member_remind"`
	IsDeleted      int16     `json:"is_deleted"`
	NeedNotify     int16     `json:"need_notify"`
	Agent          string    `json:"agent"`
	Model
}

// 查询用户是否有订购记录
func ExistPurchaseRecord(appId string, userId string) (bool, error) {
	var num []int
	err := db.Table("t_purchase").Select("1 AS num").Where("app_id = ? and user_id = ? and is_deleted != 1", appId, userId).Limit(1).Pluck("num", &num).Error
	return len(num) >= 1 && num[0] == 1, err
}

func (Purchase) TableName() string {
	return "t_purchase"
}
