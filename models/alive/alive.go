package alive

import (
	"abs/pkg/cache/redis_alive"
	"abs/pkg/provider/json"
	"abs/pkg/util"
	jsonUtil "encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"

	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// 直播角色列表缓存 key, app_id: alive_id,
const AliveRoleCacheKey = "alive:roles:%s:%s"

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
	IsTry      uint8 `json:"is_try"`
	AliveState int   `json:"alive_state"`
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
	AppId             string `json:"app_id"`
	AliveId           string `json:"alive_id"`
	IsCouponOn        uint8  `json:"is_coupon_on"`
	IsCardOn          uint8  `json:"is_card_on"`
	IsShowRewardOn    uint8  `json:"is_show_reward_on"`
	IsInviteOn        uint8  `json:"is_invite_on"`
	IsMessageOn       uint8  `json:"is_message_on"`
	IsSignInOn        uint8  `json:"is_sign_in_on"`
	IsMessageVerify   uint8  `json:"is_message_verify"`
	IsPrizeOn         uint8  `json:"is_prize_on"`
	MessageAhead      int    `json:"message_ahead"`
	AliveMode         uint8  `json:"alive_mode"`
	CompleteTime      uint16 `json:"complete_time"`
	LookbackName      string `json:"lookback_name"`
	LookbackTime      string `json:"lookback_time"`
	IsRedirectIndex   uint8  `json:"is_redirect_index"`
	IsRoundTableOn    uint8  `json:"is_round_table_on"`
	IsRedPacketOn     uint8  `json:"is_red_packet_on"`
	IsPictureOn       uint8  `json:"is_picture_on"`
	IsAuditFirstOn    uint8  `json:"is_audit_first_on"`
	IsOnlineOn        uint8  `json:"is_online_on"`
	IsHeatOn          uint8  `json:"is_heat_on"`
	IsOpenPromoter    uint8  `json:"is_open_promoter"`
	IsAntiScreen      uint8  `json:"is_anti_screen"`
	AliveJson         string `json:"alive_json"`
	IsOpenShareReward uint8  `json:"is_open_share_reward"`
	IsOpenQus         uint8  `json:"is_open_qus"`
	IsOpenVote        uint8  `json:"is_open_vote"`
	WarmUp            uint8  `json:"warm_up"`
}

type AliveTab struct {
	TabOn string `json:"tab_on"`
}

const (
	StateLiving    = 1 //直播中
	SubscribeState = 0 //查询用户已订阅直播
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
	con, err := redis_alive.GetLiveBusinessConn()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	cacheKey := fmt.Sprintf(AliveRoleCacheKey, appId, aliveId)
	data, _ := redis.Bytes(con.Do("GET", cacheKey))
	if data != nil {
		err = jsonUtil.Unmarshal(data, &ar)
	} else {
		err = db.Select("role_name,user_id,user_name,is_current_lecturer,is_can_exceptional").
			Where("app_id=? and alive_id=? and state=?", appId, aliveId, 0).
			Find(&ar).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		data, err = jsonUtil.Marshal(ar)
		if err == nil {
			con.Do("SET", cacheKey, data, "EX", 5)
		}
	}
	return ar, err
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
func GetLivingAliveListByAppId(appIds string, s []string) ([]*Alive, error) {
	var (
		aliveList []*Alive
	)
	appIdSlice := strings.Split(appIds, ",")
	if count := len(appIdSlice); count > 0 && count <= 5 {
		err := db.Table("t_alive").Select(s).
			Where("app_id in (?) and  push_state=?", appIdSlice, StateLiving).
			Find(&aliveList).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	} else {
		return nil, errors.New("app_id数量错误")
	}
	return aliveList, nil
}

// 根据app_id获取正在直播的直播
func GetUnStartAliveListByAppId(appIds string, s []string) ([]*Alive, error) {
	var (
		aliveList []*Alive
	)
	nowTime := time.Now().Format(util.TIME_LAYOUT)
	appIdSlice := strings.Split(appIds, ",")
	if count := len(appIdSlice); count > 0 && count <= 5 {
		err := db.Table("t_alive").Select(s).
			Where("app_id in (?) and zb_start_at > ? and manual_stop_at is null", appIdSlice, nowTime).
			Find(&aliveList).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	} else {
		return nil, errors.New("app_id数量错误")
	}
	return aliveList, nil
}

// 根据直播开始时间和直播类型查询直播列表
func GetAliveListByZbStartTimeAndType(appIds string, startTime string, endTime string, aliveType []string, s []string) ([]*Alive, error) {
	var aliveList []*Alive
	appIdSlice := strings.Split(appIds, ",")
	if count := len(appIdSlice); count > 0 && count <= 5 {
		err := db.Table("t_alive").Select(s).
			Where("app_id in (?) and zb_start_at>= ? and zb_start_at<=? and alive_type in (?)", appIdSlice, startTime, endTime, aliveType).
			Find(&aliveList).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	} else {
		return nil, errors.New("app_id数量错误")
	}
	return aliveList, nil
}

// CountAlive 资源计数
func CountAlive(appId string, resourceIds []string, startAt string) (total int, err error) {
	err = db.Table("t_alive").
		Where("app_id = ? and id in (?) and start_at < ? and state = 0", appId, resourceIds, startAt).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
