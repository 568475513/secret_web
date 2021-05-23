package course

import (
	"fmt"

	"abs/models/alive"
	"abs/models/user"
	e "abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
)

// 次级业务类
type Secondary struct {
	Alive    *alive.Alive
	UserInfo *user.User
	BuzUri   string
}

// 获取直播邀请讲师链接
func (s *Secondary) GetTeacherInvitationUrl() (url string) {
	tempParam := make(map[string]interface{})
	tempParam["title"] = s.Alive.Title.String
	tempParam["alive_id"] = s.Alive.Id
	tempParam["wx_nickname"] = s.UserInfo.WxNickname
	tempParam["inviteImg"] = s.UserInfo.WxAvatar
	tempParam["user_id"] = s.UserInfo.UserId
	tempParam["alive_type"] = s.Alive.AliveType

	base64Str, err := util.PutParmToStr(tempParam)
	if err != nil {
		logging.Error(fmt.Sprintf("GetTeacherInvitationUrl Error: alive_id: %s, user_id: %s", s.Alive.Id, s.UserInfo.UserId))
		return
	}
	url = util.UrlWrapper(fmt.Sprintf("/teacherInvitation/%s", base64Str), s.BuzUri, s.Alive.AppId)
	return
}

// 获取邀请卡链接
func (s *Secondary) GetInvitationCardUrl() (url string) {
	tempParam := make(map[string]interface{})
	tempParam["payment_type"] = e.PaymentTypeSingle
	tempParam["resource_type"] = e.ResourceTypeLive
	tempParam["resource_id"] = s.Alive.Id
	tempParam["product_id"] = ""

	base64Str, err := util.PutParmToStr(tempParam)
	if err != nil {
		logging.Error(fmt.Sprintf("GetInvitationCardUrl Error: alive_id: %s, user_id: %s", s.Alive.Id, s.UserInfo.UserId))
		return
	}
	url = util.UrlWrapper(fmt.Sprintf("/inviteCard/%s", base64Str), s.BuzUri, s.Alive.AppId)
	return
}

// 获取共享文件列表链接
func (s *Secondary) GetShareFileListUrl() (url string) {
	tempParam := make(map[string]interface{})
	tempParam["alive_id"] = s.Alive.Id

	base64Str, err := util.PutParmToStr(tempParam)
	if err != nil {
		logging.Error(fmt.Sprintf("GetShareFileListUrl Error: alive_id: %s, user_id: %s", s.Alive.Id, s.UserInfo.UserId))
		return
	}
	url = util.UrlWrapper(fmt.Sprintf("/share_file_list_page/%s", base64Str), s.BuzUri, s.Alive.AppId)
	return
}

// 获取邀请达人榜链接
func (s *Secondary) GetInvitationListUrl() (url string) {
	tempParam := make(map[string]interface{})
	tempParam["payment_type"] = e.PaymentTypeSingle
	tempParam["resource_type"] = e.ResourceTypeLive
	tempParam["resource_id"] = s.Alive.Id
	tempParam["product_id"] = ""

	base64Str, err := util.PutParmToStr(tempParam)
	if err != nil {
		logging.Error(fmt.Sprintf("GetInvitationListUrl Error: alive_id: %s, user_id: %s", s.Alive.Id, s.UserInfo.UserId))
		return
	}

	url = util.UrlWrapper(fmt.Sprintf("/inviteList/%s", base64Str), s.BuzUri, s.Alive.AppId)
	return
}

//// 获取邀请卡开关
//func (b *Secondary) GetInviteState(hasInvite int, paymentType int) (needInvite bool) {
//	if hasInvite == 1 &&
//		(b.Alive.PaymentType == e.PaymentTypeFree || (paymentType == e.PaymentTypeSingle && b.Alive.PaymentType == e.PaymentTypeSingle) || paymentType == e.PaymentTypeProduct) {
//		needInvite = true
//	}
//	return
//}
