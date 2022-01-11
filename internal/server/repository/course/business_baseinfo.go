package course

import (
	"abs/pkg/cache/redis_im"
	"abs/pkg/cache/redis_xiaoe_im"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"abs/internal/server/rules/validator"
	"abs/models/alive"
	"abs/models/business"
	"abs/pkg/cache/alive_static"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/cache/redis_gray"
	e "abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"
)

// 业务类型的类
type BaseInfo struct {
	Alive    *alive.Alive
	AliveRep *AliveInfo
	UserType uint
}

// 快直播结构体
type LiveUrl struct {
	PcAliveVideoUrl           string                   `json:"pc_alive_video_url"`            //pc播放地址
	MiniAliveVideoUrl         string                   `json:"mini_alive_video_url"`          //小程序播放地址
	AliveVideoUrl             string                   `json:"alive_video_url"`               //上传视频播放地址
	AliveFastWebrtcurl        string                   `json:"alive_fast_webrtcurl"`          //快直播播放地址
	NewAliveVideoUrl          string                   `json:"new_alive_video_url"`           //录播新方式的播放链接
	FastAliveSwitch           bool                     `json:"fast_alive_switch"`             //快直播开关
	VideoAliveUseCos          bool                     `json:"video_alive_use_cos"`           //置为使用cos录播方式
	AliveVideoMoreSharpness   []map[string]interface{} `json:"alive_video_more_sharpness"`    //普通直播多清晰度
	PcAliveVideoMoreSharpness []map[string]interface{} `json:"pc_alive_video_more_sharpness"` //pc普通直播多清晰度
	AliveFastMoreSharpness    []map[string]interface{} `json:"alive_fast_more_sharpness"`     //快直播多清晰度
	ReCordedUsePullStream     bool                     `json:"recorded_use_pull_stream"`      //录播直播是否伪直播
}

const (
	aliveOnlineUV = "POMELO:USER_LIST:%s"
)

// 组装直播的基本信息
func (b *BaseInfo) GetAliveInfoDetail() map[string]interface{} {
	aliveInfoDetail, now := make(map[string]interface{}), time.Now()
	aliveInfoDetail["app_id"] = b.AliveRep.AppId
	aliveInfoDetail["alive_id"] = b.Alive.Id
	aliveInfoDetail["room_id"] = b.Alive.RoomId
	aliveInfoDetail["room_id"] = b.AliveRep.GetAliveRommId(b.Alive)
	aliveInfoDetail["resource_type"] = e.ResourceTypeLive
	// 直播间标题
	aliveInfoDetail["title"] = b.Alive.Title.String
	// 直播间描述
	aliveInfoDetail["descrb"] = b.Alive.Descrb
	// 直播间简介
	aliveInfoDetail["summary"] = b.Alive.Summary.String
	// 直播专栏名称
	aliveInfoDetail["product_name"] = b.Alive.ProductName.String
	// 直播专栏ID
	aliveInfoDetail["product_id"] = b.Alive.ProductId.String
	// 直播封面或者暖场图
	aliveInfoDetail["img_url"] = b.Alive.ImgUrl.String
	if b.Alive.AliveroomImgUrl.String != "" {
		aliveInfoDetail["img_url"] = b.Alive.AliveroomImgUrl.String
	}
	// 首页展示的直播图
	aliveInfoDetail["alive_img_url"] = b.Alive.AliveImgUrl.String
	// 直播封面图
	aliveInfoDetail["cover_img_url"] = b.Alive.ImgUrl.String
	// 直播压缩图
	aliveInfoDetail["img_url_compressed"] = b.Alive.ImgUrlCompressed.String
	// 直播类型（语音/视频）0-语音直播，1-视频直播 2-推流直播
	aliveInfoDetail["alive_type"] = b.Alive.AliveType
	// 推流直播开始时间
	aliveInfoDetail["pushzb_start_at"] = b.Alive.ZbStartAt
	// 推流直播结束时间
	aliveInfoDetail["pushzb_stop_at"] = b.Alive.ZbStopAt
	// 获取直播状态
	aliveInfoDetail["alive_state"] = b.AliveRep.GetAliveStates(b.Alive)
	// 推流状态，0推流结束，1推流中，2推流未开始
	aliveInfoDetail["push_state"] = b.Alive.PushState
	// 直播剩余时长（秒）
	aliveInfoDetail["remainder_time"] = 0
	if b.Alive.AliveType == e.AliveTypeVideo {
		aliveInfoDetail["remainder_time"] = b.Alive.ZbStartAt.Unix() + b.Alive.VideoLength - now.Unix()
	}
	// 直播开始时间（时间戳：秒）
	aliveInfoDetail["zb_start_at"] = b.Alive.ZbStopAt.Unix()
	// 直播结束时间（时间戳：秒）
	aliveInfoDetail["zb_stop_at"] = b.Alive.ZbStopAt.Unix()
	// 距离直播开始倒计时（单位：秒）
	aliveInfoDetail["zb_countdown_time"] = b.Alive.ZbStartAt.Unix() - now.Unix()
	// 系统时间 进入直播间时的时间（时间戳：秒）这个是重复的 => aliveInfoDetail["enter_time"]
	aliveInfoDetail["checktimestamp"] = now.Unix()
	aliveInfoDetail["manual_stop_at"] = b.Alive.ManualStopAt
	aliveInfoDetail["view_count"] = b.Alive.ViewCount
	aliveInfoDetail["comment_count"] = b.Alive.CommentCount
	// 只有讲师才需要push_url
	if b.UserType == 1 {
		aliveInfoDetail["push_url"] = b.Alive.PushUrl
	} else {
		aliveInfoDetail["push_url"] = ""
	}
	aliveInfoDetail["push_ahead"] = b.Alive.PushAhead
	aliveInfoDetail["can_select"] = b.Alive.CanSelect
	aliveInfoDetail["org_content"] = b.Alive.OrgContent
	// 人次隐藏逻辑（设置后只有讲师看到人次）
	if b.Alive.ConfigShowViewCount == 1 && b.UserType == 0 {
		aliveInfoDetail["comment_count"] = 0
		aliveInfoDetail["view_count"] = 0
	}
	// 用户类型学员、讲师
	aliveInfoDetail["user_type"] = b.UserType
	// 参数加密串
	aliveInfoDetail["param_str"], _ = util.PutParmToStr(map[string]interface{}{
		"payment_type":  e.PaymentTypeSingle,
		"resource_type": e.ResourceTypeLive,
		"resource_id":   b.Alive.Id,
		"product_id":    "",
	})
	//拼接直播间链接
	aliveInfoDetail["alive_room_url"] = util.GetNewAliveRoom(b.Alive.AppId, b.Alive.Id, strconv.Itoa(e.PaymentTypeSingle), b.Alive.ProductId.String)

	// 录播底层优化新增 - 预期视频或推流的结束时间
	aliveInfoDetail["record_push_end_time"] = b.Alive.ZbStartAt.Time.Add(time.Duration(b.Alive.VideoLength) * 1e9).Format("2006-01-02 15:04:05")
	return aliveInfoDetail
}

