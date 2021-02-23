package marketing

import (
	"database/sql"
	"github.com/jinzhu/gorm"

	"abs/models/business"
	e "abs/pkg/enums"
	"abs/pkg/logging"
)

type InviteBusiness struct {
	AppId  string
	UserId string
}

type InviteUserInfo struct {
	ShareUserId  string
	ResourceId   string
	ResourceType int
	PaymentType  int
	ProductId    string
}

// 冗余一点方便以后扩展
type InviteRelation struct {
	InviteUserInfo
	ShareType int
}

// 0 = 音频分享 1 = 日签分享 2 = 专栏分享 4-邀请卡分享
const (
	Audio        int = 0
	DaySignature int = 1
	SpecialClumn int = 2
	InvitaCard   int = 4
	NotKnow          = 5 // 5不晓得是什么鬼

)

// 更新邀请关系及邀请数量 邀请卡与推广合并后的新逻辑
func (businesss *InviteBusiness) AddInviteCountUtilsNew(inviteUserInfo InviteUserInfo) {
	if inviteUserInfo.PaymentType == e.PaymentTypeSingle {
		inviteUserInfo.ProductId = ""
	}
	inviteUserInfoPo := businesss.transformInviteUserPo(inviteUserInfo)
	if inviteUserInfo.ShareUserId != "" && businesss.UserId != "" && businesss.UserId != inviteUserInfo.ShareUserId {
		inviteRelation := InviteRelation{ShareType: NotKnow, InviteUserInfo: inviteUserInfo}
		inviteRelationPo := businesss.transformInviteRelationPo(inviteRelation)
		_, err := business.GetInviteUserByInvitedUser(inviteRelationPo)
		if err == nil {
			return
		}
		if err != gorm.ErrRecordNotFound {
			logging.Error(err)
			return
		}
		count := business.UpdateInviteUserByInviteCount(inviteUserInfoPo)
		if count <= 0 {
			// 通过错误来避免一次数据库查询
			inviteUserInfoPo.InviteCount = 1
			if err := business.SetInviteUser(inviteUserInfoPo); err != nil {
				logging.Error(err)
				return
			}
		}
		if err = business.SetInviteRelation(inviteRelationPo); err != nil {
			logging.Error(err)
		}

	}
}

// dto模型向底层po模型转换
func (businesss *InviteBusiness) transformInviteUserPo(inviteUserInfo InviteUserInfo) business.InviteUser {
	return business.InviteUser{
		AppId:        businesss.AppId,
		PaymentType:  inviteUserInfo.PaymentType,
		ResourceId:   transformGormString(inviteUserInfo.ResourceId),
		ResourceType: inviteUserInfo.ResourceType,
		ShareUserId:  inviteUserInfo.ShareUserId,
		ProductId:    inviteUserInfo.ProductId,
	}
}

// dto模型向底层po模型转换
func (businesss *InviteBusiness) transformInviteRelationPo(inviteUserInfo InviteRelation) business.InviteRelation {
	return business.InviteRelation{
		AppId:         businesss.AppId,
		PaymentType:   inviteUserInfo.PaymentType,
		ResourceId:    transformGormString(inviteUserInfo.ResourceId),
		ResourceType:  inviteUserInfo.ResourceType,
		ProductId:     inviteUserInfo.ProductId,
		ShareUserId:   inviteUserInfo.ShareUserId,
		ShareType:     inviteUserInfo.ShareType,
		InvitedUserId: businesss.UserId,
	}
}

// 好坑，先这样搞一波
func transformGormString(str string) (nullString sql.NullString) {
	nullString.String = str
	nullString.Valid = true
	return nullString
}
