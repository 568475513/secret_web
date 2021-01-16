package data

import (
	"fmt"
	"strconv"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/gin-gonic/gin"

	"abs/pkg/enums"
	"abs/pkg/job"
	"abs/pkg/logging"
)

type AsynData struct {
	AppId       string
	UserId      string
	ResourceId  string
	ProductId   string
}

// 购买关系埋点异步数据上报
// 未完善，由于数据量过大放弃！
func (a *AsynData) AsynDataUserPurchase(c *gin.Context, available bool) error {
	addTask := &tasks.Signature{
		Name: "insert_user_purchase_log",
		Args: []tasks.Arg{
			{
				Name:  "app_id",
				Type:  "string",
				Value: a.AppId,
			},
			{
				Name:  "user_id",
				Type:  "string",
				Value: a.UserId,
			},
			{
				Name:  "raw_url",
				Type:  "string",
				Value: c.Request.URL.String(),
			},
			{
				Name:  "url",
				Type:  "string",
				Value: c.Request.URL.Path,
			},
			{
				Name:  "referer",
				Type:  "string",
				Value: c.Request.Referer(),
			},
			{
				Name:  "app_version",
				Type:  "string",
				Value: c.GetString("app_version"),
			},
			{
				Name:  "agent",
				Type:  "string",
				Value: c.Request.UserAgent(),
			},
			{
				Name:  "client",
				Type:  "string",
				Value: c.GetString("client"),
			},
			{
				Name:  "use_collection",
				Type:  "string",
				Value: c.DefaultQuery("use_collection", "1"),
			},
			{
				Name:  "ip",
				Type:  "string",
				Value: c.GetString("client_ip"),
			},
			{
				Name:  "resource_type",
				Type:  "string",
				Value: strconv.Itoa(enums.ResourceTypeLive),
			},
			{
				Name:  "resource_id",
				Type:  "string",
				Value: a.ResourceId,
			},
			{
				Name:  "product_id",
				Type:  "string",
				Value: a.ProductId,
			},
			{
				Name:  "is_resource_pay",
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