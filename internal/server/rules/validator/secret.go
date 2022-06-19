package validator

type SecretUserLoginRule struct {
	UserId     string `form:"user_id" json:"user_id"`         // 用户id
	UserIp     string `form:"user_ip" json:"user_ip"`         //拦截用户ip
	RegisterId string `form:"register_id" json:"register_id"` //极光推送注册id
}

type SecretUserInfoRule struct {
	UserId     string `form:"user_id" json:"user_id" binding:"required"`         // 用户id
	DomainType string `form:"domain_type" json:"domain_type" binding:"required"` //拦截类型
	UserIp     string `form:"user_ip" json:"user_ip"`                            //拦截用户ip
	PageSize   int    `form:"page_size" json:"page_size" binding:"required"`     // 页码大小
	Page       int    `form:"page" json:"page" binding:"required"`               // 页码
}

type SecretUserInfoListRule struct {
	UserId   string `form:"user_id" json:"user_id" binding:"required"`     // 用户id
	HighRisk string `form:"high_risk" json:"high_risk"`                    // 是否高风险
	PageSize int    `form:"page_size" json:"page_size" binding:"required"` // 页码大小
	Page     int    `form:"page" json:"page" binding:"required"`           // 页码
}

type SecretUserInfoClassifyRule struct {
	UserId string `form:"user_id" json:"user_id" binding:"required"` // 用户id
}

type SecretUserInfoClassifyDetailRule struct {
	UserId    string `form:"user_id" json:"user_id" binding:"required"`       // 用户id
	DomainTag string `form:"domain_tag" json:"domain_tag" binding:"required"` //拦截类型
}

type SecretUserClassifySwitchRule struct {
	UserId        string `form:"user_id" json:"user_id" binding:"required"` // 用户id
	IsBusMonitor  int    `form:"is_bus_monitor" json:"is_bus_monitor"`
	IsLargeData   int    `form:"is_large_data" json:"is_large_data"`
	IsSpy         int    `form:"is_spy" json:"is_spy"`
	IsCollectInfo int    `form:"is_collect_info" json:"is_collect_info"`
}

type DomainPreventRule struct {
	Domain           string `form:"domain" json:"domain" binding:"required"`               // 拦截域名
	UserId           string `form:"user_id" json:"user_id"`                                //拦截用户id
	UserIp           string `form:"user_ip" json:"user_ip" binding:"required"`             //拦截用户ip
	DomainType       string `form:"domain_type" json:"domain_type" binding:"required"`     //拦截类型
	DomainTag        string `form:"domain_tag" json:"domain_tag" binding:"required"`       //拦截类型
	DomainSource     string `form:"domain_source" json:"domain_source" binding:"required"` //拦截类型
	DomainSourceInfo string `form:"domain_source_info" json:"domain_source_info"`          //拦截类型
	RiskLevel        string `form:"risk_level" json:"risk_level"`                          //拦截类型
}
