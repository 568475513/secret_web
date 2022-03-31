package prevent

import (
	secret "abs/models/secret"
	"abs/pkg/logging"
	"abs/pkg/util"
)

type U struct {
	UserId     string
	UserIp     string
	Domain     string
	DomainType string
}

//拦截信息结构体
type Prevent struct {
	PreventName string                 `json:"prevent_name"`
	PreventType int                    `json:"prevent_type"`
	PreventNum  int                    `json:"prevent_num"`
	PreventData []secret.PreventDomain `json:"prevent_data"`
}

//获取用户拦截信息
func (u *U) GetPreventById() (ps []Prevent, err error) {

	var (
		p     Prevent
		Types []int
	)
	pi, err := secret.GetPreventCountByUserId(u.UserId, u.UserIp)
	if err != nil {
		logging.Error(err)
		return
	}

	TypesNums := make(map[int]int)
	for _, v := range pi {
		Types = append(Types, v.PreventName)
		TypesNums[v.PreventName] = v.PreventNum
	}
	domains, err := secret.GetDomainNameByType(Types)
	if err != nil {
		logging.Error(err)
		return
	}
	//TypeName :=  make(map[int]string)
	for _, v := range domains {
		p.PreventName = v.DomainName
		p.PreventNum = TypesNums[v.DomainType]
		p.PreventType = v.DomainType
		p.PreventData, err = secret.GetPreventDetailByUserId(u.UserId, u.UserIp, v.DomainType)
		ps = append(ps, p)
	}
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

	err = secret.InsertPreventInfo(u.UserId, u.UserIp, u.Domain, list.DomainType)
	if err != nil {
		logging.Error(err)
		return
	}
	return nil
}
