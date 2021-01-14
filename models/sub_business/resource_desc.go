package sub_business

import (
	"github.com/jinzhu/gorm"
)

type ResourceDesc struct {
	Id                     int
	ImgUrl                 string
	State                  uint8
	Title                  string
	ImgUrlCompressed       string
	ImgUrlCompressedLarger string
	OrgSummary             string
	OrgDescrb              string
	Descrb                 string
	AliveImgUrl            string
}

// 设置表名 ResourceDesc
func (ResourceDesc) TableName() string {
	return DatabaseMicroPage + ".t_resource_desc"
}

// 获取资源resource_desc
func GetSpecInfo(appId string, resourceId string, resourceType int) (*ResourceDesc, error) {
	var rd ResourceDesc
	err := db.Select("id,img_url,state,title,img_url_compressed,img_url_compressed_larger,org_summary,org_descrb,descrb").
		Where("app_id=? and resource_id=? and resource_type=? and state in (?)", appId, resourceId, resourceType, [2]int{0, 1}).
		First(&rd).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &rd, nil
}

// 获取资源resource_desc By ids
func GetSpecInfoByIds(appId string, ids []string) ([]*ResourceDesc, error) {
	var rds []*ResourceDesc
	err := db.Select("id,img_url,state,title,img_url_compressed,img_url_compressed_larger,org_summary,org_descrb,descrb").
		Where("app_id=? and resource_id in (?) and state<>?", appId, ids, 2).
		First(&rds).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return rds, nil
}
