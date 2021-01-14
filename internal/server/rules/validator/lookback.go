package validator

// 指明控制器所属！
// v2/lookback.go/GetLookBack

type LookBackRuleV2 struct {
	AliveId string `json:"alive_id" form:"alive_id" binding:"required,startswith=l_"`
	Client  int    `json:"client" form:"client" `
}
