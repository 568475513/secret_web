package sub_business

import (
	"github.com/jinzhu/gorm"
)

type SvipResRelation struct {
	Id           int
	AppId        string
	SvipId       string
	ResourceId   string
	ResourceType uint8
	RightsType   uint8
	IsSelectShow uint8
	State        uint8
}

// 设置表名 SvipResRelation
func (SvipResRelation) TableName() string {
	return DatabaseSvip + ".t_svip_res_relation"
}

type Svip struct {
	Id             string
	AppId          string
	EffactiveRange uint8
}

// 设置表名 Svip
func (Svip) TableName() string {
	return DatabaseSvip + ".t_svip"
}

// 获取资源svip关联信息
func GetResourceSvipRelation(appId string, resourceId string, resourceType int) (*SvipResRelation, error) {
	var sr SvipResRelation
	err := db.Select("id,app_id,svip_id,resource_id,resource_type,rights_type,is_select_show,state").
		Where("app_id=? and resource_id=? and resource_type=? and state=?", appId, resourceId, resourceType, 0).
		First(&sr).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &sr, nil
}

// 获取店铺svip信息
func GetSvipList(appId string) ([]*Svip, error) {
	var s []*Svip
	err := db.Select("id,app_id,effactive_range").
		Where("app_id=? and state in (0,1)", appId).
		Find(&s).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return s, nil
}
