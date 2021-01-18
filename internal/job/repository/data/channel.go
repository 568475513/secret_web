package data

import (
	"strconv"

	"go.uber.org/zap"

	"abs/models/business"
	e "abs/pkg/enums"
	"abs/pkg/logging"
)

type Channels struct {
	AppId       string `json:"app_id"`
	ChannelId   string `json:"channel_id"`   // 渠道id
	ResourceId  string `json:"resource_id"`
	ProductId   string `json:"product_id"`   // payment_type为2时-NULL, payment_type为3时-绑定的付费产品包id
	PaymentType string `json:"payment_type"` // 付费类型：2-单笔、3-付费产品包
}

// 增加渠道浏览量
func (c *Channels) AddChannelViewCount() {
	// 如果是有渠道,渠道浏览量+1
	var hasChannel bool

	if c.ChannelId != "" {
		channelInDB, err := business.GetChannelInfo(c.AppId, c.ChannelId)
		if err != nil {
			logging.JLogger.Error(err.Error(), zap.Stack("stack"))
			// logging.Error(err)
			return
		}
		if channelInDB.ChannelType == 0 {
			if c.PaymentType == strconv.Itoa(e.PaymentTypeSingle) {
				if c.ResourceId == channelInDB.ResourceId {
					hasChannel = true
				} else if c.ProductId == channelInDB.ProductId {
					hasChannel = true
				}
			} else {
				hasChannel = false
			}
		} else {
			hasChannel = false
		}

		if hasChannel {
			business.UpdateChannelViewCount(c.AppId, c.ChannelId)
		}
	}
	return
}
