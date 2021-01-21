package marketing

import (
	"abs/models/alive"
	"abs/models/business"
	e "abs/pkg/enums"
	"abs/pkg/logging"
)

// 没时间改了，就先用copy的代码吧
type Share struct {
	AppId          string
	UserId         string
	ProductId      string
	Alive          *alive.Alive
	shareResources []*business.ShareRecords
}

type ShareInfo struct {
	IsShareFree      uint8                    `json:"is_share_free"`
	ShareUserId      string                   `json:"share_user_id"`
	Num              int                      `json:"num"`
	SurplusNum       int                      `json:"surplus_num"`
	ShareResource    int                      `json:"share_resource"`
	HasShareResource []*business.ShareRecords `json:"has_share_resource"`
	ProductInfo      *business.PayProducts    `json:"product_info"`
}

type ShareListenInfo struct {
	IsShareListen    bool   `json:"is_share_listen"`
	ShareListenUser  string `json:"share_listen_user"`
	IsShowShareCount bool   `json:"is_show_share_count"`
}

func (share *Share) GetShareInfo(available, availableProduct bool, shareInfo ShareInfo) ShareInfo {
	// 免费 非加密 // 如果该资源或者当前专栏不可用 查询分享者信息
	if available == false || availableProduct == false {
		// 查询用户是否领取了 免费听
		availShareInfo, err := business.GetAvailShareInfo(share.AppId, e.ResourceTypeLive, share.Alive.Id, share.ProductId, share.UserId)
		if err != nil {
			logging.Error(err)
		}
		if availShareInfo.Id != 0 {
			// 如果已经领取了，取出分享人id，并查询领取人数
			shareInfo.ShareUserId = availShareInfo.ShareUserId
			// 已领取人数
			hasReceivedShare, err := business.GetHasReceivedShare(share.AppId, e.ResourceTypeLive, share.Alive.Id, share.ProductId, shareInfo.ShareUserId)
			if err != nil {
				logging.Error(err)
			}
			for k, v := range hasReceivedShare {
				if v.ReceiveUserId == share.UserId {
					shareInfo.Num = k + 1 // 第几个领取
				}
			}
			// 剩余可领取数量
			surplusNum := shareInfo.SurplusNum - len(hasReceivedShare)
			if surplusNum > 0 {
				shareInfo.SurplusNum = surplusNum
			} else {
				shareInfo.SurplusNum = 0
			}
		}
	} else {
		// 该资源可用  置空share_user_id
		shareInfo.ShareUserId = ""
		// 如果分享人不存在，即为分享者，查询该专栏可继续分享的资源数
		if shareInfo.ShareUserId == "" || shareInfo.ShareUserId == share.UserId {
			shareResource, _ := business.GetShareResource(share.AppId, share.ProductId, share.UserId)
			shareInfo.ShareResource = len(shareResource)
			shareInfo.HasShareResource = shareResource
			share.shareResources = shareResource
		}
	}
	return shareInfo
}

//
func (share *Share) GetShareInfoInit(parentColumns []*business.PayProducts) ShareInfo {
	shareInfo := ShareInfo{}
	// 如果专栏信息无误，且参与分享免费听
	if share.ProductId != "" {
		// 从该专栏信息中取出 邀请免费听的信息
		for _, value := range parentColumns {
			if value.Id == share.ProductId && share.Alive.PaymentType == 3 {
				shareInfo.IsShareFree = value.IsShareListen
				shareInfo.SurplusNum = value.ShareListenCount
				shareInfo.ProductInfo = value
				break
			}
		}
	}
	return shareInfo
}

// 分享免费听逻辑
func (sh *Share) GetShareListenInfo(shareInfo *ShareInfo, available bool) ShareListenInfo {
	var shareListenInfo ShareListenInfo = ShareListenInfo{IsShareListen: false, ShareListenUser: sh.UserId, IsShowShareCount: true}
	// shareInfo["is_share_free"]有可能是bool 有可能是int  只能判断下了
	if sh.Alive.PaymentType == e.PaymentTypeProduct && shareInfo.IsShareFree == 1 && available {
		if shareInfo.SurplusNum > 0 {
			shareProductInfo := shareInfo.ProductInfo
			if shareProductInfo.ShareListenResource < 1 || shareProductInfo.ShareListenResource > shareInfo.ShareResource {
				shareListenInfo.IsShareListen = true
			}
		}
		if shareListenInfo.IsShareListen {
			for _, v := range sh.shareResources {
				if v.ResourceId == sh.Alive.Id {
					shareListenInfo.IsShowShareCount = false
					break
				}
			}
			if shareInfo.ShareUserId != "" {
				shareListenInfo.ShareListenUser = shareInfo.ShareUserId
			}
		}
	}
	return shareListenInfo
}
