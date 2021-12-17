package sub_business

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"abs/pkg/enums"
)

type ChannelRepertory struct {
	Id                        string  `json:"id"`
	ChannelAppId              string  `json:"channel_app_id"`
	ContentAppId              string  `json:"content_app_id"`
	ResourceId                string  `json:"resource_id"`
	ImgUrl                    string  `json:"img_url"`
	ImgUrlCompressed          string  `json:"img_url_compress"`
	ImgUrlCompressedLarger    string  `json:"img_url_compressed_larger"`
	ResourceType              int     `json:"resource_type"`
	ResourceName              string  `json:"resource_name"`
	Summary                   string  `json:"summary"`
	Descrb                    string  `json:"descrb"`
	OrgContent                string  `json:"org_content"`
	LinePrice                 uint    `json:"line_price"`
	HasDistribute             int8    `json:"has_distribute"`
	NewHasDistribute          int8    `json:"new_has_distribute"`
	IsShowUserinfo            int8    `json:"is_show_userinfo"`
	InvitePoster              string  `json:"invite_poster"`
	IsDistributeShowUserinfo  int8    `json:"is_distribute_show_userinfo"`
	DistributePoster          string  `json:"distribute_poster"`
	IsMember                  int8    `json:"is_member"`
	FirstDistributeDefault    int8    `json:"first_distribute_default"`
	FirstDistributePercent    float64 `json:"first_distribute_percent"`
	SuperiorDistributeDefault int8    `json:"superior_distribute_default"`
	SuperiorDistributePercent float64 `json:"superior_distribute_percent"`
	DistributePercent         float64 `json:"distribute_percent"`
	DistributeState           int8    `json:"distribute_state"`
	ContentDistributeState    int     `json:"content_distribute_state"`
	CustomEditState           int8    `json:"custom_edit_state"`
	SyncContentState          int8    `json:"sync_content_state"`
	Weight                    int     `json:"weight"`
	IsCompleteInfo            int8    `json:"is_complete_info"`
	ChannelSource             int8    `json:"channel_source"`
	Status                    int8    `json:"status"`
}

const (
	tableName = "t_channel_repertory"
)

// GetChannelRepertoryList 查询专栏数据
func GetChannelRepertoryList(contentAppId string, appId string, resourceIds []string, s []string) (result []*ChannelRepertory, err error) {
	err = db.Table(fmt.Sprintf("%s.%s", DatabaseContentMarket, tableName)).Select(s).
		Where("content_app_id = ? and resource_id in (?) and resource_type in (?) and channel_app_id = ?",
			appId, resourceIds, []int{enums.ResourceTypeTopic, enums.ResourceTypePackage}, contentAppId).Find(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return result, nil
}
