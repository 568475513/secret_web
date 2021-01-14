package sub_business

import (
	"github.com/jinzhu/gorm"
)

type ShopConfig struct {
	Module string `json:"module"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// 获取店铺Shop Config
func GetAppShopConfig(appId string) ([]*ShopConfig, error) {
	var sc []*ShopConfig
	err := db.Table(DatabaseShopConfig+".t_shop_config").
		Select("module,name,value").
		Where("app_id=? and is_deleted=?", appId, 0).
		Find(&sc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return sc, nil
}
