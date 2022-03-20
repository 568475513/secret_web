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

//根据域名类型获取名称
func GetDomainNameByType(Types []int) (ui []DomainType, err error) {

	err = db.Table("t_secret_domain_type").Where("domain_type IN (?)", Types).Find(&ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return ui, err
	}
	return ui, nil
}
