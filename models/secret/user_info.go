package user

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 用户信息结构体
type SecretUser struct {
	Model

	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UId struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
}

// 获取用户信息
func GetUserInfo(userId string) (*SecretUser, error) {
	var ui SecretUser
	err := db.Table("t_secret_user").Where("user_id=?", userId).First(&ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &ui, nil
}

//注册用户信息
func RegisterUser(userId string) (err error) {
	var ui UId
	ui.UserId = userId
	err = db.Table("t_secret_user").Create(ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}
