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
	IsPrevent        int
	HighRisk         string
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

//拦截信息列表结构体
type PreventList struct {
	PreventDomain    string `json:"domain"`
	DomainTag        string `json:"domain_tag"`
	DomainType       string `json:"domain_type"`
	DomainSource     string `json:"domain_source"`
	DomainSourceInfo string `json:"domain_source_info"`
	RiskLevel        string `json:"risk_level"`
	CreatedAt        string `json:"created_at"`
}

//拦截信息分类结构体
type PreventClassify struct {
	DomainTag        string `json:"domain_tag"`
	DomainSource     string `json:"domain_source"`
	DomainSourceInfo string `json:"domain_source_info"`
	Count            int    `json:"count"`
	CreatedAt        string `json:"created_at"`
}

//拦截信息分类结构体
type PreventClassifyDetail struct {
	PreventDomain    string `json:"domain"`
	DomainTag        string `json:"domain_tag"`
	DomainSource     string `json:"domain_source"`
	DomainSourceInfo string `json:"domain_source_info"`
	RiskLevel        string `json:"risk_level"`
	CreatedAt        string `json:"created_at"`
}

type PreventSwitch struct {
	UserId        string `json:"user_id"`
	IsBusMonitor  int    `json:"is_bus_monitor"`
	IsLargeData   int    `json:"is_large_data"`
	IsSpy         int    `json:"is_spy"`
	IsCollectInfo int    `json:"is_collect_info"`
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

//获取用户拦截信息列表
func (u *U) GetPreventListById() (ps []PreventList, err error) {

	//获取拦截类型
	d, err := secret.GetAllDomainType()

	rs, err := secret.GetAllPreventDetailByUserId(u.UserId, u.HighRisk, u.Page, u.PageSize)
	if err != nil {
		logging.Error(err)
		return
	}
	for _, v := range rs {
		ps = append(ps, PreventList{PreventDomain: v.PreventDomain, CreatedAt: time.Unix(v.CreatedAt.Unix(), 0).Format("2006-01-02 15:04:05"), DomainTag: v.DomainTag, DomainType: d[v.DomainType].DomainName, DomainSource: v.DomainSource, DomainSourceInfo: v.DomainSourceInfo, RiskLevel: v.RiskLevel})
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

//获取用户拦截类型信息
func (u *U) GetPreventClassifyById() (ps []PreventClassify, err error) {

	rs, err := secret.GetAllPreventClassifyByUserId(u.UserId)
	if err != nil {
		logging.Error(err)
		return
	}
	for _, v := range rs {
		ps = append(ps, PreventClassify{CreatedAt: time.Unix(v.CreatedAt.Unix(), 0).Format("2006-01-02 15:04:05"), DomainTag: v.DomainTag, DomainSource: v.DomainSource, DomainSourceInfo: v.DomainSourceInfo, Count: v.Count})
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

//获取用户拦截类型详情信息
func (u *U) GetPreventClassifyDetailById() (ps []PreventClassifyDetail, err error) {

	rs, err := secret.GetAllPreventClassifyDetailByUserId(u.UserId, u.DomainTag)
	if err != nil {
		logging.Error(err)
		return
	}
	for _, v := range rs {
		ps = append(ps, PreventClassifyDetail{PreventDomain: v.PreventDomain, CreatedAt: time.Unix(v.CreatedAt.Unix(), 0).Format("2006-01-02 15:04:05"), DomainTag: v.DomainTag, DomainSource: v.DomainSource, DomainSourceInfo: v.DomainSourceInfo, RiskLevel: v.RiskLevel})
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

	err = secret.InsertPreventInfo(u.UserId, u.UserIp, u.Domain, u.DomainTag, u.DomainSource, u.DomainSourceInfo, u.RiskLevel, list.DomainType, u.IsPrevent)
	if err != nil {
		logging.Error(err)
		return
	}
	return nil
}

//更新用户拦截开关
func (p *PreventSwitch) UpdateUserPreventSwitch() (err error) {

	ui, err := secret.UpdateUserConfig(p.UserId, p.IsLargeData, p.IsSpy, p.IsBusMonitor, p.IsCollectInfo)
	if err != nil || ui == nil {
		logging.Error(err)
		return
	}
	return nil
}
