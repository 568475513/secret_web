package course

import (
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"abs/models/user"
	"abs/pkg/cache/alive_static"
	"abs/pkg/cache/redis_gray"
	"abs/pkg/logging"
)

type StaticData struct {
	IsFree           int    `redis:"is_free"`
	Title            string `redis:"title"`
	RoomId           string `redis:"room_id"`
	ImgUrl           string `redis:"img_url"`
	AliveType        int    `redis:"alive_type"`
	ForbidTalk       int    `redis:"forbid_talk"`
	Summary          string `redis:"summary"`
	ZbStartAt        string `redis:"zb_start_at"`
	RoomUrl          string `redis:"room_url"`
	ImgUrlCompressed string `redis:"img_url_compressed"`
	AliveVideoUrl    string `redis:"alive_video_url"`
	CommentCount     int    `redis:"comment_count"`
	PptImgs          string `redis:"ppt_imgs"`
	OrgContent       string `redis:"org_content"`
	AliveImgUrl      string `redis:"alive_img_url"`
	ViewCount        int    `redis:"view_count"`
	HavePassword     int    `redis:"have_password"`
	ZbStopAt         string `redis:"zb_stop_at"`
	AliveroomImgUrl  string `redis:"aliveroom_img_url"`
	ManualStopAt     string `redis:"manual_stop_at"`
	PaymentType      string `redis:"payment_type"`
	Descrb           string `redis:"descrb"`
	VideoAliveUseCos string `redis:"video_alive_use_cos"`
}

type AliveStatic struct {
	AppId     string
	AliveId   string
	UserId    string
	Type      string
	ExtraData string
}

const (
	staticAliveSwitch        = "aliveGw:switch"
	staticAlivePageCacheList = "aliveGw:page:list"
	hashStaticAliveUser      = "hash_static_alive_user_"
	currentDayAliveRole      = "current_day_alive_role:%s:%s"
	currentDayAliveInfo      = "current_day_alive_info:%s:%s"
)

//直播静态化切换主流程
func (c *AliveStatic) AliveStaticMain(agentType int) (RoomData map[string]interface{},staticSwitch bool, err error) {

	StaticRedisCon, err := alive_static.GetStaticRedisCon()
	if err != nil {
		logging.Error(err)
	}
	defer StaticRedisCon.Close()
	RoomData = make(map[string]interface{})

	if staticSwitch =  c.CheckAliveStaticSwitch(StaticRedisCon);staticSwitch {
		staticDataValues, err := redis.Values(StaticRedisCon.Do("HGETALL", fmt.Sprintf(currentDayAliveInfo, c.AppId, c.AliveId))) //获取直播静态数据
		staticData := &StaticData{}
		if err := redis.ScanStruct(staticDataValues, staticData); err != nil {
			logging.Error(err)
		}
		var (
			yymmdd         = time.Now().Format("2006-01-02")
			key            = c.AliveId + c.UserId
			ava, _         = redis.Values(StaticRedisCon.Do("HMGET", hashStaticAliveUser+yymmdd, key))
			teacher_val, _ = redis.Values(StaticRedisCon.Do("HGETALL", fmt.Sprintf(currentDayAliveRole, c.AppId, c.AliveId)))
		)
		teacher := ""
		for _, v := range teacher_val {
			teacher += string(v.([]byte))
		}
		is_teacher := strings.Index(teacher, c.UserId) > -1
		if c.Type == "12" {
			if (staticData.IsFree == 1 && !is_teacher) || (staticData.Title != "" && staticData.IsFree != 1 && len(ava) > 0 && !is_teacher) {
				RoomData["alive_info"] = map[string]interface{}{
					"app_id":   c.AppId,
					"alive_id": c.AliveId,
					"room_id":  staticData.RoomId,
					// 直播间标题
					"title": staticData.Title,
					// 直播间描述
					"descrb": staticData.Descrb,
					// 直播间简介
					"summary": staticData.Summary,
					// 直播专栏名称
					"product_name": "",
					// 直播专栏ID
					"product_id": "",
					// 直播封面或者暖场图
					"img_url": staticData.ImgUrl,
					// 首页展示的直播图
					"alive_img_url": staticData.AliveImgUrl,
					// 直播封面图
					"cover_img_url": staticData.ImgUrl,
					// 直播压缩图
					"img_url_compressed": staticData.ImgUrlCompressed,
					// 直播类型（语音/视频）0-语音直播，1-视频直播 2-推流直播
					"alive_type": staticData.AliveType,
					// 获取直播状态
					"alive_state": 1,
					// 推流状态，0推流结束，1推流中，2推流未开始
					"push_state": "",
					// 直播剩余时长（秒）
					"remainder_time": "",
					// 推流直播开始时间
					"pushzb_start_at": staticData.ZbStartAt,
					// 推流直播结束时间
					"pushzb_stop_at": staticData.ZbStopAt,
					// 直播开始时间（时间戳：秒）
					"zb_start_at": staticData.ZbStartAt,
					// 直播结束时间（时间戳：秒）
					"zb_stop_at":     staticData.ZbStopAt,
					"checktimestamp": "",
					"manual_stop_at": staticData.ManualStopAt,
					"view_count":     staticData.ViewCount,
					"comment_count":  staticData.CommentCount,
					"push_url":       "",
					"push_ahead":     "",
					"can_select":     "",
					"org_content":    staticData.OrgContent,

					// 用户类型学员、讲师
					"user_type": 0,
					"user_id":   c.UserId,
				}
				RoomData["available_info"] = map[string]interface{}{
					"available":         true,
					"available_product": true,
					"expire_at":         "",
					"have_password":     "",
					"is_public":         "",
					"is_stop_sell":      "",
					"is_try":            "",
					"payment_type":      1,
					"recycle_bin_state": "",
				}
				isGrayBool := redis_gray.InGrayShop("video_alive_not_use_cos", c.AppId)
				// 不为小程序--不在O端名单内
				RoomData["alive_play"] = map[string]interface{}{
					"alive_video_url":     staticData.AliveVideoUrl,
					"video_alive_use_cos": false,
				}
				if !isGrayBool && staticData.VideoAliveUseCos == "1" && agentType != 14 {
					RoomData["alive_play"] = map[string]interface{}{
						"video_alive_use_cos": true,
						"new_alive_video_url": staticData.AliveVideoUrl,
						"alive_video_url":     staticData.AliveVideoUrl,
					}
				}
				RoomData["alive_conf"] = map[string]interface{}{
					"alive_mode":       0,
					"show_on_wall":     1,
					"alive_type_state": 1,
				}
				RoomData["share_info"] = ""
				RoomData["caption_define"] = ""
				RoomData["index_url"] = ""
				RoomData["is_static_switch"] = true

				return RoomData,staticSwitch, err
			}
		}

		if (staticData.IsFree == 1 && !is_teacher) || (staticData.Title != "" && staticData.IsFree != 1 && len(ava) > 0 && !is_teacher) {
			RoomData["alive_info"] = map[string]interface{}{
				"image":       staticData.ImgUrl,
				"title":       staticData.Title,
				"summary":     staticData.Summary,
				"start_time":  staticData.ZbStartAt,
				"org_content": "", //详情会出现尖括号，不好处理
				"alive_room":  staticData.RoomUrl,
			}
			RoomData["available_info"] = ""
			RoomData["alive_play"] = ""
			RoomData["alive_conf"] = ""
			RoomData["share_info"] = ""
			RoomData["caption_define"] = ""
			RoomData["index_url"] = ""
			RoomData["is_static_switch"] = true

			return RoomData,staticSwitch, err
		}
	}
	return RoomData,staticSwitch, err
}

