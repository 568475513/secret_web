package validator

type SecretUserLoginRule struct {
	UserId string `form:"user_id" json:"user_id"` // 用户id
	UserIp string `form:"user_ip" json:"user_ip"` //拦截用户ip

}

type SecretUserInfoRule struct {
	UserId string `form:"user_id" json:"user_id" binding:"required"` // 用户id
	UserIp string `form:"user_ip" json:"user_ip"`                    //拦截用户ip
}

type DomainPreventRule struct {
	Domain     string `form:"domain" json:"domain" binding:"required"`           // 拦截域名
	UserId     string `form:"user_id" json:"user_id"`                            //拦截用户id
	UserIp     string `form:"user_ip" json:"user_ip" binding:"required"`         //拦截用户ip
	DomainType string `form:"domain_type" json:"domain_type" binding:"required"` //拦截类型
}
