package alive

import (
	// "strings"
	// "time"

	"abs/pkg/util"
	"github.com/jinzhu/gorm"
	"time"

	"abs/pkg/provider/json"
)

type Alive struct {
	Model

	AppId                  string              `json:"app_id"`
	Id                     string              `json:"alive_id"`
	RoomId                 string              `json:"room_id"`
	IsCompleteInfo         uint8               `json:"is_complete_info"`
	ProductId              json.JSONNullString `json:"product_id"`
	PaymentType            uint8               `json:"payment_type"`
	Summary                json.JSONNullString `json:"summary"`
	OrgContent             string              `json:"org_content"`
	Descrb                 string              `json:"descrb"`
	ZbStartAt              json.JSONTime       `json:"zb_start_at"`
	ZbStopAt               json.JSONTime       `json:"zb_stop_at"`
	ProductName            json.JSONNullString `json:"product_name"`
	Title                  json.JSONNullString `json:"title"`
	AliveVideoUrl          json.JSONNullString `json:"alive_video_url"`
	AliveImgUrl            json.JSONNullString `json:"alive_img_url"`
	ManualStopAt           json.JSONTime       `json:"manual_stop_at"`
	FileId                 string              `json:"file_id"`
	AliveType              uint8               `json:"alive_type"`
	ImgUrl                 json.JSONNullString `json:"img_url"`
	AliveroomImgUrl        json.JSONNullString `json:"aliveroom_img_url"`
	ImgUrlCompressed       json.JSONNullString `json:"img_url_compressed"`
	CanSelect              uint8               `json:"can_select"`
	DistributePercent      float64             `json:"distribute_percent"`
	HasDistribute          uint8               `json:"has_distribute"`
	DistributePoster       string              `json:"distribute_poster"`
	FirstDistributeDefault uint8               `json:"first_distribute_default"`
	FirstDistributePercent float64             `json:"first_distribute_percent"`
	RecycleBinState        uint8               `json:"recycle_bin_state"`
	State                  uint8               `json:"state"`
	StartAt                json.JSONTime       `json:"start_at"`
	IsStopSell             uint8               `json:"is_stop_sell"`
	ConfigShowViewCount    uint                `json:"config_show_view_count"`
	ConfigShowReward       uint                `json:"config_show_reward"`
	HavePassword           uint8               `json:"have_password"`
	IsDiscount             uint8               `json:"is_discount"`
	IsPublic               uint8               `json:"is_public"`
	PiecePrice             *uint               `json:"piece_price"`
	LinePrice              uint                `json:"line_price"`
	CommentCount           int                 `json:"comment_count"`
	ViewCount              int                 `json:"view_count"`
	ChannelId              string              `json:"channel_id"`
	PushState              uint8               `json:"push_state"`
	RewindTime             json.JSONTime       `json:"rewind_time"`
	PushUrl                string              `json:"push_url"`
	PlayUrl                string              `json:"play_url"`
	PptImgs                json.JSONNullString `json:"ppt_imgs"`
	PushAhead              string              `json:"push_ahead"`
	IsLookback             uint8               `json:"is_lookback"`
	IsTakegoods            uint8               `json:"is_takegoods"`
	IfPush                 uint8               `json:"if_push"`
	CreateMode             uint8               `json:"create_mode"`
	VideoLength            int64               `json:"video_length"`
	VideoSize              float64             `json:"video_size"`
	AliveM3u8HighSize      float64             `json:"alive_m3u8_high_size"`
	ForbidTalk             uint8               `json:"forbid_talk"`
	ShowOnWall             uint8               `json:"show_on_wall"`
	CanRecord              uint8               `json:"can_record"`

	// 非db数据
	IsTry uint8 `json:"is_try"`
}

type AliveRole struct {
	RoleName          json.JSONNullString `json:"role_name"`
	UserId            json.JSONNullString `json:"user_id"`
	UserName          json.JSONNullString `json:"user_name"`
	IsCurrentLecturer uint8               `json:"is_current_lecturer"`
	IsCanExceptional  uint8               `json:"is_can_exceptional"`
}