// 组装直播的权益模块
func (b *BaseInfo) GetAvailableInfo(available, availableProduct bool, expireAt string) map[string]interface{} {
	availableInfo := make(map[string]interface{})
	availableInfo["available"] = available
	availableInfo["available_product"] = availableProduct
	availableInfo["expire_at"] = expireAt
	availableInfo["payment_type"] = b.Alive.PaymentType
	availableInfo["recycle_bin_state"] = b.Alive.RecycleBinState
	availableInfo["is_stop_sell"] = b.Alive.IsStopSell
	availableInfo["have_password"] = b.Alive.HavePassword
	availableInfo["is_try"] = b.Alive.IsTry
	availableInfo["is_public"] = b.Alive.IsPublic
	return availableInfo
}

// 组装直播店铺配置信息
func (b *BaseInfo) GetAliveConfInfo(baseConf *service.AppBaseConf, aliveModule *alive.AliveModuleConf, available bool, userId string) map[string]interface{} {
	aliveConf := make(map[string]interface{})
	// 店铺名称
	aliveConf["wx_app_name"] = baseConf.ShopName
	// 店铺头像
	aliveConf["wx_app_avatar"] = baseConf.ShopLogo
	// 是否全局禁言
	aliveConf["forbid_talk"] = b.Alive.ForbidTalk
	// 用户发言上墙开关
	aliveConf["show_on_wall"] = b.Alive.ShowOnWall
	// 用户上麦开关
	aliveConf["can_record"] = b.Alive.CanRecord
	// 是否有打赏功能
	aliveConf["has_reward"] = baseConf.HasReward
	// 打赏提醒是否显示 0-全部显示 1-仅讲师和打赏者可见
	aliveConf["is_show_reward"] = b.Alive.ConfigShowReward
	// 共享文件插件开关 0-不可用  1-可用
	aliveConf["share_file_switch"] = 0
	// 学员上墙开关 0-不可用  1-可用
	aliveConf["show_on_wall_switch"] = 0
	// 店铺是否认证
	aliveConf["authentic_state"] = baseConf.AuthenticState
	// 打赏功能是否可用 0-不可用 1-可用
	aliveConf["reward_switch"] = baseConf.HasReward
	// 邀请功能是否开启 0-不可用 1-可用
	aliveConf["has_invite"] = baseConf.HasInvite
	//店铺名称为空 conf_hub服务异常,邀请功能置 1 (默认值)
	if baseConf.ShopName == "" {
		aliveConf["has_invite"] = 1
	}
	// 版本信息
	aliveConf["version_type"] = baseConf.VersionType
	// 是否显示关联售卖界面，默认1-显示，0-不显示
	aliveConf["relate_sell_info"] = baseConf.RelateSellInfo
	// 是否开启回看
	aliveConf["is_lookback"] = b.Alive.IsLookback
	// 是否开启直播带货
	aliveConf["is_takegoods"] = b.Alive.IsTakegoods
	// 是否只有h5观看直播
	aliveConf["only_h5_play"] = baseConf.OnlyH5Play
	// 模板消息推送状态
	aliveConf["if_push"] = b.Alive.IfPush
	// 播放器 0-默认播放器 1-自研播放器 (默认0)
	aliveConf["video_player_type"] = baseConf.VideoPlayerType
	// 隐私保护开关 0-关闭 1-开启
	aliveConf["is_privacy_protection"] = baseConf.IsPrivacyProtection
	// 标准版是否可用 0-不可用  1-可用
	versionUse := 1
	// 标准版或试用版过期
	if versionState := util.JudgeDate(baseConf.VersionType, baseConf.ExpireTime); versionState["type"] == 1 || versionState["type"] == 3 {
		versionUse = 0
	}
	//开启PC网校 0为关闭 1为开启
	if baseConf.IsEnable == 1 && baseConf.IsValid == 1 {
		aliveConf["open_pc_network_school"] = 1
	} else {
		aliveConf["open_pc_network_school"] = 0
	}

	//PC网校自定义域名
	aliveConf["pc_network_school_index_url"] = baseConf.PcCustomDomain
	aliveConf["is_open_promoter"] = aliveModule.IsOpenPromoter
	// 版本过期信息
	versionState := b.GetAppExpireTime(baseConf.Profit)
	aliveIsRemind := 0
	// 版本是否过期和功能是否过期0-过期 1-没过期，2-流量余额为0，-1未购买
	// 该类型直播间是否过期
	switch b.Alive.AliveType {
	case 0:
		aliveIsRemind = versionState["alive_voice_is_remind"].(int)
	case 1:
		aliveIsRemind = versionState["alive_video_voice_is_remind"].(int)
	case 2:
		aliveIsRemind = versionState["alive_push_is_remind"].(int)
	case 3:
		aliveIsRemind = versionState["alive_PPT_is_remind"].(int)
	}
	// 版本是否过期和功能是否过期0-未过期 1-即将过期，2-过期，-1未购买
	// 推流直播暂不加权限限制
	// 该类型直播间是否可用开关 0-不可用  1-可用
	if b.Alive.AliveType == e.AliveTypePush {
		aliveConf["alive_type_state"] = 1
	} else if (aliveIsRemind == 1 || aliveIsRemind == 0) && versionUse == 1 {
		aliveConf["alive_type_state"] = 1
	} else {
		aliveConf["alive_type_state"] = 0
	}
	// 共享文件功能是否过期
	if versionState["alive_share_file_is_remind"].(int) == 1 || versionState["alive_share_file_is_remind"].(int) == 0 {
		aliveConf["share_file_switch"] = 1
	} else if versionState["alive_share_file_is_remind"] == 2 {
		aliveConf["share_file_switch"] = 0
	} else {
		aliveConf["share_file_switch"] = 0
	}
	// 学员上墙功能是否过期
	if versionState["student_show_is_remind"].(int) == 1 || versionState["student_show_is_remind"].(int) == 0 {
		aliveConf["show_on_wall_switch"] = 1
	}
	// 打赏功能是否过期
	if versionState["alive_reward_is_remind"].(int) == 2 || versionState["alive_reward_is_remind"].(int) == -1 {
		aliveConf["reward_switch"] = 0
	}
	// 隐藏直播间人次功能是否过期
	// if versionState["alive_show_man_time_is_remind"].(int) == 1 || versionState["alive_show_man_time_is_remind"].(int) == 0 {
	// 	aliveConf["is_show_view_count_switch"] = 1
	// }

	// 是否显示直播人次
	aliveConf["is_show_view_count"] = 1
	if (versionState["alive_show_man_time_is_remind"].(int) == 1 || versionState["alive_show_man_time_is_remind"].(int) == 0) && b.Alive.ConfigShowViewCount == 1 && b.UserType == 0 {
		aliveConf["is_show_view_count"] = 0
	}
	/**
	 * @Description: 隐藏人次服务过期，重新显示人次
	 * @param versionState["alive_show_man_time_is_remind"].(int)
	 * @return {
	 */
	if versionState["alive_show_man_time_is_remind"].(int) != 1 && versionState["alive_show_man_time_is_remind"].(int) != 0 && b.Alive.ConfigShowViewCount == 1 && b.UserType == 0 {
		aliveConf["is_show_view_count"] = 0
	}

	// 获取直播配置表相关配置
	// 邀请达人榜需要灰度控制
	if redis_gray.InGrayShopNew("invite_forbid", b.AliveRep.AppId) {
		aliveConf["is_invite_on"] = 0
	} else {
		aliveConf["is_invite_on"] = aliveModule.IsInviteOn
	}
	aliveConf["is_message_on"] = aliveModule.IsMessageOn
	aliveConf["alive_mode"] = aliveModule.AliveMode
	aliveConf["is_picture_on"] = aliveModule.IsPictureOn
	aliveConf["is_audit_first_on"] = aliveModule.IsAuditFirstOn
	aliveConf["is_coupon_on"] = aliveModule.IsCouponOn
	aliveConf["is_card_on"] = aliveModule.IsCardOn
	aliveConf["is_prize_on"] = aliveModule.IsPrizeOn
	aliveConf["is_redirect_index"] = aliveModule.IsRedirectIndex
	aliveConf["complete_time"] = aliveModule.CompleteTime
	// 是否开启打赏， 0-关闭 1-开启
	aliveConf["is_show_reward_on"] = aliveModule.IsShowRewardOn
	// 是否开启签到，0-未开启，1-开启
	aliveConf["is_sign_in_on"] = aliveModule.IsSignInOn
	// 红包功能是否开启，0-关闭，1-开启
	aliveConf["is_red_packet_on"] = aliveModule.IsRedPacketOn

	if aliveModule.CompleteTime == 0 {
		aliveConf["is_open_complete_time"] = 0
	} else {
		aliveConf["is_open_complete_time"] = 1
	}
	//该直播是否开启圆桌会议模式，0关闭，1开启
	aliveConf["is_round_table_on"] = aliveModule.IsRoundTableOn

	//是否开启防录屏
	aliveConf["anti_screen_jump"] = 0
	aliveConf["anti_screen_jump_url"] = ""
	if aliveModule.IsAntiScreen == 1 && b.UserType == 0 && available == true {
		aliveConf["anti_screen_jump"] = 1
		aliveConf["anti_screen_jump_url"] = os.Getenv("APP_REDIRECT_DOMAIN") + "open_app?app_id=" + b.AliveRep.AppId + "&params=" + b.GetWakeUpAppParams(userId)
	}

	tab := &alive.AliveTab{}
	//该直播是否自定义tab
	if err := util.JsonDecode([]byte(aliveModule.AliveJson), tab); err != nil || tab.TabOn == "0" {
		aliveConf["alive_tab"] = 0
	} else {
		aliveConf["alive_tab"] = 1
	}

	return aliveConf
}

