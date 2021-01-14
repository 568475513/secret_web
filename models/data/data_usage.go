package data

import (
	"fmt"
	"time"
)

type DataUage struct {
	AppId          string  `json:"app_id"`
	UserId         string  `json:"user_id"`
	ResourceId     string  `json:"resource_id"`
	ResourceType   int     `json:"resource_type"`
	ResourceName   string  `json:"resource_name"`
	Size           float32 `json:"size"`
	SizeCompressed float32 `json:"size_compressed"`
	ImgSizeTotal   float32 `json:"img_size_total"`
	SizeTotal      float32 `json:"size_total"`
	WxAppType      int     `json:"wx_app_type"`
	Way            int     `json:"way"`
	Model
}

func InsertFlowRecord(dataUage DataUage) bool {
	err := db.Create(&dataUage).Error
	if err != nil {
		return false
	}
	return true
}

func CreateDataUsageTable() bool {
	db.Set("gorm:table_options", "ENGINE=InnoDB")
	err := db.Exec(`CREATE TABLE IF NOT EXISTS ` + GetTableName() + ` (
	id INT(12) NOT NULL AUTO_INCREMENT COMMENT 'id',
		app_id VARCHAR(64) NOT NULL COMMENT '应用Id',
		user_id VARCHAR(64) NOT NULL COMMENT '用户Id',
		resource_id VARCHAR(64) NOT NULL DEFAULT '' COMMENT '资源id',
		resource_type INT(11) NOT NULL COMMENT '资源类型：0-无、1-音频、2-视频、3-直播、4-图文、5-直播回放',
		resource_name VARCHAR(128) DEFAULT NULL COMMENT '资源名',
		size FLOAT DEFAULT '0' COMMENT '消耗流量(M)',
		size_compressed FLOAT DEFAULT '0' COMMENT '压缩后的流量(M)',
		img_size_total FLOAT DEFAULT '0' COMMENT '图片大小(M)',
		size_total FLOAT DEFAULT '0' COMMENT '资源大小和图片大小之和(M)',
		wx_app_type INT(11) DEFAULT '1' COMMENT '数据来源 0-小程序 1-公众号',
		way INT(11) NOT NULL DEFAULT '0' COMMENT '统计方式：0-前端页面上报  1-后台直接上报 2-循环播放统计',
		created_at TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间，有修改自动更新',
		PRIMARY KEY (id),
		KEY index_user (app_id,resource_type,resource_id,user_id)
	) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='流量使用表'`)
	if err == nil {
		return false
	}
	return true
}

func (DataUage) TableName() string {
	return GetTableName()
}

func GetTableName() string {
	time.Now().Year()
	time.Now().Month()
	x := time.Unix(time.Now().Unix(), 0)
	return "t_data_usage_" + x.Format("2006_01_02")
}

func IsHaveTable(tableName string) bool {
	fmt.Println("mingz")
	return db.HasTable(tableName)
}
