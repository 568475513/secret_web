package business

import (
	// "database/sql"
	// "time"

	"github.com/jinzhu/gorm"

	"abs/pkg/provider/json"
)

type AppConf struct {
	Model

	// 关联关系【慎用】
	AppModule AppModule `gorm:"ForeignKey:AppId;AssociationForeignKey:AppId" json:"-"`

	AppId                  string        `json:"app_id"`
	ShopLogo               string        `json:"shop_logo"`
	ShopName               string        `json:"shop_name"`
	IsNewer                uint8         `json:"isNewer"`
	Balance                int           `json:"balance"`
	ModuleConf             string        `json:"module_conf,omitempty"` // 如果为空置则忽略字段
	WxAppId                *string       `json:"wx_app_id"`
	WxAppAvatar            *string       `json:"wx_app_avatar"`
	WxAppName              *string       `json:"wx_app_name"`
	UseCollection          uint8         `json:"use_collection"`
	VersionType            uint8         `json:"version_type"`
	ExpireTime             json.JSONTime `json:"expire_time"`
	WxSecreteKey           string        `json:"wx_secrete_key"`
	WxAccessToken          *string       `json:"wx_access_token"`
	WxAccessTokenRefreshAt *string       `json:"wx_access_token_refresh_at"`
}

// 设置表名 AppConf
func (AppConf) TableName() string {
	return DatabaseConfig + ".t_app_conf"
}

type AppModule struct {
	HasInvite                 uint8  `json:"has_invite"`
	HasReward                 uint8  `json:"has_reward"`
	GiftBuy                   uint8  `json:"gift_buy"`
	GiftBuyHideIfNotAvailable uint8  `json:"gift_buy_hide_if_not_available"`
	HasDistribute             uint8  `json:"has_distribute"`
	HideSubCount              uint8  `json:"hide_sub_count"`
	IsShowResourcecount       uint8  `json:"is_show_resourcecount"`
	HideViewCount             uint8  `json:"hide_view_count"`
	IsPersonMessagePush       uint8  `json:"is_person_message_push"`
	HideCommentCount          uint8  `json:"hide_comment_count"`
	IsShowExerciseSystem      uint8  `json:"is_show_exercise_system"`
	UserInfoConfig            string `json:"userinfo_config"`
	CaptionDefine             string `json:"caption_define"`
	AuthenticState            uint8  `json:"authentic_state"`
}

// 设置表名 AppModule
func (AppModule) TableName() string {
	return DatabaseConfig + ".t_app_module"
}

type ShopConfig struct {
	Module string `json:"module"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// 设置表名 AppModule
func (ShopConfig) TableName() string {
	return DatabaseConfig + ".t_shop_config"
}

// 获取店铺详情
func GetAppConfDetail(appId string) (*AppConf, error) {
	var ac AppConf
	err := db.Select([]string{
		"shop_logo",
		"shop_name",
		"isNewer",
		"created_at",
		"balance",
		"module_conf",
		"wx_app_id",
		"wx_app_avatar",
		"wx_app_name",
		"use_collection",
		"expire_time",
		"version_type",
		"wx_secrete_key",
		"wx_access_token",
		"wx_access_token_refresh_at",
		"updated_at"}).
		Where("app_id=? and wx_app_type=?", appId, 1).First(&ac).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ac, nil
}

// 获取店铺AppModule
func GetAppModule(appId string) (*AppModule, error) {
	var am AppModule
	err := db.Select([]string{
		"has_invite",
		"has_reward",
		"gift_buy",
		"gift_buy_hide_if_not_available",
		"has_distribute",
		"hide_sub_count",
		"is_show_resourcecount",
		"hide_view_count",
		"is_person_message_push",
		"hide_comment_count",
		"is_show_exercise_system",
		"userinfo_config",
		"caption_define",
		"authentic_state"}).
		Where("app_id=?", appId).First(&am).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &am, nil
}