// 获取直播间相关的链接
func (b *BaseInfo) GetAliveLiveUrl(agentType, version, enableWebRtc int, UserId string) (liveUrl LiveUrl) {
	var (
		playUrls     []string
		err          error
		isUserWebRtc bool
		// isEnableWebRtc bool
	)

	timeStamp := time.Now().Unix()

	supportSharpness := map[string]interface{}{
		"default": "原画", //默认原画
		"hd":  "高清", //高清（720P）
		"fluent":  "流畅", //流畅（480P）
	}

	if err = util.JsonDecode([]byte(b.Alive.PlayUrl), &playUrls); err != nil {
		logging.Warn(fmt.Sprintf("获取直播间播放链接JsonDecode有错误【非致命，不慌】：%s", err.Error()))
		// 不能返回，有特殊的PlayUrl
		// return
	}
	if len(playUrls) >= 3 && (b.Alive.AliveType == e.AliveTypePush || b.Alive.AliveType == e.AliveOldTypePush) {
		liveUrl.PcAliveVideoUrl = playUrls[1]
		liveUrl.AliveVideoUrl, liveUrl.MiniAliveVideoUrl = playUrls[2], playUrls[2]
		// 店铺设置是否开启快直播【无用了现在】
		// isEnableWebRtc = b.canUseFastLive(version)
		// 快直播功能判断
		// 用户是否可用快直播
		if isUserWebRtc, err = b.isUseFastLive(UserId); err != nil {
			logging.Error(fmt.Sprintf("获取用户是否可用快直播错误：%s", err.Error()))
			// 这里需要返回吗？
			// return
		}

		currentUv := b.getCurrentUv()

		// 普通直播多清晰度
		liveUrl.AliveVideoMoreSharpness = make([]map[string]interface{}, len(supportSharpness))
		liveUrl.PcAliveVideoMoreSharpness = make([]map[string]interface{}, len(supportSharpness))
		i := 0
		for k, v := range supportSharpness {
			i = b.getIndex(currentUv, k)

			liveUrl.AliveVideoMoreSharpness[i] = map[string]interface{}{
				"definition_name": v,
				"definition_p":    k,
				"url":             b.getPlayUrlBySharpness(k, playUrls[2], b.Alive.ChannelId),
				"encrypt":         "",
			}
			liveUrl.PcAliveVideoMoreSharpness[i] = map[string]interface{}{
				"definition_name": v,
				"definition_p":    k,
				"url":             b.getPlayUrlBySharpness(k, playUrls[1], b.Alive.ChannelId),
				"encrypt":         "",
			}
		}

		// 快直播O端名单目录
		isGray := redis_gray.InGrayShop("fast_alive_switch", b.AliveRep.AppId)
		if isGray && isUserWebRtc && enableWebRtc == 1 && util.Substr(playUrls[0], 0, 4) == "rtmp" {
			limitUv, _ := strconv.Atoi(os.Getenv("WEBRTC_SWITCH_RTMP_UV"))

			//成本控制的白名单
			inCostOptWhiteMenu := redis_gray.InGrayShopSpecialHit("webrtc_cost_opt_white_menu", b.Alive.AppId)

			if inCostOptWhiteMenu || currentUv < limitUv {
				liveUrl.AliveFastWebrtcurl = "webrtc" + util.Substr(playUrls[0], 4, len(playUrls[0]))
				liveUrl.FastAliveSwitch = true
				//快直播多清晰度
				liveUrl.AliveFastMoreSharpness = make([]map[string]interface{}, len(supportSharpness))
				i := 0
				for k, v := range supportSharpness {
					i = b.getIndex(currentUv, k)

					liveUrl.AliveFastMoreSharpness[i] = map[string]interface{}{
						"definition_name": v,
						"definition_p":    k,
						"url":             b.getPlayUrlBySharpness(k, liveUrl.AliveFastWebrtcurl, b.Alive.ChannelId),
						"encrypt":         "",
					}
				}
			} else {
				// 触发成本控制了，记录下
				logging.Info(fmt.Sprintf("cost_optimization app_id:%s alive_id:%s uv:%d limit:%d",
					b.Alive.AppId, b.Alive.Id, currentUv, limitUv))
			}
		}
	} else {
		liveUrl.MiniAliveVideoUrl = fmt.Sprintf("https://%s/%s.m3u8?%d", util.GetH5Domain(b.AliveRep.AppId, true), b.AliveRep.AliveId, timeStamp)
		liveUrl.AliveVideoUrl = liveUrl.MiniAliveVideoUrl
		// todo 录播底层优化-直播链接下发逻辑 start
		/**
		 * 先通过是否存在channel_id判断是否走伪直播逻辑
		 * 通过直播状态值判断是否存在转推任务,因为推流了这边会由记录为1
		 * redis也没有那么及时key 设置为 alive_recorded_{channel_id}
		 */
		isUsePullStream := b.GetNowRecordedIsPush()
		if isUsePullStream {
			liveUrl.AliveVideoMoreSharpness = make([]map[string]interface{}, 2)
			recordedUrl := "http://" + os.Getenv("LIVE_PLAY_HOST") + b.Alive.ChannelId + ".m3u8"
			liveUrl.AliveVideoMoreSharpness[0] = map[string]interface{}{
				"definition_name": "原画",
				"definition_p":    "default",
				"url":             b.getPlayUrlBySharpness("default", recordedUrl, b.Alive.ChannelId),
				"encrypt":         "",
			}
			liveUrl.AliveVideoMoreSharpness[1] = map[string]interface{}{
				"definition_name": "流畅",
				"definition_p":    "fluent",
				"url":             b.getPlayUrlBySharpness("fluent", recordedUrl, b.Alive.ChannelId),
				"encrypt":         "",
			}
		}
		liveUrl.ReCordedUsePullStream = isUsePullStream
		// todo 录播底层优化-直播链接下发逻辑 end
		isGrayBool := redis_gray.InGrayShop("video_alive_not_use_cos", b.AliveRep.AppId)
		// play_url不为空--不为小程序--不在O端名单内
		if !isGrayBool && b.Alive.PlayUrl != "" && agentType != 14 {
			// Es日志
			// logging.LogToEs("新录播方式log", map[string]interface{}{
			// 	"app_id": a.AppId,
			// 	"redis_gray": grayBool,
			// 	"playUrl": playUrl,
			// 	"agentType": agentType,
			// })
			// 置为使用cos录播方式
			liveUrl.VideoAliveUseCos = true
			if len(playUrls) > 3 {
				liveUrl.NewAliveVideoUrl = playUrls[3]
			} else {
				liveUrl.NewAliveVideoUrl = b.Alive.PlayUrl
			}
		}
	}
	return
}