type AliveForbid struct {
	AppId    string `json:"app_id"`
	UseId    string `json:"user_id"`
	RoomId   string `json:"room_id"`
	IsUseful int    `json:"is_useful"`
}

type AliveModuleConf struct {
	AppId           string `json:"app_id"`
	AliveId         string `json:"alive_id"`
	IsCouponOn      uint8  `json:"is_coupon_on"`
	IsCardOn        uint8  `json:"is_card_on"`
	IsShowRewardOn  uint8  `json:"is_show_reward_on"`
	IsInviteOn      uint8  `json:"is_invite_on"`
	IsMessageOn     uint8  `json:"is_message_on"`
	IsSignInOn      uint8  `json:"is_sign_in_on"`
	IsMessageVerify uint8  `json:"is_message_verify"`
	IsPrizeOn       uint8  `json:"is_prize_on"`
	MessageAhead    int    `json:"message_ahead"`
	AliveMode       uint8  `json:"alive_mode"`
	CompleteTime    uint8  `json:"complete_time"`
	LookbackName    string `json:"lookback_name"`
	LookbackTime    string `json:"lookback_time"`
	IsRoundTableOn  uint8  `json:"is_round_table_on"`
}

const (
	//直播状态
	//ALIVE_STATE_UNSTART  = 0 //未开始
	AliveStateLiving = 1 //直播中
	//ALIVE_STATE_INTERACT = 2 //互动时间（已经断流，但还未回放）
	//ALIVE_STATE_END      = 3 //已结束（回放中）
	//ALIVE_STATE_LEAVE    = 4 //讲师离开
)

// 获取直播详情
func GetAliveInfo(appId string, aliveId string, s []string) (*Alive, error) {
	var a Alive
	err := db.Select(s).Where("app_id=? and id=?", appId, aliveId).First(&a).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &a, nil
}

// 通过channelId获取直播详情
func GetAliveInfoByChannelId(channelId string, s []string) (*Alive, error) {
	var a Alive
	err := db.Select(s).Where("channel_id=? and create_mode=0", channelId).First(&a).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &a, nil
}

// 获取直播讲师信息详情
func GetAliveRole(appId string, aliveId string) ([]*AliveRole, error) {
	var ar []*AliveRole
	err := db.Select("role_name,user_id,user_name,is_current_lecturer,is_can_exceptional").
		Where("app_id=? and alive_id=? and state=?", appId, aliveId, 0).
		Find(&ar).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return ar, nil
}

// 获取直播是否被禁言数据
func GetAliveForbid(appId, roomId, userId string) (*AliveForbid, error) {
	var af AliveForbid
	err := db.Select("is_useful").Where("app_id=? and room_id=? and user_id=?", appId, roomId, userId).First(&af).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &af, nil
}

// 获取直播配置
func GetAliveModuleConf(appId string, aliveId string, s []string) (*AliveModuleConf, error) {
	var amc AliveModuleConf
	err := db.Select(s).Where("app_id=? and alive_id=?", appId, aliveId).First(&amc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &amc, nil
}

// 更新直播观看人数
func UpdateViewCount(appId, aliveId string, viewCount int) error {
	var a Alive
	return db.Model(&a).Where("app_id=? and id=? and view_count<?", appId, aliveId, viewCount).
		Update("view_count", viewCount).
		Limit(1).Error
}

// 根据直播开始时间查询直播列表
func GetAliveListByZbStartTime(appId string, startTime string, endTime string, s []string) ([]*Alive, error) {
	var aliveList []*Alive
	err := db.Table("t_alive").Select(s).
		Where("app_id=? and zb_start_at>= ? and zb_start_at<=?", appId, startTime, endTime).Find(&aliveList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return aliveList, nil
}

// 根据app_id获取正在直播的直播
func GetLivingAliveListByAppId(appId string, s []string) ([]*Alive, error) {
	var (
		aliveList  []*Alive
		nowTimeStr = time.Now().Format(util.TIME_LAYOUT)
	)
	err := db.Table("t_alive").Select(s).
		Where("app_id=? and (push_state=? or (zb_start_at>? and zb_stop_at<? and manual_stop_at is null))", appId, AliveStateLiving, nowTimeStr, nowTimeStr).
		Find(&aliveList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return aliveList, nil
}
