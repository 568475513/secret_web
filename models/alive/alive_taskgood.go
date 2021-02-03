package alive

import (
	"github.com/jinzhu/gorm"
)

type TaskGoodsDetail struct {
	Model

	AppId        string `json:"app_id"`
	AliveId      string `json:"alive_id"`
	ResourceId   string `json:"resource_id"`
	ResourceType int    `json:"resource_type"`
	ViewCount    int    `json:"view_count"`
	State        int    `json:"state"`
}

// 获取带货管理详情
func GetTaskGoodsInfo(appId, sourceId, resourceId string, s []string) (*TaskGoodsDetail, error) {
	var tgd TaskGoodsDetail

	err := db.Table("t_takegoods_detail").Select(s).Where("app_id=? and alive_id=? and resource_id=?", appId, sourceId, resourceId).First(&tgd).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &tgd, nil
}

// 初始化带货管理
func InsertTaskGoodsInfo(tgd TaskGoodsDetail) error {
	result := db.Table("t_takegoods_detail").Create(&tgd)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 更新带货PV
func UpdateTaskGoodsViewCount(appId, sourceId, resourceId string, viewCount int) error {
	err := db.Table("t_takegoods_detail").Where("app_id=? and id=? and resource_id=?", appId, sourceId, resourceId).Update("view_count", viewCount).Limit(1).Error
	if err != nil {
		return err
	}

	return nil
}