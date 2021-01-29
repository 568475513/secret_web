package alive

import (
	// "database/sql"
	// "time"

	"github.com/jinzhu/gorm"
	"strings"
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
	CommentCount           uint                `json:"comment_count"`
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
	IsMessageVerify uint8  `json:"is_message_verify"`
	IsPrizeOn       uint8  `json:"is_prize_on"`
	MessageAhead    uint8  `json:"message_ahead"`
	AliveMode       uint8  `json:"alive_mode"`
	CompleteTime    uint8  `json:"complete_time"`
	LookbackName    string `json:"lookback_name"`
	LookbackTime    string `json:"lookback_time"`
}

type AliveLookBack struct {
	Id             int    `json:"id"`
	AppId          string `json:"app_id"`
	AliveId        string `json:"alive_id"`
	LookbackFileId string `json:"lookback_file_id"`
	RegionFileId   string `json:"region_file_id"`
	LookbackMp4    string `json:"lookback_mp4"`
	LookbackM3u8   string `json:"lookback_m3u8"`
	FileName       string `json:"file_name"`
	TranscodeState uint8  `json:"transcode_state"`
	State          uint8  `json:"state"`
	OriginType     uint8  `json:"origin_type"`
}

type AliveConcatHlsResult struct {
	ChannelId                string `json:"channel_id"`
	LatestM3u8FileId         string `json:"latest_m3u8_file_id"`
	ConcatLatestFileId       string `json:"concat_latest_file_id"`
	ConcatM3u8Url            string `json:"concat_m3u8_url"`
	TranscodeState           uint8  `json:"transcode_state"`
	TranscodeSuccessLastTime string `json:"transcode_success_last_time"`
	ConcatSuccessLastTime    string `json:"concat_success_last_time"`
	TranscodeM3u8Url         string `json:"transcode_m3u8_url"`
	ConcatTimes              uint8  `json:"concat_times"`
	TranscodeTimes           uint8  `json:"transcode_times"`
	ComposeLatestFileId      string `json:"compose_latest_file_id"`
	ConcatMp4Url             string `json:"concat_mp4_url"`
	IsUseConcatMp4           uint8  `json:"is_use_concat_mp4"`
	IsDrm                    uint8  `json:"is_drm"`
	DrmM3u8Url               string `json:"drm_m3u8_url"`
}

type TaskGoodsDetail struct {
	Model

	AppId        string `json:"app_id"`
	AliveId      string `json:"alive_id"`
	ResourceId   string `json:"resource_id"`
	ResourceType int    `json:"resource_type"`
	ViewCount    int    `json:"view_count"`
	State        int    `json:"state"`
}

type CourseWareRecords struct {
	AppId              string              `json:"app_id"`
	AliveId            string              `json:"alive_id"`
	AliveTime          int                 `json:"alive_time"`
	CourseUseTime      int                 `json:"course_use_time"`
	UserId             json.JSONNullString `json:"user_id"`
	CoursewareId       json.JSONNullString `json:"courseware_id"`
	CurrentPreviewPage int                 `json:"current_preview_page"`
	CurrentImageUrl    json.JSONNullString `json:"current_image_url"`
	CutFileId          int                 `json:"cut_file_id"`
	IsShow             uint8               `json:"is_show"`
	CreatedAt          json.JSONTime       `json:"created_at"`
	UpdatedAt          json.JSONTime       `json:"updated_at"`
}

type CourseWare struct {
	Id                 json.JSONNullString `json:"id"`
	AppId              string              `json:"app_id"`
	AliveId            string              `json:"alive_id"`
	FileName           string              `json:"filename"`
	FileUrl            string              `json:"file_url"`
	FileType           uint8               `json:"file_type"`
	UseState           uint8               `json:"use_state"`
	ChangeToImageState uint8               `json:"change_to_image_state"`
	PageCount          int                 `json:"page_count"`
	State              int                 `json:"state"`
	CurrentPreviewPage int                 `json:"current_preview_page"`
	CoursewareImage    string              `json:"courseware_image"`
	CourseImageArray   []map[string]interface{}
}

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
	err := db.Select(s).Where("channel_id=?", channelId).First(&a).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &a, nil
}

