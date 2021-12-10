package course

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gomodule/redigo/redis"

	"abs/models/alive"
	"abs/models/sub_business"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
)

type AliveInfo struct {
	AppId   string
	AliveId string
}

const (
	// Redis key
	aliveInfoKey        = "base_info_alive_info:%s:%s"
	aliveModuleConf     = "alive_module_conf:%s:%s"
	aliveCircuitBreaker = "alive:circuitBreaker"
	// 直播静态相关
	staticAliveHashId   = "hash_static_alive_id_%s"
	staticAliveHashUser = "hash_static_alive_user_%s"
	// view_count店铺id跟直播id集合
	viewCountSetKey      = "view_count_set_key"
	viewCountTimeKeyNew  = "view_count_update_set_time_new:%s:%s"
	aliveViewCountNew    = "alive_view_count_new:%s:%s"    // 直播访问量
	forbiddenUserListKey = "forbidden_user_list_key:%s:%s" // 直播禁言
	// 带货PV
	pvCacheKeyPre    = "alive_take_goods_pv:%s:%s:%s:"              // pv缓存键
	timeCacheKeyPre  = "alive_take_goods_pv_refresh_time:%s:%s:%s:" // pv缓存上一次刷新时间键
	allPvSetCacheKey = "alive_take_goods_pv_set:"                   // 所有带货商品pv集合缓存
	expirationTime   = 300                                          // pv缓存有效时间，单位秒

	hitImActive      = "active_im_group_all_cache_%s"     // IM活跃群组
	oldHitImActive   = "old_active_im_group_all_cache_%s" // IM活跃群组
	imGroupActive    = "im_active_:%s"
	oldImGroupActive = "old_im_active_:%s"
	// 不使用快直播名单
	notUseFastLiveKey = "notUseFastLive"

	// 直播间排行榜
	aLiveRankingKey = "abs:alive_ic_grkn_%s_%s"

	// 缓存时间控制(秒)
	// 直播详情
	aliveInfoCacheTime = "60"
	// 直播的ModuleConf
	aliveModuleConfCacheTime = "60"
)

// 获取直播详情
func (a *AliveInfo) GetAliveInfo() (cacheAliveInfo *alive.Alive, err error) {
	conn, err := redis_alive.GetLiveBusinessConn()
	if err != nil {
		logging.Error(err)
	}
	defer conn.Close()

	cacheKey := fmt.Sprintf(aliveInfoKey, a.AppId, a.AliveId)
	info, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err == nil {
		if err = util.JsonDecode(info, &cacheAliveInfo); err != nil {
			logging.Error(err)
		} else {
			return
		}
	}
	// 数据库获取
	cacheAliveInfo, err = alive.GetAliveInfo(a.AppId, a.AliveId, []string{
		"id",
		"app_id",
		"is_complete_info",
		"product_id",
		"payment_type",
		"room_id",
		"summary",
		"org_content",
		"zb_start_at",
		"zb_stop_at",
		"product_name",
		"title",
		"alive_video_url",
		"video_length",
		"manual_stop_at",
		"file_id",
		"alive_type",
		"img_url",
		"img_url_compressed",
		"alive_img_url",
		"aliveroom_img_url",
		"can_select",
		"distribute_percent",
		"has_distribute",
		"distribute_poster",
		"first_distribute_default",
		"first_distribute_percent",
		"recycle_bin_state",
		"state",
		"start_at",
		"is_stop_sell",
		"is_public",
		"config_show_view_count",
		"config_show_reward",
		"have_password",
		"is_discount",
		"is_public",
		"piece_price",
		"line_price",
		"comment_count",
		"view_count",
		"channel_id",
		"push_state",
		"rewind_time",
		"play_url",
		"push_url",
		"ppt_imgs",
		"push_ahead",
		"if_push",
		"is_lookback",
		"is_takegoods",
		"create_mode",
		"forbid_talk",
		"show_on_wall",
		"can_record",
	})
	// 未查到在此处返回
	if err != nil || cacheAliveInfo.Id == "" {
		return
	}

	// 缓存
	if value, err := util.JsonEncode(cacheAliveInfo); err == nil {
		if _, err = conn.Do("SET", cacheKey, value, "EX", aliveInfoCacheTime); err != nil {
			logging.Error(err)
		}
	} else {
		logging.Error(err)
	}

	// Redis错误不影响返回
	return cacheAliveInfo, nil
}

