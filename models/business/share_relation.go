package business

import (
	"time"
	
	"github.com/jinzhu/gorm"
)

type ShareRelation struct {
	Model
	Id            int
	AppId         string    `json:"app_id"`
	PaymentType   string    `json:"payment_type"`
	ResourceId    string    `json:"resource_id"`
	ResourceType  int       `json:"resource_type"`
	ProductId     string    `json:"product_id"`
	ShareUserId   string    `json:"share_user_id"`
	ReceiveUserId string    `json:"receive_user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetAvailShareInfo(appId string, resourceType int, resourceId string, productId string, userId string) (*ShareRelation, error) {
	var availShareInfo ShareRelation
	err := db.Table("t_share_relation").Where("app_id = ? and payment_type = ? and resource_type = ? and resource_id = ? and product_id = ? and receive_user_id = ? and received_at is not null",
		appId, 3, resourceType, resourceId, productId, userId).First(&availShareInfo).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &availShareInfo, nil

}

func GetHasReceivedShare(appId string, resourceType int, resourceId string, productId string, userId string) ([]*ShareRelation, error) {
	var hasReceivedShare []*ShareRelation
	err := db.Table("t_share_relation").Select("receive_user_id").Where("app_id = ? and payment_type = ? and resource_type = ? and resource_id = ? and product_id = ? and share_user_id = ? and received_at is not null",
		appId, 3, resourceType, resourceId, productId, userId).Order("received_at").Find(&hasReceivedShare).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return hasReceivedShare, nil
}
