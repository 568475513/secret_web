package user

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 用户信息结构体
type DomainType struct {
	Model

	Id         string    `json:"id"`
	DomainType int       `json:"domain_type"`
	DomainName string    `json:"domain_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

//获取所有拦截类型
func GetAllDomainType() (rs map[int]DomainType, err error) {

	var (
		d DomainType
	)
	rs = make(map[int]DomainType)
	r, err := db.Table("t_secret_domain_type").Rows()
	if err != nil && err != gorm.ErrRecordNotFound || r == nil {
		return nil, err
	}
	for r.Next() {
		r.Scan(&d.Id, &d.DomainType, &d.DomainName, &d.CreatedAt, &d.UpdatedAt)
		rs[d.DomainType] = d
	}
	return
}

//根据域名类型获取名称
func GetDomainNameByType(Types []int) (ui []DomainType, err error) {

	err = db.Table("t_secret_domain_type").Where("domain_type IN (?)", Types).Find(&ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return ui, err
	}
	return ui, nil
}

//获取域名类型
func GetDomainType(TypeName string) (ui DomainType, err error) {
	err = db.Table("t_secret_domain_type").Where("domain_name = ?", TypeName).Find(&ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return ui, err
	}
	return ui, nil
}
