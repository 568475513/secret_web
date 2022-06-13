package prevent

import (
	secret "abs/models/secret"
	"abs/pkg/logging"
	"abs/pkg/util"
	"time"
)

type U struct {
	UserId           string
	UserIp           string
	Domain           string
	DomainType       string
	DomainTag        string
	DomainSource     string
	DomainSourceInfo string
	RiskLevel        string
	Page             int
	PageSize         int
}

//拦截信息结构体
type Prevent struct {
	PreventDomain string `json:"domain"`
	DomainTag     string `json:"domain_tag"`
	DomainSource  string `json:"domain_source"`
	CreatedAt     string `json:"created_at"`
}

//type Prevent struct {
//	PreventName string                 `json:"prevent_name"`
//	PreventType int                    `json:"prevent_type"`
//	PreventNum  int                    `json:"prevent_num"`
//	PreventData []secret.PreventDomain `json:"prevent_data"`
//}

//获取用户拦截信息
func (u *U) GetPreventById() (ps []Prevent, err error) {

	rs, err := secret.GetPreventDetailByUserId(u.UserId, u.UserIp, u.DomainType, u.Page, u.PageSize)
	if err != nil {
		logging.Error(err)
		return
	}
	for _, v := range rs {
		ps = append(ps, Prevent{PreventDomain: v.PreventDomain, CreatedAt: time.Unix(v.CreatedAt.Unix(), 0).Format("2006-01-02 15:04:05"), DomainTag: v.DomainTag, DomainSource: v.DomainSource})
	}

	//var (
	//	p     Prevent
	//	//Types []int
	//)

	//pi, err := secret.GetPreventCountByUserId(u.UserId, u.UserIp)
	//if err != nil {
	//	logging.Error(err)
	//	return
	//}
	//
	//TypesNums := make(map[int]int)
	//for _, v := range pi {
	//	Types = append(Types, v.PreventName)
	//	TypesNums[v.PreventName] = v.PreventNum
	//}
	//domains, err := secret.GetDomainNameByType(Types)
	//if err != nil {
	//	logging.Error(err)
	//	return
	//}
	////TypeName :=  make(map[int]string)
	//for _, v := range domains {
	//	p.PreventName = v.DomainName
	//	p.PreventNum = TypesNums[v.DomainType]
	//	p.PreventType = v.DomainType
	//	p.PreventData, err = secret.GetPreventDetailByUserId(u.UserId, u.UserIp, v.DomainType)
	//	ps = append(ps, p)
	//}
	return
}

//记录用户拦截信息
func (u *U) InsertUserPreventInfo() (err error) {

	ui, err := secret.GetUserInfo(u.UserId, u.UserIp)
	if err != nil || ui == nil {
		logging.Error(err)
		return
	}

	list, err := secret.GetDomainType(u.DomainType)
	if err != nil {
		logging.Error(err)
		return
	}
	//获取用户所得积分
	ui.UserPrice = util.GetPrice(ui.UserPrice)

	err = secret.UpdateUserPrice(ui.UserId, ui.UserIp, ui.UserPrice)
	if err != nil {
		logging.Error(err)
		return
	}

	err = secret.InsertPreventInfo(u.UserId, u.UserIp, u.Domain, u.DomainTag, u.DomainSource, u.DomainSourceInfo, u.RiskLevel, list.DomainType)
	if err != nil {
		logging.Error(err)
		return
	}
	return nil
}