// 获取缓存里面的直播评论ViewCount
func (a *AliveInfo) GetAliveViewCountFromCache() (viewCount int, err error) {
	conn, err := redis_alive.GetLiveInteractConn()
	if err != nil {
		logging.Error(err)
		return
	}
	defer conn.Close()

	viewCount, err = redis.Int(conn.Do("GET", fmt.Sprintf(aliveViewCountNew, a.AppId, a.AliveId)))
	if err != nil {
		logging.Warn(fmt.Sprintf("获取缓存里面的直播评论ViewCount失败：%s", err.Error()))
		return
	}
	return
}

// 获取直播的ModuleConf
func (a *AliveInfo) GetAliveModuleConf() (*alive.AliveModuleConf, error) {
	var cacheAliveModuleConf *alive.AliveModuleConf
	conn, _ := redis_alive.GetSubBusinessConn()
	defer conn.Close()

	cacheKey := fmt.Sprintf(aliveModuleConf, a.AppId, a.AliveId)
	info, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err == nil {
		if err = util.JsonDecode(info, &cacheAliveModuleConf); err != nil {
			logging.Error(err)
		}
		return cacheAliveModuleConf, nil
	}

	cacheAliveModuleConf, err = alive.GetAliveModuleConf(a.AppId, a.AliveId, []string{"*"})
	if err != nil {
		logging.Error(err)
		return cacheAliveModuleConf, err
	}

	if value, err := util.JsonEncode(cacheAliveModuleConf); err == nil {
		if _, err = conn.Do("SET", cacheKey, value, "EX", aliveModuleConfCacheTime); err != nil {
			logging.Error(err)
		}
	}

	return cacheAliveModuleConf, nil
}

// 获取讲师信息
func (a *AliveInfo) GetAliveRole(userId string) (isRole uint, roleInfo map[string]interface{}, err error) {
	roleInfo = map[string]interface{}{
		"user_title":         "",
		"is_can_exceptional": 0,
		"main_teacher":       "",
		"main_user_id":       "",
		"role_user_id":       "",
	}
	aliveRoles, err := alive.GetAliveRole(a.AppId, a.AliveId)
	if err != nil {
		logging.Error(err)
		return
	}

	for _, v := range aliveRoles {
		if v.IsCurrentLecturer == 1 {
			roleInfo["main_teacher"] = v.UserName.String
			roleInfo["main_user_id"] = v.UserId.String
		}
		if v.UserId.String == userId {
			isRole = 1
			roleInfo["user_title"] = v.RoleName.String
			roleInfo["is_can_exceptional"] = int(v.IsCanExceptional)
		}
		roleInfo["role_user_id"] = v.UserId.String
	}
	return
}

// 查询直播是否被禁言
func (a *AliveInfo) GetAliveImIsShow(roomId, userId string) (isShow int) {
	isShow = 1
	aliveForbids, err := alive.GetAliveForbid(a.AppId, roomId, userId)
	if err != nil {
		logging.Error(err)
		return
	}
	if aliveForbids.IsUseful > 0 {
		isShow = 0
	}
	return
}

