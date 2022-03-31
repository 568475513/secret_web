package user

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 用户信息结构体
type SecretUser struct {
	Model

	Id                   string    `json:"id"`
	UserId               string    `json:"user_id"`
	UserIp               string    `json:"user_ip"`
	UserDnsPreventDomain string    `json:"user_dns_prevent_domain"`
	UserPrice            float64   `json:"user_price"`
	UserName             string    `json:"user_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type UId struct {
	UserId               string `json:"user_id"`
	UserName             string `json:"user_name"`
	UserDnsPreventDomain string `json:"user_dns_prevent_domain"`
}

// 获取用户信息
func GetUserInfo(userId, userIp string) (su *SecretUser, err error) {

	var ui SecretUser
	if userId != "" {
		err = db.Table("t_secret_user").Where("user_id=?", userId).First(&ui).Error
	} else if userIp != "" {
		err = db.Table("t_secret_user").Where("user_ip=?", userIp).First(&ui).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &ui, nil
}

//注册用户信息
func RegisterUser(userId, userDPD string) (err error) {
	var ui UId
	ui.UserId = userId
	ui.UserDnsPreventDomain = userDPD
	err = db.Table("t_secret_user").Create(ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

//更新用户积分
func UpdateUserPrice(userId, userIp string, price float64) (err error) {

	if userId != "" {
		err = db.Table("t_secret_user").Where("user_id=?", userId).Update("user_price", price).Error
	} else if userIp != "" {
		err = db.Table("t_secret_user").Where("user_ip=?", userIp).Update("user_price", price).Error
	}
	return
}
