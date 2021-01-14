package business

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ShareRecords struct {
	Model
	Id           int
	AppId        string    `json:"app_id"`
	ResourceType int       `json:"resource_type"`
	ResourceId   string    `json:"resource_id"`
	ProductId    string    `json:"product_id"`
	UserId       string    `json:"user_id"`
	ShareUserId  string    `json:"share_user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetShareResource(appId string, productId string, shareUserId string) ([]*ShareRecords, error) {
	var shareResource []*ShareRecords

	err := db.Table("t_share_records").Select("resource_id").Group("resource_id").Where("app_id = ? and product_id = ? and share_user_id = ?",
		appId, productId, shareUserId).Find(&shareResource).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return shareResource, err
}
