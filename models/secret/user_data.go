package user

import (
	"github.com/jinzhu/gorm"
	"time"
)

//用户数据结构体
type UserData struct {
	Model

	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	DomainType int       `json:"domain_type"`
	Domain     string    `json:"domain"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Prevent struct {
	PreventName int `json:"prevent_name"`
	PreventNum  int `json:"prevent_num"`
}

type PreventDomain struct {
	PreventDomain string `json:"prevent_domain"`
	PreventNums   int    `json:"prevent_nums"`
}

//获取用户拦截类型数
func GetPreventCountByUserId(userId string) (tcs []Prevent, err error) {

	var (
		tc Prevent
	)
	rs, err := db.Table("t_secret_user_data").Select("domain_type , count(id) as count").Where("user_id = ? ", userId).Group("domain_type").Rows()
	if err != nil && err != gorm.ErrRecordNotFound || rs == nil {
		return nil, nil
	}
	for rs.Next() {
		rs.Scan(&tc.PreventName, &tc.PreventNum)
		tcs = append(tcs, tc)
	}
	return
}

//获取用户类型详细数据
func GetPreventDetailByUserId(userId string, dt int) (ps []PreventDomain, err error) {

	var (
		p PreventDomain
	)
	rs, err := db.Table("t_secret_user_data").Select("count(id) as prevent_nums ,domain").Where("user_id = ? and domain_type=? ", userId, dt).Group("domain").Rows()
	if err != nil && err != gorm.ErrRecordNotFound || rs == nil {
		return ps, nil
	}
	for rs.Next() {
		rs.Scan(&p.PreventNums, &p.PreventDomain)
		ps = append(ps, p)
	}
	return ps, nil
}
