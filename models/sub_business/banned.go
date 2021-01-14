package sub_business

import (
	"github.com/jinzhu/gorm"
)

type BannedResource struct {
	ShelfState    uint8
	ResourceState uint8
}

// 资源是否被封禁
func ResourceIsBan(appId string, resourceId string) (bool, error) {
	var br BannedResource
	err := db.Select("resource_state").Where("app_id=? and resource_id=?", appId, resourceId).First(&br).Error
	// 报错默认返回未封禁
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// 被封禁
	if br.ResourceState > 0 {
		return true, nil
	}

	return false, nil
}
