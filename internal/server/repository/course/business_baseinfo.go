package course

import (
	"strconv"
	"time"

	"abs/internal/server/rules/validator"
	"abs/models/alive"
	"abs/models/business"
	"abs/pkg/cache/alive_static"
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

// 组装直播的基本信息
func (b *BaseInfo) GetAliveInfoDetail(userId string) map[string]interface{} {
	aliveInfoDetail, now := make(map[string]interface{}), time.Now()
	aliveInfoDetail["app_id"] = b.AliveRep.AppId
	aliveInfoDetail["alive_id"] = b.Alive.Id
	aliveInfoDetail["room_id"] = b.Alive.RoomId
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
	// 获取直播状态
	aliveInfoDetail["alive_state"] = b.AliveRep.GetAliveStates(b.Alive)
	// 推流状态，0推流结束，1推流中，2推流未开始
	aliveInfoDetail["push_state"] = b.Alive.PushState
	// 直播剩余时长（秒）
	aliveInfoDetail["remainder_time"] = 0
	if b.Alive.AliveType == e.AliveTypeVideo {
		aliveInfoDetail["remainder_time"] = b.Alive.ZbStartAt.Unix() + b.Alive.VideoLength - now.Unix()
	}
	// 推流直播开始时间
	aliveInfoDetail["pushzb_start_at"] = b.Alive.ZbStartAt
	// 推流直播结束时间
	aliveInfoDetail["pushzb_stop_at"] = b.Alive.ZbStopAt
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
	aliveInfoDetail["push_url"] = b.Alive.PushUrl
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
	aliveInfoDetail["user_id"] = userId
	// 参数加密串
	aliveInfoDetail["param_str"], _ = util.PutParmToStr(map[string]interface{}{
		"payment_type":  e.PaymentTypeSingle,
		"resource_type": e.ResourceTypeLive,
		"resource_id":   b.Alive.Id,
		"product_id":    "",
	})

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
	// 判断是否是讲师,讲师不用付费
	if !available && b.UserType == 1 {
		availableInfo["available"] = true
	}
	// 买赠功能
	// availableInfo["gift_buy"] = 0
	// if b.Alive.State == 1 {
	// 	availableInfo["in_recycle"] = b.Alive.RecycleBinState
	// } else {
	// 	availableInfo["in_recycle"] = 0
	// }
	// if b.Alive.IsStopSell == 0 && b.Alive.RecycleBinState == 0 && b.Alive.StartAt.Time.After(now) {
	// 	availableInfo["time_left"] = b.Alive.StartAt.Time.Sub(now).Seconds()
	// } else {
	// 	availableInfo["time_left"] = 0
	// }
	return availableInfo
}

// 组装直播店铺配置信息
func (b *BaseInfo) GetAliveConfInfo(baseConf *service.AppBaseConf, aliveModule *alive.AliveModuleConf) map[string]interface{} {
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
	// 标准版是否可用 0-不可用  1-可用
	versionUse := 1
	// 标准版或试用版过期
	if versionState := util.JudgeDate(baseConf.VersionType, baseConf.ExpireTime); versionState["type"] == 1 || versionState["type"] == 3 {
		versionUse = 0
	}
	//开启PC网校 0为关闭 1为开启
	aliveConf["open_pc_network_school"] = baseConf.IsEnable
	//PC网校自定义域名
	aliveConf["pc_network_school_index_url"] = baseConf.PcCustomDomain

	// 版本过期信息
	versionState := b.getAppExpireTime(baseConf.Profit)
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

	// 获取直播配置表相关配置
	aliveConf["is_message_on"] = aliveModule.IsMessageOn
	aliveConf["alive_mode"] = aliveModule.AliveMode
	aliveConf["is_invite_on"] = aliveModule.IsInviteOn
	aliveConf["is_coupon_on"] = aliveModule.IsCouponOn
	aliveConf["is_card_on"] = aliveModule.IsCardOn
	aliveConf["is_prize_on"] = aliveModule.IsPrizeOn
	aliveConf["complete_time"] = aliveModule.CompleteTime
	// 是否开启打赏， 0-关闭 1-开启
	aliveConf["is_show_reward_on"] = aliveModule.IsShowRewardOn
	if aliveModule.CompleteTime == 0 {
		aliveConf["is_open_complete_time"] = 0
	} else {
		aliveConf["is_open_complete_time"] = 1
	}

	return aliveConf
}

// 直播静态页的信息采集【用户】
func (b *BaseInfo) SetAliveUserToStaticRedis(userId string) {
	err := alive_static.HsetNxString(staticAliveHashUser, b.Alive.Id+userId, 1, 3600*24)
	if err != nil {
		logging.Error(err)
	}
}

// 直播静态页的信息采集【ID】
func (b *BaseInfo) SetAliveIdToStaticRedis() {
	err := alive_static.HsetNxString(staticAliveHashId, b.Alive.Id, b.Alive.Id, 3600*24)
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
	if len(products) == 0 && b.Alive.PaymentType == 3 && !available {
		svipReq := Svip{AppId: b.Alive.AppId, ResourceId: b.Alive.Id}
		if redirect := svipReq.GetResourceSvipRedirect(); redirect != "" {
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
			Type:        strconv.Itoa(e.PaymentTypeProduct),
			ProductId:   b.Alive.ProductId.String,
			ChannelId:   req.ChannelId,
			ShareUserId: req.ShareUserId,
			ShareType:   strconv.Itoa(req.ShareType),
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
				urlColumParams.ResourceType = strconv.Itoa(e.ResourceTypeLive)
				urlColumParams.ShareUserId = req.ShareUserId
				urlColumParams.ShareType = strconv.Itoa(req.ShareType)
				url = util.ParentColumnsUrl(urlColumParams)
			} else if len(products) == 1 {
				urlColumParams.Type = strconv.Itoa(e.PaymentTypeProduct)
				urlColumParams.ResourceId = ""
				urlColumParams.ResourceType = ""
				urlColumParams.ProductId = products[0].Id
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
// 获取版本过期信息
func (b *BaseInfo) getAppExpireTime(profit map[string]interface{}) map[string]interface{} {
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
	}

	permissionArray := profit["fp_conf"].(map[string]interface{}) //店铺配置
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
