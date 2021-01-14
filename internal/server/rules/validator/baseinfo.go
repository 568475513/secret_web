package validator

// 指明控制器所属！
// v2/baseinfo.go/GetBaseInfo
type BaseInfoRuleV2 struct {
	AppId      string `form:"app_id" json:"app_id" binding:"startswith=app"`
	ResourceId string `form:"resource_id" json:"resource_id" binding:"required"`
	Type       string `form:"type" json:"type" binding:"required"`
	// ResourceType  string `form:"resource_type" json:"resource_type"`
	PaymentType int    `form:"payment_type" json:"payment_type" binding:"required"`
	ProductId   string `form:"product_id" json:"product_id"`
	// 老的BaseInfo兼容参数
	ChannelId    string `json:"channel_id" form:"channel_id"`
	ShareUserId  string `json:"share_user_id" form:"share_user_id"`
	ShareAgent   string `json:"share_agent" form:"share_agent"`
	ShareType    int    `json:"share_type" form:"share_type"`
	ShareFrom    string `json:"share_from" form:"share_from"`
	ContentAppId string `json:"content_app_id" form:"content_app_id"`
	MoreWay      string `json:"more_way" form:"more_way"`
}

// v2/baseinfo.go/GetSecondaryInfo
type SecondaryInfoRuleV2 struct {
	ResourceId  string `form:"resource_id" json:"resource_id" binding:"required"`
	ShareUserId string `json:"share_user_id" form:"share_user_id"`
	PaymentType int    `form:"payment_type" json:"payment_type" binding:"required"`
	ProductId   string `form:"product_id" json:"product_id"`
}

// v2/baseinfo.go/DataReported
type DataReportedV2 struct {
	AppId         string `form:"app_id" json:"app_id" binding:"startswith=app"`
	ResourceId    string `form:"resource_id" json:"resource_id" binding:"required"`   // 直播id
	ProductId     string `form:"product_id" json:"product_id"`                        // payment_type为2时-NULL, payment_type为3时-绑定的付费产品包id
	ChannelId     string `form:"channel_id" json:"channel_id"`                        // 渠道id
	PaymentType   string `form:"payment_type" json:"payment_type" binding:"required"` // 付费类型：2-单笔、3-付费产品包
	AppVersion    string `form:"app_version" json:"app_version"`
	Client        int8   `form:"client" json:"client"`                 // 来源终端
	UseCollection bool   `form:"use_collection" json:"use_collection"` // 判断小程序是否使用个人模式
}
