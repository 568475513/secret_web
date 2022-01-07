package sub_business

import (
	"github.com/jinzhu/gorm"
)

type ResourceDesc struct {
	Id                     int    `json:"id"`
	ImgUrl                 string `json:"img_url"`
	State                  uint8  `json:"state"`
	Title                  string `json:"title"`
	ImgUrlCompressed       string `json:"img_url_compressed"`
	ImgUrlCompressedLarger string `json:"img_url_compressed_larger"`
	OrgSummary             string `json:"org_summary"`
	OrgDescrb              string `json:"org_descrb"`
	Descrb                 string `json:"descrb"`
	AliveImgUrl            string `json:"alive_img_url"`
	Summary                string `json:"summary"`
}

// 设置表名 ResourceDesc
func (ResourceDesc) TableName() string {
	return DatabaseMicroPage + ".t_resource_desc"
}

// 获取资源resource_desc
func GetSpecInfo(appId string, resourceId string) (*ResourceDesc, error) {
	var rd ResourceDesc
	err := db.Select("id,img_url,state,title,img_url_compressed,img_url_compressed_larger,org_summary,org_descrb,descrb").
		Where("app_id=? and resource_id=? and state in (?)", appId, resourceId, []int{0, 1}).
		Take(&rd).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return &rd, nil
}

// 获取资源resource_desc By ids
func GetSpecInfoByIds(appId string, ids []string) ([]*ResourceDesc, error) {
	var rds []*ResourceDesc
	err := db.Select("id,img_url,state,title,img_url_compressed,img_url_compressed_larger,org_summary,org_descrb,descrb").
		Where("app_id=? and resource_id in (?) and state<>?", appId, ids, 2).
		Take(&rds).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return rds, nil
}