// 查询直播是否被禁言【Redis版】
// 暂时不可用，可问jessica
func (a *AliveInfo) GetAliveImIsShowForRedis(roomId, userId string) (isShow int) {
	isShow = 1
	conn, _ := redis_alive.GetForbiddenUserConn()
	defer conn.Close()

	cacheKey := fmt.Sprintf(forbiddenUserListKey, a.AppId, roomId)
	isExist, err := redis.Int(conn.Do("HEXISTS", cacheKey, userId, userId))
	if err != nil {
		logging.Error(err)
		return
	}

	// 存在即被禁言
	if isExist == 1 {
		isShow = 0
	}
	return
}

// 获取直播间是否被封禁
func (a *AliveInfo) GetAliveRoomIsBan() bool {
	isBan, err := sub_business.ResourceIsBan(a.AppId, a.AliveId)
	if err != nil {
		logging.Error(err)
		isBan = false
	}
	return isBan
}

// 获取直播回放链接的状态
func (a *AliveInfo) GetAliveLookBackStates(aliveInfo *alive.Alive) (aliveState int) {
	if aliveInfo.AliveType == enums.AliveTypePush || aliveInfo.AliveType == enums.AliveOldTypePush {
		// 视频推流直播
		aliveState = a.GetAliveState(aliveInfo.ZbStartAt.Time, aliveInfo.ZbStopAt.Time, aliveInfo.ManualStopAt.Time, aliveInfo.RewindTime.Time, aliveInfo.PushState)
	} else if aliveInfo.AliveType == enums.AliveTypeVideo {
		// 视频直播状态
		aliveState = a.GetAliveStateUtils(aliveInfo.ZbStartAt.Time, aliveInfo.VideoLength, aliveInfo.ManualStopAt.Time, aliveInfo.ZbStopAt.Time, aliveInfo.PushState)
	} else {
		// 语音或ppt直播
		aliveState = a.GetAliveStateForOthers(aliveInfo.ZbStartAt.Time, aliveInfo.ManualStopAt.Time, aliveInfo.ZbStopAt.Time)
	}
	return
}

// 获取直播的状态
func (a *AliveInfo) GetAliveStates(aliveInfo *alive.Alive) (aliveState int) {
	now := time.Now()
	if aliveInfo.AliveType == enums.AliveTypePush || aliveInfo.AliveType == enums.AliveOldTypePush {
		// 视频推流直播
		aliveState = a.GetAliveState(aliveInfo.ZbStartAt.Time, aliveInfo.ZbStopAt.Time, aliveInfo.ManualStopAt.Time, aliveInfo.RewindTime.Time, aliveInfo.PushState)
		// 互动状态默认互动时间为五分钟
		if aliveState == enums.AliveTypePush {
			aliveInfo.ZbStopAt.Time = aliveInfo.ZbStopAt.Time.Add(300 * time.Second)
		}
		// 直播已经开始（提前开始解决提前开始倒计时问题）
		if aliveState != 0 {
			if aliveInfo.ZbStartAt.Time.Add(60 * time.Second).After(now) {
				aliveInfo.ZbStartAt.Time = now
			}
			// 提前一分钟，避免客户端与服务器的时间差
			m, _ := time.ParseDuration("-1m")
			aliveInfo.ZbStartAt.Time = aliveInfo.ZbStartAt.Time.Add(m)
		}
	} else if aliveInfo.AliveType == enums.AliveTypeVideo {
		// 视频直播状态
		aliveState = a.GetAliveStateUtils(aliveInfo.ZbStartAt.Time, aliveInfo.VideoLength, aliveInfo.ManualStopAt.Time, aliveInfo.ZbStopAt.Time, aliveInfo.PushState)
	} else {
		// 语音或ppt直播
		aliveState = a.GetAliveStateForOthers(aliveInfo.ZbStartAt.Time, aliveInfo.ManualStopAt.Time, aliveInfo.ZbStopAt.Time)
	}
	return
}

