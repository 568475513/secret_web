package course

import (
	"abs/models/sub_business"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"
)

type Svip struct {
	AppId        string
	ResourceId   string
	ResourceType int
}

// 获取svip的跳转链接
func (s *Svip) GetResourceSvipRedirect() (redirect string) {
	relation, err := sub_business.GetResourceSvipRelation(s.AppId, s.ResourceId, s.ResourceType)
	if err != nil {
		logging.Error(err)
	} else {
		contentParam := util.ContentParam{
			Type:         15,
			ResourceType: 23,
			ResourceId:   "",
			ProductId:    relation.SvipId,
			AppId:        s.AppId,
		}
		if relation.Id != 0 {
			redirect = util.ContentUrl(contentParam)
		} else {
			svips, err := sub_business.GetSvipList(s.AppId)
			if err != nil {
				logging.Error(err)
			}
			for _, v := range svips {
				if v.EffactiveRange == 1 {
					contentParam.ProductId = v.Id
					redirect = util.ContentUrl(contentParam)
					break
				}
			}
		}
	}
	return
}

// GetResourceSvipRedirectV2 超级会员三期 获取svip的跳转链接
func (s *Svip) GetResourceSvipRedirectV2() (redirect string) {
	ss := service.SvipService{AppId: s.AppId}
	relation, err := ss.GetSvipBindRes(s.ResourceId, s.ResourceType)
	if err != nil {
		logging.Error(err)
	} else {
		if len(relation) == 1 {
			contentParam := util.ContentParam{
				Type:         15,
				ResourceType: 23,
				ResourceId:   "",
				ProductId:    relation[0].SvipID,
				AppId:        s.AppId,
			}
			redirect = util.ContentUrl(contentParam)
		} else if len(relation) > 1 {
			contentParam := util.ContentParam{
				ResourceType: s.ResourceType,
				ResourceId:   s.ResourceId,
			}
			redirect = util.ParentColumnsUrl(contentParam)
		}
	}
	return
}
