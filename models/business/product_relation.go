package business

import (
	"github.com/jinzhu/gorm"

	"abs/pkg/provider/json"
)

type ProductRelation struct {
	//Model

	AppId         string `json:"app_id"`
	ProductId     string `json:"product_id"`
	ProductType   int    `json:"product_type"`
	ProductName   string `json:"product_name"`
	ResourceType  int    `json:"resource_type"`
	ResourceId    string `json:"resource_id"`
	RelationState int    `json:"relation_state"`
	IsTry         int    `json:"is_try"`
	OrderWeight   int    `json:"order_weight"`
	IsTop         int    `json:"is_top"`

	//这里不直接继承model是为了方便service层json decode响应数据
	CreatedAt json.JSONTime `json:"created_at"`
	UpdatedAt json.JSONTime `json:"updated_at"`
}

const (
	tableName           = "t_pro_res_relation"
	relationStateNormal = 0
	relationStateDelete = 1
)

// GetResRelation 根据资源id查询父级关联关系
func GetResRelation(appId string, resourceId string, s []string) (data []*ProductRelation, err error) {
	err = db.Table(tableName).Select(s).
		Where("app_id = ? and resource_id = ? and relation_state = ?", appId, resourceId, relationStateNormal).
		Find(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return data, nil
}

// GetResByProductId 根据父级id查询关联关系
func GetResByProductId(appId string, productId string, s []string) (data []*ProductRelation, err error) {
	err = db.Table(tableName).Select(s).
		Where("app_id = ? and product_id = ? and relation_state = ?", appId, productId, relationStateNormal).
		Find(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return data, nil
}
