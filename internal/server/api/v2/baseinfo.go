package v2

import (
	"abs/internal/server/repository/app_conf"
	"abs/internal/server/repository/course"
	"abs/internal/server/repository/data"
	"abs/internal/server/repository/marketing"
	"abs/internal/server/repository/material"
	ruser "abs/internal/server/repository/user"
	"abs/internal/server/rules/validator"
	"abs/pkg/app"
	"abs/pkg/enums"
	"abs/pkg/util"

	// service做变量初始化
	"abs/service"
	// Model层不可以直接调用，这里只能做变量初始化
	malive "abs/models/alive"
	mbusiness "abs/models/business"
	muser "abs/models/user"

	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

// @Summary 直播间基础信息
func GetBaseInfo(c *gin.Context) {
	var (
		err error
		req validator.BaseInfoRuleV2
	)
	userId := app.GetUserId(c)
	req.AppId = app.GetAppId(c)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("直播主业务调用链路", opentracing.ChildOf(app.GetTracingSpan(c)))
	span.SetTag("params", req)
	span.SetTag("userid", userId)
	defer span.Finish()

	// 获取直播详情内容
	childSpan := tracer.StartSpan("获取直播详情内容", opentracing.ChildOf(span.Context()))
	aliveRep := course.AliveInfo{AppId: req.AppId, AliveId: req.ResourceId}
	aliveInfo, err := aliveRep.GetAliveInfo()
	childSpan.Finish()
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("获取直播基础信息错误:%s", err.Error()), enums.ERROR, c)
		return
	}
	if aliveInfo == nil || aliveInfo.Id == "" || aliveInfo.State == enums.AliveStateDelete {
		app.FailWithMessage("内容已被删除", enums.Code_Db_Not_Find, c)
		return
	}

	// 直播静态化查询操作
	aliveStaticRep := course.AliveStatic{AppId: req.AppId, AliveId: req.ResourceId, UserId: userId, Type: req.Type}
	aliveStaticData, err := aliveStaticRep.AliveStaticMain(c.GetInt("agent_type"))
	if len(aliveStaticData) > 0 {
		app.OkWithData(aliveStaticData, c)
		return
	}
	// 直播专栏关联信息
	childSpan = tracer.StartSpan("获取直播专栏关联信息", opentracing.ChildOf(span.Context()))
	proRep := course.Product{AppId: req.AppId, ResourceId: req.ResourceId}
	aliveRelations, err := proRep.GetResourceRelation()
	childSpan.Finish()
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("获取直播专栏关联信息错误:%s", err.Error()), enums.ERROR, c)
		return
	}

	// 协程组查询数据包
	// 标准来说请尽量不要用共享内存来实现协程通信（GC压力高等），少量服务查询请用channel通信，此处为了代码可读性
	childSpan = tracer.StartSpan("协程组查询数据包", opentracing.ChildOf(span.Context()))
	bT := time.Now()
	var (
		aliveModule      *malive.AliveModuleConf
		available        bool
		availableProduct bool
		expireAt         string
		products         []*mbusiness.PayProducts
		termList         []*mbusiness.PayProducts
		baseConf         *service.AppBaseConf
		roleInfo         map[string]interface{}
		userType         uint
	)
	// 初始化权益实例
	ap := ruser.UserPowerBusiness(req.AppId, userId, c.GetInt("agent_type"))
	// 初始化店铺配置相关
	appRep := app_conf.AppInfo{AppId: req.AppId}
	// 此处需要补充致命错误输出后立刻返回
	err = app.GoroutineNotPanic(func() (err error) {
		goSpan := tracer.StartSpan("查询次级业务数据", opentracing.ChildOf(childSpan.Context()))
		defer goSpan.Finish()
		// 查询讲师数据
		userType, roleInfo, err = aliveRep.GetAliveRole(userId)
		// 获取直播配置信息
		aliveModule, err = aliveRep.GetAliveModuleConf()
		// 资源一对多（对应所有专栏信息）
		products, err = proRep.GetParentColumns(aliveRelations)
		return nil
	}, func() (err error) {
		// 获取店铺配置相关
		goSpan := tracer.StartSpan("获取店铺配置相关", opentracing.ChildOf(childSpan.Context()))
		defer goSpan.Finish()
		baseConf, err = appRep.GetConfHubInfo()
		// 报错不抛出
		return nil
	}, func() (err error) {
		// 获取营期内容
		goSpan := tracer.StartSpan("获取营期内容", opentracing.ChildOf(childSpan.Context()))
		defer goSpan.Finish()
		termList, err = proRep.GetCampTermListByIds(aliveRelations)
		return nil
	}, func() (err error) {
		// 用户权益
		if aliveInfo.IsPublic == 0 {
			goSpan := tracer.StartSpan("用户内部培训权益", opentracing.ChildOf(childSpan.Context()))
			defer goSpan.Finish()
			available, err = ap.IsInsideAliveAccess(req.ResourceId)
		} else {
			if aliveInfo.PaymentType == enums.PaymentTypeFree && aliveInfo.HavePassword != 1 && aliveInfo.State == 0 {
				available = true
			} else {
				if aliveInfo.HavePassword == 1 {
					goSpan := tracer.StartSpan("用户密码直播权益", opentracing.ChildOf(childSpan.Context()))
					defer goSpan.Finish()
					available, err = ap.IsEncryAliveAccess(req.ResourceId)
				} else {
					goSpan := tracer.StartSpan("用户权益", opentracing.ChildOf(childSpan.Context()))
					defer goSpan.Finish()
					expireAt, available = ap.IsHaveAlivePower(req.ResourceId, strconv.Itoa(enums.ResourceTypeLive), true)
				}
			}
		}
		return
	}, func() (err error) {
		if aliveInfo.PaymentType == enums.PaymentTypeFree && aliveInfo.HavePassword != 1 {
			// 目前这里只是针对免费直播不进行查询专栏的订购关系，赋默认值, 默认值为false，这个参数目前只用于微信初始化接口，慎用其他地方
			availableProduct = false
		} else {
			// 专栏权益
			goSpan := tracer.StartSpan("专栏权益", opentracing.ChildOf(childSpan.Context()))
			defer goSpan.Finish()
			// 如果该资源或者当前专栏不可用 查询分享者信息
			_, availableProduct = ap.IsHaveSpecialColumnPower(req.ProductId)
		}
		return nil
	}, func() (err error) {
		// 直播异步操作
		goSpan := tracer.StartSpan("直播异步操作", opentracing.ChildOf(childSpan.Context()))
		defer goSpan.Finish()
		// 直播Pv数加一
		aliveInfo.ViewCount, _ = aliveRep.UpdateViewCountToCache(aliveInfo.ViewCount)
		// 直播带货商品PV加一
		aliveRep.IncreasePv(c.Request.Referer(), aliveInfo.Id, int(aliveInfo.AliveType))
		return nil
	})
	childSpan.Finish()
	// fmt.Println("BaseInfo的协程处理时间: ", time.Since(bT))
	// 错误处理【需要扔掉一些不要的】
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("并行请求组错误: %s[%s]", err.Error(), time.Since(bT)), enums.ERROR, c)
		return
	}
	// 公开课跳转
	if aliveInfo.IsPublic == 0 && !available && userType == 0 && !util.IsQyApp(baseConf.VersionType) {
		app.FailWithMessage("内部课程，暂无权限", enums.FORBIDDEN, c)
		return
	}
	// 数据组装阶段
	childSpan = tracer.StartSpan("数据组装阶段", opentracing.ChildOf(span.Context()))
	// 替换Redis里面的真实ViewCount
	// Todo: 0308发现跟老abs逻辑不符合，注释了
	// if viewCount, err := aliveRep.GetAliveViewCountFromCache(); err == nil {
	// 	aliveInfo.ViewCount = viewCount
	// }
	// 替换第一个专栏内容【显示用】
	if len(aliveRelations) > 0 {
		aliveInfo.ProductId.String = aliveRelations[0].ProductId
		aliveInfo.ProductName.String = aliveRelations[0].ProductName.String
		aliveInfo.IsTry = aliveRelations[0].IsTry
	} else {
		aliveInfo.ProductId.String = ""
		aliveInfo.ProductName.String = ""
	}
	// 替换全局设置
	for _, val := range products {
		val.IsShowResourceCount = baseConf.IsShowResourcecount
	}
	// 给专栏添加活动标签
	// products = marketing.GetActivityTags(products, 2, c.GetString("client"), c.GetString("app_version"))
	// 组合营期内容
	products = append(products, termList...)
	// 判断是否是讲师,讲师不用付费
	if available == false && userType == 1 {
		available = true
	}
	// 邀请好友免费听逻辑 免费 非加密
	shareRes := marketing.Share{AppId: req.AppId, UserId: userId, ProductId: req.ProductId, Alive: aliveInfo}
	shareInfo := shareRes.GetShareInfoInit(products)
	// 邀请好友免费听：如果领取了免费听 则将该资源置位可用！
	if aliveInfo.PaymentType != enums.PaymentTypeFree || aliveInfo.HavePassword == 1 {
		shareInfo = shareRes.GetShareInfo(available, availableProduct, shareInfo)
		if shareInfo.ShareUserId != "" && shareInfo.Num > 0 {
			available = true
		}
	}
	// 分享免费听逻辑
	shareListenInfo := shareRes.GetShareListenInfo(&shareInfo, available)
	// 业务数据封装
	baseInfoRep := course.BaseInfo{Alive: aliveInfo, AliveRep: &aliveRep, UserType: userType}
	aliveInfoDetail := baseInfoRep.GetAliveInfoDetail()
	aliveConf := baseInfoRep.GetAliveConfInfo(baseConf, aliveModule)
	availableInfo := baseInfoRep.GetAvailableInfo(available, availableProduct, expireAt)
	// 回放服务
	lookBackRep := material.LookBack{AppId: req.AppId, AliveId: req.ResourceId}
	lookBackExpire, _ := lookBackRep.GetLookbackExpire(int(aliveInfo.IsLookback), aliveModule.LookbackTime)
	// 补充回放过期信息
	aliveConf["lookback_time"] = lookBackExpire["lookback_time"]
	// 补充讲师信息
	aliveInfoDetail["user_id"] = userId
	aliveInfoDetail["user_title"] = roleInfo["user_title"]
	aliveConf["is_can_exceptional"] = roleInfo["is_can_exceptional"]
	// 补充老直播间链接 util.GetAliveRoomUrl(req.ResourceId, req.ProductId, req.ChannelId, req.AppId, enums.AliveRoomPage)
	aliveInfoDetail["old_live_room_url"] = baseInfoRep.GetAliveRoomUrl(req)
	// 获取播放连接【错误处理需要仓库层打印】
	alivePlayInfo := baseInfoRep.GetAliveLiveUrl(c.GetInt("agent_type"), baseConf.VersionType, baseConf.EnableWebRtc, userId)
	// 直播静态化写入操作
	if available && (aliveInfoDetail["alive_state"].(int) == 1 || aliveInfo.ZbStartAt.Format("2006-01-02") == time.Now().Format("2006-01-02")) {
		baseInfoRep.SetAliveIdToStaticRedis()
		if aliveInfo.PaymentType != 1 {
			baseInfoRep.SetAliveUserToStaticRedis(userId)
		}
	}
	childSpan.Finish()

	// 数据上报服务
	childSpan = tracer.StartSpan("异步队列处理时间", opentracing.ChildOf(span.Context()))
	dataAsyn := data.AsynData{AppId: req.AppId, UserId: userId, ResourceId: req.ResourceId, ProductId: req.ProductId, PaymentType: int(aliveInfo.PaymentType)}
	// 用户购买关系埋点上报
	dataAsyn.AsynDataUserPurchase(c, available)
	// 增加渠道浏览量
	dataAsyn.AsynChannelViewCount(req.ChannelId)
	// 直接上报流量
	dataAsyn.AsynFlowRecord(aliveInfo, available, aliveInfoDetail["alive_state"].(int))
	childSpan.Finish()

	// 开始组装数据
	data := make(map[string]interface{})
	// 父级专栏信息列表
	// data["parent_columns"] = products
	// 直播权益信息
	data["available_info"] = availableInfo
	// 直播基本信息
	data["alive_info"] = aliveInfoDetail
	// 直播播放信息
	data["alive_play"] = alivePlayInfo
	// 直播配置信息
	data["alive_conf"] = aliveConf
	// 直播分享邀请免费听逻辑
	data["share_info"] = map[string]interface{}{
		"share_info":        shareInfo,
		"share_listen_info": shareListenInfo,
	}
	// 直播自定义文案
	data["caption_define"] = baseInfoRep.GetCaptionDefine(baseConf.CaptionDefine)
	// 首页链接
	data["index_url"] = util.UrlWrapper("homepage", c.GetString("buz_uri"), req.AppId)
	// 页面是否跳转
	if url, code, msg := baseInfoRep.BaseInfoPageRedirect(products, available, baseConf.VersionType, req); code != 0 {
		app.OkWithCodeData(msg, map[string]string{"url": url}, code, c)
		return
	} else {
		data["payment_url"] = url
	}
	app.OkWithData(data, c)
}

