package data

import (
	"fmt"
	"strconv"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/gin-gonic/gin"

	"abs/models/alive"
	"abs/pkg/enums"
	"abs/pkg/job"
	"abs/pkg/logging"
)

// Todo: 请注意参数顺序，这个很重要！！！！
type AsynData struct {
	AppId       string
	UserId      string
	ResourceId  string
	ProductId   string
	PaymentType int
}

// 购买关系埋点异步数据上报
// 未完善，由于数据量过大需要慎重！
func (a *AsynData) AsynDataUserPurchase(c *gin.Context, available bool) error {
	addTask := &tasks.Signature{
		Name: "insert_user_purchase_log",
		Args: []tasks.Arg{
			{
				// Name:  "app_id",
				Type:  "string",
				Value: a.AppId,
			},
			{
				// Name:  "user_id",
				Type:  "string",
				Value: a.UserId,
			},
			{
				// Name:  "raw_url",
				Type:  "string",
				Value: c.Request.URL.String(),
			},
			{
				// Name:  "url",
				Type:  "string",
				Value: c.Request.URL.Path,
			},
			{
				// Name:  "referer",
				Type:  "string",
				Value: c.Request.Referer(),
			},
			{
				// Name:  "app_version",
				Type:  "string",
				Value: c.GetString("app_version"),
			},
			{
				// Name:  "agent",
				Type:  "string",
				Value: c.Request.UserAgent(),
			},
			{
				// Name:  "client",
				Type:  "string",
				Value: c.GetString("client"),
			},
			{
				// Name:  "use_collection",
				Type:  "string",
				Value: c.DefaultQuery("use_collection", "true"),
			},
			{
				// Name:  "ip",
				Type:  "string",
				Value: c.GetString("client_ip"),
			},
			{
				// Name:  "resource_type",
				Type:  "string",
				Value: strconv.Itoa(enums.ResourceTypeLive),
			},
			{
				// Name:  "resource_id",
				Type:  "string",
				Value: a.ResourceId,
			},
			{
				// Name:  "product_id",
				Type:  "string",
				Value: a.ProductId,
			},
			{
				// Name:  "is_resource_pay",
				Type:  "string",
				Value: strconv.FormatBool(available),
			},
		},
	}
	_, err := job.Machinery.SendTask(addTask)
	if err != nil {
		fmt.Printf("AsynDataUserPurchase异步上报错误: %s\n", err.Error())
		logging.Error(err)
	}

	return err
}

// 增加渠道浏览量
func (a *AsynData) AsynChannelViewCount(channelId string) error {
	addTask := &tasks.Signature{
		Name: "add_channel_view_count",
		Args: []tasks.Arg{
			{
				// Name:  "app_id",
				Type:  "string",
				Value: a.AppId,
			},
			{
				// Name:  "resource_id",
				Type:  "string",
				Value: a.ResourceId,
			},
			{
				// Name:  "product_id",
				Type:  "string",
				Value: a.ProductId,
			},
			{
				// Name:  "payment_type",
				Type:  "string",
				Value: strconv.Itoa(a.PaymentType),
			},
			{
				// Name:  "channel_id",
				Type:  "string",
				Value: channelId,
			},
		},
	}
	_, err := job.Machinery.SendTask(addTask)
	if err != nil {
		fmt.Printf("AsynChannelViewCount异步上报错误: %s\n", err.Error())
		logging.Error(err)
	}

	return err
}

// 直接上报流量
func (a *AsynData) AsynFlowRecord(aliveInfo *alive.Alive, available bool, aliveState int) error {
	resourceType, vidioSize, aliveM3u8HighSize := enums.ResourceTypeVideo, aliveInfo.VideoSize, aliveInfo.AliveM3u8HighSize
	if aliveInfo.AliveType == enums.AliveTypeVideo && available {                                                                                                                           //视频直播
		//直播类型（如果直播结束就是回看类型）
		switch aliveState {
		case 1:
			resourceType = 3
		case 3:
			resourceType = 5
		}
	} else if (aliveInfo.AliveType == enums.AliveTypePush || aliveInfo.AliveType == enums.AliveOldTypePush) && available { // 推流直播上报流量
		resourceType = 6
		if aliveState == 3 {
			resourceType = 5
		} else {
			vidioSize, aliveM3u8HighSize = float64(0), float64(0)
		}
	}
	addTask := &tasks.Signature{
		Name: "insert_flow_record",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: a.AppId,
			},
			{
				Type:  "string",
				Value: a.UserId,
			},
			{
				Type:  "string",
				Value: strconv.Itoa(resourceType), // 为什么是3 ???
			},
			{
				Type:  "string",
				Value: a.ResourceId,
			},
			{
				Type:  "string",
				Value: aliveInfo.Title.String,
			},
			{
				Type:  "string",
				Value: strconv.FormatFloat(vidioSize, 'g', 6, 64),
			},
			{
				Type:  "string",
				Value: strconv.FormatFloat(aliveM3u8HighSize, 'g', 6, 64),
			},
			{
				Type:  "string",
				Value: "0",
			},
			{
				Type:  "string",
				Value: "1",
			},
			{
				Type:  "string",
				Value: "1",
			},
		},
	}
	_, err := job.Machinery.SendTask(addTask)
	if err != nil {
		fmt.Printf("AsynFlowRecord异步上报错误: %s\n", err.Error())
		logging.Error(err)
	}

	return err
}