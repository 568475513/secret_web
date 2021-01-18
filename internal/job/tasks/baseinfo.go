package tasks

import (
	"strconv"
	"time"

	"github.com/RichardKnop/machinery/v1/log"

	"abs/internal/job/repository/data"
)

// 购买关系埋点 ...
func InsertUserPurchaseLog(args ...string) (bool, error) {
	client, _ := strconv.Atoi(args[7])
	uCollect, _ := strconv.ParseBool(args[8])
	available, _ := strconv.ParseBool(args[13])
	dataRep := data.BuryingPoint{}
	userPurchase := &data.UserPurchaseData{
		AppId:          args[0],
		UserId:         args[1],
		RawUrl:         args[2],
		Url:            args[3],
		Referer:        args[4],
		AppVersion:     args[5],
		Agent:          args[6],
		Client:         int8(client),
		UserCollection: uCollect,
		Ip:             args[9],
		ResourceType:   args[10],
		ResourceId:     args[11],
		ProductId:      args[12],
		IsResourcePay:  available,
	}

	dataRep.InsertUserPurchaseLog(userPurchase)
	return true, nil
}

// 增加渠道浏览量 ...
func AddChannelViewCount(args ...string) (bool, error) {
	// 渠道上报实例
	channelRepository := &data.Channels{
		AppId:       args[0],
		ChannelId:   args[1],
		ResourceId:  args[2],
		PaymentType: args[3],
		ProductId:   args[4],
	}

	channelRepository.AddChannelViewCount()
	return true, nil
}

// 直接上报流量 ...
func InsertFlowRecord(args ...string) (bool, error) {
	reourceType, _ := strconv.Atoi(args[2])
	vidioSize, _ := strconv.ParseFloat(args[5], 64)
	aliveM3u8HighSize, _ := strconv.ParseFloat(args[6], 64)
	imgSizeTotal, _ := strconv.ParseFloat(args[7], 64)
	wxAppType, _ := strconv.Atoi(args[8])
	way, _ := strconv.Atoi(args[9])
	// 流量上报实例
	dataUageBusiness := &data.DataUageBusiness{}
	// 流量上报处理
	// 流量上报结构体
	flowReportData := data.FlowReportData{
		AppId:             args[0],
		UserId:            args[1],
		ResourceType:      reourceType,
		AliveId:           args[3],
		Title:             args[4],
		VidioSize:         vidioSize,
		AliveM3u8HighSize: aliveM3u8HighSize,
		ImgSizeTotal:      imgSizeTotal,
		WxAppType:         wxAppType,
		Way:               way,
	}

	return dataUageBusiness.InsertFlowRecord(flowReportData), nil
}

// LongRunningTask ...
func LongRunningTask() error {
	log.INFO.Print("Long running task started")
	for i := 0; i < 10; i++ {
		log.INFO.Print(10 - i)
		time.Sleep(1 * time.Second)
	}
	log.INFO.Print("Long running task finished")
	return nil
}