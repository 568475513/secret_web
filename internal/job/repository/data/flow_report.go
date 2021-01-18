package data

import (
	"fmt"

	"github.com/gomodule/redigo/redis"

	"abs/models/data"
	"abs/pkg/cache/redis_alive"
)

const (
	// 数据来源
	Applets          = 0 // 小程序
	OfficialAccounts = 1 // 公众号

	// 统计方式
	FrontEnd     = 0 // 前端页面上报
	BackEnd      = 1 // 后台直接上报
	LoopPlayback = 2 // 循环播放统计

	// Redis keys
	cacheRedisDataUage string = "alive_data_uage_%s"
	// 缓存时间控制(秒)
	// -
	redisDataUageCacheTime = 60 * 60 * 3
)

// 上报 DTO对象，直接操作PO多余得字段可能不好造成参数模糊
type FlowReportData struct {
	AppId             string // 店铺ID
	AliveId           string // 直播ID
	UserId            string // 用户ID
	ResourceType      int
	Title             string  // 直播名称
	VidioSize         float64 // 视频大小
	AliveM3u8HighSize float64 // meu8大小
	ImgSizeTotal      float64
	WxAppType         int     // 数据来源 0-小程序 1-公众号
	Way               int     // 统计方式：0-前端页面上报  1-后台直接上报 2-循环播放统计
}

type DataUageBusiness struct {
}

// 直接上报流量,同abs insertFlowRecord
func (business *DataUageBusiness) InsertFlowRecord(flowReportData FlowReportData) bool {
	// todo 先用着，后期增加缓存减少一次数据库得查询
	return business.BaseInsertFlowRecord(flowReportData)
}

// 直接上报流量,同abs insertAlive2FlowRecord
func (business *DataUageBusiness) InsertAlivePushFlowRecord(flowReportData FlowReportData, aliveState int8) bool {
	flowReportData.ImgSizeTotal = 0
	// 从abs迁过来得，数据库注释是（0-无、1-音频、2-视频、3-直播、4-图文、5-直播回放 ），不明白为什么是6
	flowReportData.ResourceType = 6
	if aliveState == 3 {
		flowReportData.ResourceType = 5
	} else {
		flowReportData.AliveM3u8HighSize = 0
		flowReportData.VidioSize = 0
	}
	return business.BaseInsertFlowRecord(flowReportData)
}

func (business *DataUageBusiness) BaseInsertFlowRecord(flowReportData FlowReportData) bool {
	// dto 转成 po
	dataUage := data.DataUage{
		AppId:          flowReportData.AppId,
		UserId:         flowReportData.UserId,
		ResourceId:     flowReportData.AliveId,
		ResourceType:   flowReportData.ResourceType,
		ResourceName:   flowReportData.Title,
		Size:           flowReportData.VidioSize,
		SizeCompressed: flowReportData.AliveM3u8HighSize,
		ImgSizeTotal:   flowReportData.ImgSizeTotal,
		WxAppType:      flowReportData.WxAppType,
		Way:            flowReportData.Way,
	}
	var cacheKey string = fmt.Sprintf(cacheRedisDataUage, data.GetTableName())
	conn, err := redis_alive.GetLiveBusinessConn()
	defer conn.Close()
	if err != nil {
		return false
	}
	cacheData, err := redis.Bytes(conn.Do("GET", cacheKey))
	// todo 先用着，后期增加缓存减少一次数据库得查询
	// 中断或，前面为true后面则不执行
	if cacheData == nil {
		if !data.IsHaveTable(data.GetTableName()) {
			// todo 日志
			if !data.CreateDataUsageTable() {
				// todo 日志
				return false
			}
		}
		// 缓存3小时，这个缓存出异常不会影响业务，会回源到DB查询
		conn.Do("SET", cacheKey, 1, "EX", redisDataUageCacheTime)
	}
	return data.InsertFlowRecord(dataUage)
}