// @Summary 获取直播间次级业务信息
func GetSecondaryInfo(c *gin.Context) {
	userId := app.GetUserId(c)
	appId := app.GetAppId(c)
	// 参数校验
	var req validator.SecondaryInfoRuleV2
	if err := app.ParseRequest(c, &req); err != nil {
		return
	}
	// 初始化店铺配置相关
	appRep := app_conf.AppInfo{AppId: appId}
	// 直播静态化查询操作
	if req.StaticIsStart != "" { //如果携带固定参数则走静态页
		ImInit := appRep.GetCommunicationCloudInfo(userId, "", true) //获取im配置
		userRep := ruser.UserBusinessConstrct(appId, userId)
		userInfo, _ := userRep.GetUserInfo()
		aliveStaticRep := course.AliveStatic{AppId: appId, UserId: userId}
		StaticData := aliveStaticRep.SecondaryInfoStaticData(ImInit, userInfo)
		app.OkWithData(StaticData, c)
		return
	}
	//直播间信息初始化
	aliveRep := course.AliveInfo{AppId: appId, AliveId: req.ResourceId}
	aliveInfo, err := aliveRep.GetAliveInfo()
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("获取直播基础信息失败:%s", err.Error()), enums.ERROR, c)
		return
	}
	// 协程组查询数据包
	bT := time.Now()
	var (
		userInfo     muser.User
		appMsgSwitch int
		isShow       int
		blackInfo    service.UserBlackInfo
		baseConf     *service.AppBaseConf
	)
	data := map[string]interface{}{"alive_id": aliveInfo.Id}
	// 初始化用户实例
	userRep, userInfoMap := ruser.UserBusinessConstrct(appId, userId), make(map[string]interface{})
	err = app.GoroutineNotPanic(func() (err error) {
		// 获取用户的基本信息
		userInfo, err = userRep.GetUserInfo()
		// Todo 记得这里是否抛出错误
		return nil
	}, func() (err error) {
		// 查询用户是否在黑名单【返回错误不抛出】
		blackInfo, err = userRep.GetUserBlackStates()
		return nil
	}, func() (err error) {
		// 查询短信预约总开关
		appMsgSwitch, err = appRep.GetAppConfSwitchState()
		return
	}, func() (err error) {
		// 查询直播间是否被禁言
		isShow = aliveRep.GetAliveImIsShow(aliveInfo.RoomId, userId)
		return nil
	}, func() (err error) {
		// 获取店铺配置
		baseConf, err = appRep.GetConfHubInfo()
		return
	})
	// fmt.Println("GetSecondaryInfo的协程处理时间: ", time.Since(bT))
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("并行请求组错误: %s[%s]", err.Error(), time.Since(bT)), enums.ERROR, c)
		return
	}
	baseInfoRep := course.Secondary{Alive: aliveInfo, UserInfo: &userInfo, BuzUri: c.GetString("buz_uri")}
	// 写入邀请关系
	if baseInfoRep.GetInviteState(baseConf.HasInvite, req.PaymentType) && aliveInfo.PaymentType == enums.PaymentTypeFree {
		inviteBusiness := marketing.InviteBusiness{AppId: appId, UserId: userId}
		inviteBusiness.AddInviteCountUtilsNew(marketing.InviteUserInfo{
			ShareUserId:  req.ShareUserId,
			PaymentType:  2, // 这个payment_type有坑，对老代码妥协的结果
			ResourceType: enums.ResourceTypeLive,
			ResourceId:   req.ResourceId,
			ProductId:    req.ProductId,
		})
	}
	// 组装用户信息
	userInfoMap["phone"] = userInfo.Phone
	userInfoMap["wx_avatar"] = userInfo.WxAvatar
	userInfoMap["wx_nickname"] = userInfo.WxNickname
	// 用户信息
	data["user_info"] = userInfoMap
	// 短信预约总开关
	data["is_message_on"] = appMsgSwitch
	// 用户是否被禁言
	data["is_show"] = isShow
	// 用户黑名单
	data["black_list"] = blackInfo
	// 邀请卡链接
	data["invite_card_url"] = baseInfoRep.GetInvitationCardUrl()
	// 邀请讲师链接
	data["invite_teacher_url"] = baseInfoRep.GetTeacherInvitationUrl()
	// 邀请达人榜链接
	data["invite_list_url"] = baseInfoRep.GetInvitationListUrl()
	// 共享文件列表链接
	data["share_file_url"] = baseInfoRep.GetShareFileListUrl()
	// 获取云通信配置
	data["im_init"] = appRep.GetCommunicationCloudInfo(userId, aliveInfo.RoomId, false)
	app.OkWithData(data, c)
}