// 获取推流直播的状态
func (a *AliveInfo) GetAliveState(start time.Time, stop time.Time, mst time.Time, rt time.Time, pstate uint8) (state int) {
	// 直播状态:0-还未开始  1-直播中  2-互动时间  3-直播结束了（回播） 4-离开
	// 如果没有手动结束 && 现在时间比预设开始时间早 && 未推流，则未开始直播
	now := time.Now()
	if mst.IsZero() && now.Before(start) && pstate == 2 {
		return
	}

	// 播放已开始
	state = 1
	//手动结束 && 现在的时间大于手动结束时间
	if !mst.IsZero() && now.After(mst) {
		state = 3
	}
	if mst.IsZero() && pstate != 1 {
		// 设定直播时间已经到了,并且断流
		if now.After(stop) {
			// 断流超过5分钟
			if rt.Add(300 * time.Second).Before(now) {
				// 直播结束
				state = 3
			} else {
				// 等待推流
				state = 4
			}
		} else { // 直播时间内断流等待
			state = 4
			if pstate == 0 && start.After(now) {
				state = 0
			}
		}
	}
	return
}

// 视频直播状态
func (a *AliveInfo) GetAliveStateUtils(start time.Time, total int64, mst time.Time, stop time.Time, pstate uint8) (state int) {
	if mst.IsZero() && pstate == 1 {
		state = 1
		return
	}

	now := time.Now()
	if now.After(start) {
		state = 1 //播放已开始
		//判断视频是否结束了
		if now.Unix()-start.Unix() >= total {
			state = 2
		}

		//判断直播是否结束了
		if mst.IsZero() && !stop.IsZero() {
			if now.After(stop) {
				state = 3
			}
		} else if !mst.IsZero() && !stop.IsZero() {
			zbStopAt := stop
			if stop.Unix() > mst.Unix() {
				zbStopAt = mst
			}
			if now.After(zbStopAt) {
				state = 3
			}
		}
	}

	return
}

//	语音或ppt直播的直播状态
func (a *AliveInfo) GetAliveStateForOthers(start time.Time, mst time.Time, stop time.Time) (state int) {
	// 直播状态:0-还未开始  1-直播中  2-互动时间  3-直播结束了（回播）
	state = 0
	now := time.Now()
	if now.After(start) {
		//开始直播
		state = 1
		//判断直播是否结束
		if mst.IsZero() && !stop.IsZero() {
			if now.After(stop) {
				state = 3
			}
		} else if !mst.IsZero() && !stop.IsZero() {
			zbStopAt := stop
			if stop.Unix() > mst.Unix() {
				zbStopAt = mst
			}
			if now.After(zbStopAt) {
				state = 3
			}
		}
	}

	return
}

// 直播次数加一，PV+1
func (a *AliveInfo) UpdateViewCountToCache(viewCount int) (int, error) {
	redisConn, err := redis_alive.GetLiveInteractConn()
	// 直接数据库写入
	if err != nil {
		err = alive.UpdateViewCount(a.AppId, a.AliveId, viewCount+1)
		logging.Error(err)
		return viewCount, err
	}
	defer redisConn.Close()

	// 更新周期，先设置时间为5分钟
	updateTime := 300
	// 优先查询set里面有没有该直播, 如果有，则走新逻辑新key
	key := fmt.Sprintf("%s:%s", a.AppId, a.AliveId)
	setTimeKey := fmt.Sprintf(viewCountTimeKeyNew, a.AppId, a.AliveId)
	viewCountKey := fmt.Sprintf(aliveViewCountNew, a.AppId, a.AliveId)
	isExist, err := redis.Bool(redisConn.Do("sismember", viewCountSetKey, key))
	if err != nil {
		isExist = false
	}
	setTime, err := redis.Int(redisConn.Do("get", setTimeKey))
	if err != nil {
		setTime = 0
	}
	viewCountByRedis, err := redis.Bool(redisConn.Do("exists", viewCountKey))
	if err != nil {
		viewCountByRedis = false
	}

	if isExist != false && setTime != 0 && viewCountByRedis != false {
		// redis有值，判断是否到更新周期时间，到更新时间则更新到数据库，并重置key，没到更新周期则更新缓存
		viewCount, err = redis.Int(redisConn.Do("incr", viewCountKey))
		if int(time.Now().Unix())-setTime >= updateTime {
			redisConn.Do("set", setTimeKey, time.Now().Unix())
			// 直接数据库写入
			err = alive.UpdateViewCount(a.AppId, a.AliveId, viewCount)
			if err != nil {
				logging.Error(err)
			}
		}
	} else {
		// 写入redis
		viewCount = viewCount + 1
		redisConn.Do("sadd", viewCountSetKey, key)
		redisConn.Do("set", setTimeKey, time.Now().Unix())
		redisConn.Do("set", viewCountKey, viewCount)
	}

	return viewCount, err
}

