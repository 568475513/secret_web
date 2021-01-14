package sub_business

import (
	// "database/sql"
	// "time"

	"abs/pkg/provider/json"
	"github.com/jinzhu/gorm"
)

type CloneResApply struct {
	ContentAppId       string        `json:"app_id"`
	ChannelAppId       string        `json:"file_id"`
	ResourceId         string        `json:"video_url"`
	ResourceType       uint8         `json:"video_audio_url"`
	ResourceName       string        `json:"video_mp4"`
	ShopName           string        `json:"shop_name"`
	ApplyState         uint8         `json:"apply_state"`
	PassedAt           json.JSONTime `json:"passed_at"`
	State              uint8         `json:"state"`
	CreatedAt          json.JSONTime `json:"created_at"`
	UpdatedAt          json.JSONTime `json:"updated_at"`
	LookbackPermission uint8         `json:"lookback_permission"`
}

// 设置表名 CloneResApply
func (CloneResApply) TableName() string {
	return DatabaseContentMarket + ".t_clone_res_apply"
}

// 获取转播课程
func GetCloneResApply(ContentAppId string, ResourceId string, ChannelAppId string, s []string) (*CloneResApply, error) {
	var cra CloneResApply
	err := db.Select(s).Where("content_app_id=? and resource_id=? and channel_app_id=?", ContentAppId, ResourceId, ChannelAppId).First(&cra).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &cra, nil
}