// 获取直播讲师信息详情
func GetAliveRole(appId string, aliveId string) ([]*AliveRole, error) {
	var ar []*AliveRole
	// err := db.Select("role_name,user_id,user_name,is_current_lecturer,is_can_exceptional").
	// 	Where("app_id=? and alive_id=? and state=?", appId, aliveId, 0).
	// 	Find(&ar).Error
	db := db.Select("role_name,user_id,user_name,is_current_lecturer,is_can_exceptional").
		Where("app_id=? and alive_id=? and state=?", appId, aliveId, 0).
		Find(&ar)
	// 内存泄露风险
	// defer db.Close()
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return nil, db.Error
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

func GetAliveLookBackFile(appId string, aliveId string, s []string) (*AliveLookBack, error) {
	var alf AliveLookBack
	err := db.Table("t_alive_lookback").Select(s).Where("app_id=? and alive_id=? and state=?", appId, aliveId, 1).First(&alf).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &alf, nil
}

func GetAliveHlsResult(channelId string, s []string) (*AliveConcatHlsResult, error) {
	var ahr AliveConcatHlsResult

	err := db.Select(s).Where("channel_id=?", channelId).First(&ahr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &ahr, nil
}

// 更新直播观看人数
func UpdateViewCount(appId, aliveId string, viewCount int) error {
	_, err := GetAliveInfo(appId, aliveId, []string{"view_count"})
	if err != nil {
		return err
	}
	//todo 上报ES

	err = db.Where("app_id=? and id=?", appId, aliveId).Update("view_count", viewCount).Limit(1).Error
	if err != nil {
		return err
	}
	//todo 上报ES

	return nil
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

// 获取课件详情（通过courseWareId）
func GetCourseWareInfo(appId, aliveId string, courseWareId, s []string) ([]*CourseWare, error) {
	var cw []*CourseWare

	courseWareIds := strings.Join(courseWareId, ",") //用逗号,拼接

	err := db.Table("t_courseware").
		Select(s).
		Where("app_id=? and alive_id=? and id in (?)", appId, aliveId, courseWareIds).
		Find(&cw).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cw, nil
}

// 获取课件详情（通过aliveId）
func GetCourseWareInfoByAliveId(appId, aliveId string, s []string) (*CourseWare, error) {
	var cw CourseWare

	err := db.Table("t_courseware").
		Select(s).
		Where("app_id=? and alive_id=?", appId, aliveId).
		First(&cw).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &cw, nil
}

// 通过直播时间获取课件记录
func GetCourseWareByAliveTime(appId, aliveId string, lookBackId, aliveTime, pageSize int, ascOrder bool) ([]*CourseWareRecords, error) {
	var cwr []*CourseWareRecords

	//字段
	var s = []string{"app_id", "alive_id", "alive_time", "course_use_time", "user_id", "current_preview_page", "current_image_url", "courseware_id", "created_at", "updated_at"}

	where := "app_id = ? and alive_id = ? and is_show = ? and cut_file_id = ? and user_id is not null"

	conn := db.Table("t_alive_course_records").Select(s)

	if ascOrder {
		where += " and course_use_time >= ?"
		err := conn.Where(where, appId, aliveId, 1, lookBackId, aliveTime).
			Order("course_use_time ASC").
			Order("id ASC").
			Limit(pageSize).
			Find(&cwr).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	} else { //更加course_use_time从0到对应值内的课件
		where += " and courseware_id is not null and course_use_time >= ? and course_use_time <= ?"
		err := conn.Where(where, appId, aliveId, 1, lookBackId, 0, aliveTime).
			Order("course_use_time DESC").
			Order("id DESC").
			Limit(pageSize).
			Find(&cwr).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return cwr, nil
}

// 全部更新为不使用状态
func UpdateCourseWareAllNotUseState(appId, aliveId string) error {
	data := make(map[string]string)
	data["use_state"] = "0"
	data["updated_at"] = time.Now().Format("2006-01-02 15:06:04")

	err := db.Table("t_courseware").Where("app_id=? and alive_id=? and use_state=?", appId, aliveId, 1).
		Update(data).Error
	if err != nil {
		return err
	}

	return nil
}

// 更新课件使用状态
func UpdateCourseWareUseState(appId, aliveId, courseWareId string, data map[string]string) error {
	err := db.Table("t_courseware").
		Where("app_id=? and alive_id=? and id=?", appId, aliveId, courseWareId).
		Update(data).Error
	if err != nil {
		return err
	}

	return nil
}
