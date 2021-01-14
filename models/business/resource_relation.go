package business

import (
	"github.com/jinzhu/gorm"

	"abs/pkg/provider/json"
)

type ProResRelation struct {
	// Model // 不需要先

	AppId         string              `json:"app_id"`
	ProductId     string              `json:"product_id"`
	ProductType   uint8               `json:"product_type"`
	ProductName   json.JSONNullString `json:"product_name"`
	ResourceType  uint8               `json:"resource_type"`
	ResourceId    string              `json:"resource_id"`
	RelationState uint8               `json:"relation_state"`
	IsTry         uint8               `json:"is_try"`
	OrderWeight   uint                `json:"order_weight"`
	ShowInMenu    uint8               `json:"show_in_menu"`
	IsTop         uint8               `json:"is_top"`
}

// 获取资源专栏关联消息
func GetResourceProducts(appId string, resourceId string) ([]*ProResRelation, error) {
	var rp []*ProResRelation
	err := db.Select([]string{
		"app_id",
		"resource_id",
		"resource_type",
		"relation_state",
		"product_type",
		"product_id",
		"product_name",
		"is_try"}).
		Where("app_id=? and resource_id=? and resource_type=? and relation_state=?", appId, resourceId, 4, 0).Find(&rp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return rp, nil
}
