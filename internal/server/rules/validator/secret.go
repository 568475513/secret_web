package validator

type SecretUserLoginRule struct {
	UserId string `form:"user_id" json:"user_id"` // 用户id
}

type SecretUserInfoRule struct {
	UserId string `form:"user_id" json:"user_id" binding:"required"` // 用户id
}
