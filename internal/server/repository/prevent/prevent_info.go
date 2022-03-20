package prevent

import (
	secret "abs/models/secret"
	"abs/pkg/logging"
)

type Uid struct {
	UserId string
}

//拦截信息结构体
type Prevent struct {
	PreventName string                 `json:"prevent_name"`
	PreventType int                    `json:"prevent_type"`
	PreventNum  int                    `json:"prevent_num"`
	PreventData []secret.PreventDomain `json:"prevent_data"`
}

func (u *Uid) GetPreventById() (ps []Prevent, err error) {

	var (
		p     Prevent
		Types []int
	)
	pi, err := secret.GetPreventCountByUserId(u.UserId)
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
		p.PreventData, err = secret.GetPreventDetailByUserId(u.UserId, v.DomainType)
		ps = append(ps, p)
	}
	return
}
