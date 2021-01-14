package business

import (
	"github.com/jinzhu/gorm"
)

type Channels struct {
	Model
	AppId       string `json:"app_id"`
	ChannelType int    `json:"channel_type"`
	ResourceId  string `json:"resource_id"`
	ProductId   string `json:"product_id"`
}

func GetChannelInfo(appId string, channelId string) (*Channels, error) {
	var channelInDB Channels
	err := db.Select("*").Where("id = ? and app_id = ?", channelId, appId).First(&channelInDB).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &channelInDB, err
}

func UpdateChannelViewCount(appId string, channelId string) error {
	err := dbRw.Table("t_channels").
		Where("id = ? and app_id = ?", channelId, appId).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error

	return err
}