// @Todo 备份使用中
// @Summary 直播间数据上报接口
// func DataReported(c *gin.Context) {
// 	var (
// 		err error
// 		req validator.DataReportedV2
// 	)
// 	//参数校验
// 	if err = app.ParseRequest(c, &req); err != nil {
// 		return
// 	}

// 	//获取直播详情
// 	aliveRep := course.AliveInfo{AppId: req.AppId, AliveId: req.ResourceId}
// 	aliveInfo, err := aliveRep.GetAliveInfo()
// 	if err != nil {
// 		app.FailWithMessage(fmt.Sprintf("获取直播基础信息错误:%s", err.Error()), enums.ERROR, c)
// 		return
// 	}

// 	//用户ID
// 	userId := app.GetUserId(c)
// 	//初始化权益实例
// 	ap := ruser.UserPowerBusiness(req.AppId, userId, c.GetInt("agent_type"))
// 	//渠道上报实例
// 	channelRepository := &data.Channels{
// 		AppId:       req.AppId,
// 		ChannelId:   req.ChannelId,
// 		ResourceId:  req.ResourceId,
// 		PaymentType: req.PaymentType,
// 		ProductId:   req.ProductId,
// 	}
// 	//流量上报实例
// 	dataUageBusiness := &data.DataUageBusiness{}
// 	//流量上报处理
// 	//流量上报结构体
// 	flowReportData := data.FlowReportData{
// 		AppId:             req.AppId,
// 		UserId:            app.GetUserId(c),
// 		ResourceType:      3,
// 		AliveId:           req.ResourceId,
// 		Title:             aliveInfo.Title.String,
// 		VidioSize:         aliveInfo.VideoSize,
// 		AliveM3u8HighSize: aliveInfo.AliveM3u8HighSize,
// 		ImgSizeTotal:      float64(0),
// 		WxAppType:         1,
// 		Way:               1,
// 	}

