package business

import (
	"github.com/jinzhu/gorm"

	"abs/pkg/provider/json"
)

const (
	// 会员
	SINGLE_GOODS_MEMBER = 5
	// 专栏
	SINGLE_GOODS_PACKAGE = 6
)

type PayProducts struct {
	// Model

	AppId                  string              `json:"app_id"`
	Id                     string              `json:"id"`
	ImgUrl                 json.JSONNullString `json:"img_url"`
	ImgUrlCompressed       json.JSONNullString `json:"img_url_compressed"`
	Name                   json.JSONNullString `json:"name"`
	Summary                json.JSONNullString `json:"summary"`
	PurchaseCount          int                 `json:"-"` // purchase_count
	Price                  int                 `json:"price"`
	DistributePercent      float64             `json:"-"` // distribute_percent
	FirstDistributePercent float64             `json:"-"` // first_distribute_percent
	IsMember               uint8               `json:"-"` // is_member
	MemberType             uint8               `json:"-"` // member_type
	MemberIconDefault      json.JSONNullString `json:"-"` // member_icon_default
	Period                 int                 `json:"-"` // period
	IsShowResourceCount    int                 `json:"is_show_resourcecount"`
	IsShareListen          uint8               `json:"-"` // is_share_listen
	ShareListenResource    int                 `json:"-"` // share_listen_resource
	ShareListenCount       int                 `json:"-"` // share_listen_count
	ResourceCount          uint                `json:"resource_count"`
	RecycleBinState        uint8               `json:"-"` // recycle_bin_state
	State                  uint8               `json:"state"`

	RatePrice  int      `json:"rate_price"`
	SrcType    uint8    `json:"srcType"`
	InActivity int      `json:"in_activity"`
	Tags       []string `json:"tags"`
}

// 获取状态值筛选的专栏资源
func GetPayProductState(appId string, resourceId string) (*PayProducts, error) {
	var pp PayProducts
	err := db.Select([]string{"id"}).Where("app_id=? and id=? and state=?", appId, resourceId, 0).First(&pp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &pp, nil
}

// 获取指定id的专栏资源
func GetPayProductByIds(appId string, ids string) ([]*PayProducts, error) {
	var pp []*PayProducts
	err := db.Select([]string{"id",
		"app_id",
		"img_url_compressed",
		"name",
		"summary",
		"purchase_count",
		"price",
		"distribute_percent",
		"first_distribute_percent",
		"is_member",
		"member_type",
		"member_icon_default",
		"period",
		"is_show_resourcecount",
		"is_share_listen",
		"share_listen_resource",
		"share_listen_count",
		"resource_count",
		"recycle_bin_state",
		"state",
		"img_url"}).Where("app_id=? and state!=? and id in (?)", appId, 2, ids).First(&pp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return pp, nil
}
