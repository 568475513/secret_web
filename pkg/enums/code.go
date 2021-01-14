package enums

const (
	SUCCESS        = 0
	ERROR          = 500
	TIMEOUT        = 504
	BAD_REQUEST    = 400
	FORBIDDEN      = 403
	INVALID_PARAMS = 422

	// Baseinfo 神奇的值。。。
	RESOURCE_IS_BAN   = -99
	RESOURCE_REDIRECT = 302

	// 老的状态码兼容（php搬迁）
	Code_Db_Not_Find = 6
)
