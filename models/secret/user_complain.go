package user

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 用户信息结构体
type SecretComplain struct {
	Model

	Id string `json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SecretComplainData struct {
	UserId          string `json:"user_id"`
	ComplainType    int    `json:"complain_type"`
	ComplainMsg     string `json:"complain_msg"`
	ComplainContact string `json:"complain_contact"`
}

//注册用户信息
func InsertUserComplain(userId, complainMsg, complainContact string, complainType int) (err error) {
	var ui SecretComplainData
	ui.UserId = userId
	ui.ComplainMsg = complainMsg
	ui.ComplainContact = complainContact
	ui.ComplainType = complainType
	err = db.Table("t_secret_complain").Create(ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}
