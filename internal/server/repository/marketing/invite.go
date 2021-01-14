package marketing

import (
	"fmt"
	"strconv"

	"abs/pkg/util"
)

type InviteInfo struct {
	AppId        string
	ResourceId   string
	ResourceType int
	PaymentType  int
	ProductId    string
}

type TeacherInviteUrl struct {
	Title      string `json:"title"`
	AliveId    string `json:"alive_id"`
	WxNickName string `json:"wx_nickname"`
	UserId     string `json:"user_id"`
	InviteImg  string `json:"inviteImg"`
	AliveType  int    `json:"alive_type"`
}

//生成邀请卡链接（currentUrl等价于abs的buz_uri）
func (invite *InviteInfo) GetInviteUrl(currentUrl string) string {
	params := map[string]string{
		"payment_type":  strconv.Itoa(invite.PaymentType),
		"resource_type": strconv.Itoa(invite.ResourceType),
		"resource_id":   invite.ResourceId,
		"product_id":    invite.ProductId,
	}
	path := fmt.Sprintf("/inviteCard/%s", util.SafeBase64Encode(params))
	return util.UrlWrapper(path, currentUrl, invite.AppId)
}

//生成邀请讲师链接（currentUrl等价于abs的buz_uri）
func (invite *InviteInfo) GetInviteTeacherUrl(currentUrl string, teacherInviteUrl TeacherInviteUrl) string {
	path := fmt.Sprintf("/teacherInvitation/%s", util.SafeBase64Encode(teacherInviteUrl))
	return util.UrlWrapper(path, currentUrl, invite.AppId)
}
