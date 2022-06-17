package user

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 用户信息结构体
type UserConfig struct {
	Model

	Id            string    `json:"id"`
	UserId        string    `json:"user_id"`
	IsBuy         int       `json:"is_buy"`
	IsBusMonitor  int       `json:"is_bus_monitor"`
	IsLargeData   int       `json:"is_large_data"`
	IsSpy         int       `json:"is_spy"`
	IsCollectInfo int       `json:"is_collect_info"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// 获取用户信息
func GetUserConfig(userId string) (uc *UserConfig, err error) {

	var ui UserConfig
	err = db.Table("t_secret_user_config").Where("user_id=?", userId).First(&ui).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &ui, nil
}