// 获取播放链接的位置
func (b *BaseInfo) getIndex(currentUv int, k string) int{
	//获取超过多少UV默认使用【高清】播放的配置
	limitUvUseHd, _ := strconv.Atoi(os.Getenv("DEFAULT_USE_HD_LIMIT_UV"))
	//获取超过多少UV默认使用【流畅】播放的配置
	limitUvUseFluent, _ := strconv.Atoi(os.Getenv("DEFAULT_USE_FLUENT_LIMIT_UV"))

	//是否是默认使用【高清】播放的店铺
	inGrayDefaultUseHd := redis_gray.InGrayShopSpecialHit("alive_default_use_hd_switch", b.Alive.AppId)
	//是否是默认使用【流畅】播放的店铺
	inGrayDefaultUseFluent := redis_gray.InGrayShopSpecialHit("alive_default_use_fluent_switch", b.Alive.AppId)

	i := 0
	if inGrayDefaultUseFluent && currentUv > limitUvUseFluent {
		//默认使用流畅（0代表默认 default这个命名忽略 历史原因）
		switch k {
		case "hd":
			i = 1
		case "fluent":
			i = 0
		case "default":
			i = 2
		}

		if i == 0 {
			logging.Info(fmt.Sprintf("default_play_url:use_fluent,app_id:%s,alive_id:%s,current_uv:%d,limit_uv:%d", b.Alive.AppId,b.Alive.Id,currentUv,limitUvUseFluent))
		}
	}else if inGrayDefaultUseHd && currentUv > limitUvUseHd {
		//默认使用高清（0代表默认 default这个命名忽略 历史原因）
		switch k {
		case "hd":
			i = 0
		case "fluent":
			i = 1
		case "default":
			i = 2
		}

		if i == 0 {
			logging.Info(fmt.Sprintf("default_play_url:use_hd,app_id:%s,alive_id:%s,current_uv:%d,limit_uv:%d", b.Alive.AppId,b.Alive.Id,currentUv,limitUvUseHd))
		}
	} else {
		//默认使用原画（0代表默认 default这个命名忽略 历史原因）
		switch k {
		case "hd":
			i = 1
		case "fluent":
			i = 2
		case "default":
			i = 0
		}
	}

	return i
}