// 直播带货商品PV+1
func (a *AliveInfo) IncreasePv(referer, resourceId string, resourceType int) bool {
	parse, err := url.Parse(referer)
	if err != nil {
		return false
	}
	queryParam, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return false
	}

	aliveId, liveRoom := "", ""
	aliveId = queryParam.Get("aliveId")
	liveRoom = queryParam.Get("live_room")
	if aliveId == "" || liveRoom != "1" {
		return false
	}
	a.updatePv(aliveId, resourceType)
	return true
}

// 更新直播带货商品PV
func (a *AliveInfo) updatePv(resourceId string, resourceType int) {
	redisConn, err := redis_alive.GetLiveInteractConn()
	if err != nil {
		logging.Error(err)
	}
	defer redisConn.Close()
	pvCacheKey := fmt.Sprintf(pvCacheKeyPre, a.AppId, resourceId, a.AliveId)
	pvRefreshCacheKey := fmt.Sprintf(timeCacheKeyPre, a.AppId, resourceId, a.AliveId)
	pvSetValue := fmt.Sprintf("%s:%s:%s", a.AppId, resourceId, a.AliveId)
	pv := 1
	currentTime := time.Now().Unix()
	isExist, err := redis.Int(redisConn.Do("sismember", allPvSetCacheKey, pvSetValue))
	if err != nil {
		logging.Error(err)
		isExist = 0
	}
	if isExist == 0 {
		//无缓存的情况（被脚本消费了，或者是第一次带货访问）
		pvInfo, err := alive.GetTaskGoodsInfo(a.AppId, resourceId, a.AliveId, []string{"view_count"})
		if err != nil {
			logging.Error(err)
			return
		}
		if pvInfo.ViewCount == 0 {
			//初始化pv记录
			tgd := alive.TaskGoodsDetail{
				AppId:        a.AppId,
				AliveId:      resourceId,
				ResourceId:   a.AliveId,
				ResourceType: resourceType,
				ViewCount:    1,
				State:        1,
			}
			err = alive.InsertTaskGoodsInfo(tgd)
			if err != nil {
				logging.Error(err)
				return
			}
		} else {
			pv = pvInfo.ViewCount + 1
		}
		redisConn.Do("sadd", allPvSetCacheKey, pvSetValue)
		redisConn.Do("set", pvRefreshCacheKey, currentTime)
	} else {
		//有缓存的情况
		pv, _ = redis.Int(redisConn.Do("get", pvCacheKey))
		pv = pv + 1
		pvRefreshTime, _ := redis.Int64(redisConn.Do("get", pvRefreshCacheKey))
		if currentTime-pvRefreshTime >= expirationTime {
			//到了刷新时间则更新到数据库，并更新缓存刷新时间
			err = alive.UpdateTaskGoodsViewCount(a.AppId, resourceId, a.AliveId, pv)
			if err != nil {
				logging.Error(err)
				return
			}
			redisConn.Do("set", pvRefreshCacheKey, currentTime)
		}
	}
	redisConn.Do("set", pvCacheKey, pv)
}