// 	available, _ := ap.IsInsideAliveAccess(aliveInfo.Id)                                                                                                                 //权益
// 	aliveState := aliveRep.GetAliveState(aliveInfo.ZbStartAt.Time, aliveInfo.ZbStopAt.Time, aliveInfo.ManualStopAt.Time, aliveInfo.RewindTime.Time, aliveInfo.PushState) //直播状态
// 	if aliveInfo.AliveType == 1 && available {                                                                                                                           //视频直播
// 		//直播类型（如果直播结束就是回看类型）
// 		switch aliveState {
// 		case 1:
// 			flowReportData.ResourceType = 3
// 		case 3:
// 			flowReportData.ResourceType = 5
// 		}
// 	} else if (aliveInfo.AliveType == 2 || aliveInfo.AliveType == 4) && available { //推流直播上报流量
// 		flowReportData.ResourceType = 6
// 		if aliveState == 3 {
// 			flowReportData.ResourceType = 5
// 		} else {
// 			flowReportData.VidioSize, flowReportData.AliveM3u8HighSize = float64(0), float64(0)
// 		}
// 	}

// 	//协程组执行IO处理
// 	err = app.GoroutineNotPanic(
// 		func() error {
// 			//增加渠道浏览量
// 			channelRepository.AddChannelViewCount()
// 			return nil
// 		},
// 		func() error {
// 			//直接上报流量
// 			dataUageBusiness.InsertFlowRecord(flowReportData)
// 			return nil
// 		},
// 		func() error {
// 			// 用户购买关系上报
// 			if aliveInfo.IsPublic != 0 {
// 				if aliveInfo.PaymentType == enums.PaymentTypeFree && aliveInfo.HavePassword != 1 && aliveInfo.State == 0 {
// 					available = true
// 				} else {
// 					if aliveInfo.HavePassword == 1 {
// 						available, err = ap.IsEncryAliveAccess(req.ResourceId)
// 					} else {
// 						_, available = ap.IsHaveAlivePower(req.ResourceId, strconv.Itoa(enums.ResourceTypeLive), true)
// 					}
// 				}
// 			}
// 			dataRep := data.BuryingPoint{AppId: req.AppId, UserId: userId, ResourceId: req.ResourceId, ProductId: req.ProductId}
// 			dataRep.InsertDataUserPurchase(c, available)
// 			return nil
// 		},
// 	)
// 	app.OkWithData("OK", c)
// }