func (b *BaseInfo) getCurrentUv() int{
	xiaoEImRedisConn, err := redis_xiaoe_im.GetConn()
	if err != nil {
		logging.Error(err)
	}
	defer xiaoEImRedisConn.Close()

	currentUv := 0

	//查询实时在线UV
	cacheKey := fmt.Sprintf(aliveOnlineUV, b.Alive.Id)
	currentUv, err = redis.Int(xiaoEImRedisConn.Do("ZCARD", cacheKey))

	if err != nil {
		//这里只记录查询redis失败日志，不去影响主流程
		logging.Error(fmt.Sprintf("base_info 查询实时在线人数失败：%s", err.Error()))
	}

	return currentUv
}

// 直播静态页的信息采集【用户】
func (b *BaseInfo) SetAliveUserToStaticRedis(userId string) {

	err := alive_static.HsetNxString(fmt.Sprintf(staticAliveHashUser, time.Now().Format("2006-01-02")), b.Alive.Id+userId, 1, 3600*24)
	if err != nil {
		logging.Error(err)
	}
}

// 直播静态页的信息采集【ID】
func (b *BaseInfo) SetAliveIdToStaticRedis() {
	err := alive_static.HsetNxString(fmt.Sprintf(staticAliveHashId, time.Now().Format("2006-01-02")), b.Alive.Id, b.Alive.Id, 3600*24)
	if err != nil {
		logging.Error(err)
	}
}

