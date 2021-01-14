package course

import (
	"abs/models/sub_business"
	"abs/pkg/logging"
	"abs/pkg/util"
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
			Type:         "15",
			ResourceType: "23",
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
					redirect = util.ContentUrl(contentParam)
					break
				}
			}
		}
	}
	return
}