//检查是否开启直播静态化开关
func (c *AliveStatic) CheckAliveStaticSwitch(conn redis.Conn) (Switch bool) {

	var (
		StaticSwitch    bool
		StaticPageCache bool
	)
	//判断是否开启静态化开关
	StaticSwitchNum, err := redis.Int(conn.Do("GET", staticAliveSwitch)) //alive_static.GetStaticSwitch(conn, staticAliveSwitch)
	if err != nil {
		logging.Error(err)
	}
	if StaticSwitchNum > 0 {
		StaticSwitch = true
	} else {
		StaticSwitch = false
	}
	//判断页面是否有缓存
	page_strings := fmt.Sprintf("%s:%s", c.AppId, c.AliveId)
	StaticPageCache, err = redis.Bool(conn.Do("SISMEMBER", staticAlivePageCacheList, page_strings))
	if err != nil {
		logging.Error(err)
	}

	if StaticSwitch || StaticPageCache {
		return true
	}

	return false
}

//次级业务接口静态化逻辑
func (c *AliveStatic) SecondaryInfoStaticData(im map[string]string, user user.User) map[string]interface{} {

	data := make(map[string]interface{})
	// 组装用户信息
	userInfoMap := make(map[string]interface{})
	userInfoMap["phone"] = user.Phone
	userInfoMap["wx_avatar"] = user.WxAvatar
	userInfoMap["wx_nickname"] = user.WxNickname
	// 用户信息
	data["user_info"] = userInfoMap
	// 短信预约总开关
	data["is_message_on"] = ""
	// 用户是否被禁言
	data["is_show"] = 1
	// 用户黑名单
	data["black_list"] = ""
	// 邀请卡链接
	data["invite_card_url"] = ""
	// 邀请讲师链接
	data["invite_teacher_url"] = ""
	// 邀请达人榜链接
	data["invite_list_url"] = ""
	// 共享文件列表链接
	data["share_file_url"] = ""
	// 获取云通信配置
	data["im_init"] = im
	return data
}