// 页面跳转逻辑专用方法
func (b *BaseInfo) BaseInfoPageRedirect(
	products []*business.PayProducts,
	available bool,
	versionType int,
	req validator.BaseInfoRuleV2) (url string, code int, msg string) {
	// 当无自身与属于专栏/会员售卖形式时，将超级会员加入进来
	if len(products) == 0 && b.Alive.PaymentType == e.PaymentTypeProduct && !available {
		svipReq := Svip{AppId: b.Alive.AppId, ResourceId: b.Alive.Id, ResourceType: e.ResourceTypeLive}
		if redirect := svipReq.GetResourceSvipRedirectV2(); redirect != "" {
			url = redirect
			// 是否更多来源
			if req.MoreWay == "1" {
				url = url + "?more_way=1"
			}
			msg = "超级会员跳转"
			code = e.RESOURCE_REDIRECT
			return
		}
	}
	// 兼容企学院的跳转
	if util.IsQyApp(versionType) && !available {
		msg = "兼容企学院的跳转"
		url = "/training_page/noPermission"
		code = e.RESOURCE_REDIRECT
		return
	}
	// 判断是否购买,如果未购买则跳专栏
	if b.Alive.ProductId.String != "" {
		urlParams := util.ContentParam{
			Type:        e.PaymentTypeProduct,
			ProductId:   b.Alive.ProductId.String,
			ChannelId:   req.ChannelId,
			ShareUserId: req.ShareUserId,
			ShareType:   req.ShareType,
		}
		url = util.ContentUrl(urlParams)
		// 页面跳转（未购买且属于多专栏/会员）
		if available == false && b.Alive.PaymentType == e.PaymentTypeProduct {
			msg = "页面跳转（未购买且属于多专栏/会员）"
			code = e.RESOURCE_REDIRECT
			urlColumParams := util.ContentParam{
				ChannelId:  req.ChannelId,
				ShareAgent: req.ShareAgent,
				ShareFrom:  req.ShareFrom,
			}
			if len(products) > 1 {
				urlColumParams.ResourceId = b.Alive.Id
				urlColumParams.ResourceType = e.ResourceTypeLive
				urlColumParams.ShareUserId = req.ShareUserId
				urlColumParams.ShareType = req.ShareType
				url = util.ParentColumnsUrl(urlColumParams)
			} else if len(products) == 1 {
				urlColumParams.Type = e.PaymentTypeProduct
				urlColumParams.ResourceType = int(products[0].SrcType)
				urlColumParams.ProductId = products[0].Id
				urlColumParams.AppId = products[0].AppId
				urlColumParams.ResourceId = products[0].Id
				if req.ContentAppId != "" {
					urlColumParams.ContentAppId = req.ContentAppId
					urlColumParams.Source = "2"
				}
				url = util.ContentUrl(urlColumParams)
			} else {
				code = e.SUCCESS
			}
		}
	}

	return
}

