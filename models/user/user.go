package user

import (
	"github.com/jinzhu/gorm"
)

// 此结构体只用做用户服务返回，暂时不查用户表了
type User struct {
	Model

	AppId            string `json:"app_id"`
	UserId           string `json:"user_id"`
	// 注意这里开发环境和现网类型不一样，现网要int，开发环境要string
	// AccountType      int    `json:"account_type"`
	WxOpenId         string `json:"wx_open_id"`
	WxUnicoId        string `json:"wx_union_id"`
	WxAppOpenId      string `json:"wx_app_open_id"`
	UniversalUnicoId string `json:"universal_union_id"`
	UniversalOpenId  string `json:"universal_open_id"`
	WxName           string `json:"wx_name"`
	WxNickname       string `json:"wx_nickname"`
	WxAvatar         string `json:"wx_avatar"`
	WxAccount        string `json:"wx_account"`
	Phone            string `json:"phone"`
}

// 获取用户信息
func GetUser(where map[string]interface{}, s []string) (*User, error) {
	var u User
	err := db.Select(s).Where(where).First(&u).Error
	// && err != gorm.ErrRecordNotFound 找不到用户直接抛错吧~
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// 获取多个用户信息
func GetUserList(appId string, ids []string, s []string) ([]*User, error) {
	var u []*User
	err := db.Select(s).Where("app_id=? and user_ids in (?)", appId, ids).First(&u).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return u, nil
}
