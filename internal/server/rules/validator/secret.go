package validator

type SecretUserLoginRule struct {
	UserId string `form:"user_id" json:"user_id"` // 用户id
	UserIp string `form:"user_ip" json:"user_ip"` //拦截用户ip

}

type SecretUserInfoRule struct {
	UserId     string `form:"user_id" json:"user_id" binding:"required"`         // 用户id
	DomainType string `form:"domain_type" json:"domain_type" binding:"required"` //拦截类型
	UserIp     string `form:"user_ip" json:"user_ip"`                            //拦截用户ip
	PageSize   int    `form:"page_size" json:"page_size" binding:"required"`     // 页码大小
	Page       int    `form:"page" json:"page" binding:"required"`               // 页码
}

type DomainPreventRule struct {
	Domain     string `form:"domain" json:"domain" binding:"required"`           // 拦截域名
	UserId     string `form:"user_id" json:"user_id"`                            //拦截用户id
	UserIp     string `form:"user_ip" json:"user_ip" binding:"required"`         //拦截用户ip
	DomainType string `form:"domain_type" json:"domain_type" binding:"required"` //拦截类型
}