// 获取防录屏落地页唤起APP参数
func (b *BaseInfo) GetWakeUpAppParams(userId string) (url string) {
	params := make(map[string]interface{})
	params["app_id"] = b.AliveRep.AppId
	params["resource_id"] = b.AliveRep.AliveId
	params["resource_type"] = 4
	params["user_id"] = userId
	params["content_app_id"] = ""
	params["encrypt_user_id"] = util.GetEncryptUserId(userId)
	tempParam := make(map[string]interface{})
	tempParam["params"] = params
	tempParam["type"] = 4

	base64Str, err := util.PutParmToStr(tempParam)
	if err != nil {
		logging.Error(fmt.Sprintf("GetWakeUpAppParams Error: alive_id: %s, err: %v", b.AliveRep.AliveId, err))
		return
	}
	url = base64Str
	return
}

// 获取旧直播间链接
func (b *BaseInfo) GetAliveRoomUrl(req validator.BaseInfoRuleV2) string {
	params := util.ContentParam{
		Type:         e.PaymentTypeReward,
		ResourceType: e.ResourceTypeLive,
		ResourceId:   req.ResourceId,
		ProductId:    req.ProductId,
		PaymentType:  int(b.Alive.PaymentType),
		ChannelId:    req.ChannelId,
		AppId:        b.AliveRep.AppId,
		ShareUserId:  req.ShareUserId,
		ShareType:    req.ShareType,
		ShareAgent:   req.ShareAgent,
		ShareFrom:    req.ShareFrom,
		Scene:        req.Scene,
		WebAlive:     req.WebAlive,
		Token:        req.InviteToken,
		ExtraData:    e.AliveRoomPage,
	}
	return util.ContentUrl(params)
}

// 获取直播自定义文案内容
func (b *BaseInfo) GetCaptionDefine(captionDefineJson string) map[string]string {
	captionDefine := make(map[string]string)
	captionDefine["home_title"] = "首页"
	captionDefine["home_tab_message"] = "消息"
	captionDefine["column_title"] = "专栏"
	captionDefine["column_open"] = "订阅"
	captionDefine["column_pay_hint"] = "订阅专栏"
	captionDefine["audio_try_hint"] = "购买"
	captionDefine["single_product_hint"] = "订阅专栏"
	if captionDefineJson != "" {
		m := make(map[string]string)
		err := util.JsonDecode([]byte(captionDefineJson), &m)
		if err == nil {
			if v, ok := m["home_title"]; ok && v != "" {
				captionDefine["home_title"] = v
			}
			if v, ok := m["home_tab_message"]; ok && v != "" {
				captionDefine["home_tab_message"] = v
			}
			if v, ok := m["column_title"]; ok && v != "" {
				captionDefine["column_title"] = v
			}
			if v, ok := m["column_open"]; ok && v != "" {
				captionDefine["column_open"] = v
			}
			if v, ok := m["column_pay_hint"]; ok && v != "" {
				captionDefine["column_pay_hint"] = v
			}
			if v, ok := m["audio_try_hint"]; ok && v != "" {
				captionDefine["audio_try_hint"] = v
			}
			if v, ok := m["single_product_hint"]; ok && v != "" {
				captionDefine["single_product_hint"] = v
			}
		}
	}
	return captionDefine
}

// 私有方法 ===============================================================================
// 判断用户是否打开快直播
func (b *BaseInfo) isUseFastLive(userId string) (bool, error) {
	conn, err := redis_alive.GetLiveInteractConn()
	if err != nil {
		return true, err
	}
	defer conn.Close()

	flag, err := redis.Bool(conn.Do("SISMEMBER", notUseFastLiveKey, b.AliveRep.AppId+userId))
	if err != nil {
		return true, err
	}
	return !flag, nil
}

// 根据清晰度替换播放链接, sharpness可切换的清晰度：default默认，fluent流畅
func (b *BaseInfo) getPlayUrlBySharpness(sharpness, playUrl, channelId string) string {
	replaceStr := ""
	switch sharpness {
	case "fluent":
		replaceStr = fmt.Sprintf("%s_%s", channelId, os.Getenv("ALIVE_SHARPNESS_SWITCH_FLUENT"))
	case "hd":
		replaceStr = fmt.Sprintf("%s_%s", channelId, os.Getenv("ALIVE_SHARPNESS_SWITCH_HD"))
	default:
		replaceStr = ""
	}

	if replaceStr != "" {
		playUrl = strings.Replace(playUrl, channelId, replaceStr, -1)
	}
	return playUrl
}

// 店铺版本决定默认是否开启快直播，老版本默认关闭，其余开启【老逻辑，现在不在用了好像】
func (b *BaseInfo) canUseFastLive(versionType int) bool {
	//允许开快直播的版本
	switch versionType {
	case e.VERSION_TYPE_PROBATION:
		return true
	case e.VERSION_TYPE_ONLINE_EDUCATION:
		return true
	case e.VERSION_TYPE_ADVANCED:
		return true
	case e.VERSION_TYPE_STANDARD:
		return true
	case e.VERSION_TYPE_TRAINING_STD:
		return true
	case e.VERSION_TYPE_TRAINING_TRY:
		return true
	default:
		return false
	}
}

// GetAppExpireTime 获取版本过期信息
func (b *BaseInfo) GetAppExpireTime(profit map[string]interface{}) map[string]interface{} {
	// 查询直播间插件功能过期信息
	versionStateArray := []string{
		"alive_voice",
		"alive_video_voice",
		"alive_PPT",
		"alive_push",
		"alive_share_file",
		"student_show",
		"alive_show_man_time",
		"alive_reward",
		"exercise",
		"hide_resource_count",
		"hide_sub_count",
	}

	var permissionArray map[string]interface{}
	if v, ok := profit["fp_conf"]; ok {
		permissionArray = v.(map[string]interface{})
	} else {
		permissionArray = make(map[string]interface{})
	}
	// 不兼容配置服务报错
	// permissionArray := profit["fp_conf"].(map[string]interface{}) //店铺配置
	result := make(map[string]interface{})
	for _, v := range versionStateArray {
		// 0-未过期 1-即将过期，2-过期，-1未购买
		_, ok := permissionArray[v]
		if !ok {
			result[v+"_is_remind"] = -1
		} else {
			if v == "video" || v == "alive_push" || v == "alive_video_voice" {
				result[v+"_is_remind"] = 1
			} else {
				expireTime, _ := time.Parse("2006-01-02 15:04:05", permissionArray[v].(string))
				leftTime := expireTime.Unix() - time.Now().Unix()
				result[v+"_expire_time"] = permissionArray[v]
				if leftTime > 8*24*3600 {
					result[v+"_is_remind"] = 0
				} else {
					if leftTime >= 0 {
						result[v+"_is_remind"] = 1
					} else {
						result[v+"_is_remind"] = 2
					}
				}
			}
		}
	}

	return result
}

// todo 录播底层优化-直播链接下发逻辑-func start
/**
1、判断是否有channel_id，无则返回 false
2、判断是否在灰度，不在灰度则不使用伪直播 false
3、判断推流状态是否为1 ，为1说明已经在推了，true
4、redis查询拉流转推任务状态,数据为空则查mysql数据，目前判断redis使用情况不到3%
*/
func (b *BaseInfo) GetNowRecordedIsPush() bool {
	if b.Alive.ChannelId == "" {
		return false
	}
	isGrayBool := redis_gray.InGrayShopSpecialHit("recorded_use_retweet", b.AliveRep.AppId)
	if !isGrayBool {
		return false
	}
	if b.Alive.PushState == 1 {
		return true
	}
	redisConn, err := redis_im.GetLiveGroupActionConn()
	if err != nil {
		logging.Error(err)
	}
	defer redisConn.Close()
	expire := 5
	if b.Alive.ZbStartAt.Time.Add(300 * time.Second).Before(time.Now()) {
		expire = 30
	}
	existTaskCacheKey := fmt.Sprintf("alive_exist_retweet_task_%s", b.AliveRep.AliveId)
	data, _ := redis.String(redisConn.Do("get", existTaskCacheKey))
	if data == "" {
		// redis没有数据，查mysql
		info, errVod := alive.GetRecordedRetweetTaskInfo(b.AliveRep.AppId, b.AliveRep.AliveId, "task_id,task_state")
		if errVod != nil {
			return false
		}
		data = "0"
		if (info.TaskId != "" && info.TaskState != 2) || info.TaskState == 4 {
			data = "1"
		}
		redisConn.Do("setex", existTaskCacheKey, expire, data)
	}
	if data == "1" {
		return true
	}
	return false
}

// todo 录播底层优化-直播链接下发逻辑-func end
