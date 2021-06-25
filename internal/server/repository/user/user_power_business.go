package user

import (
	"abs/pkg/cache/redis_gray"
	"abs/pkg/logging"
	"abs/service"
)

type userPowerBusiness struct {
	AppId     string
	UserId    string
	AgentType int
	service   service.AvailableService
}

// 初始化用户权益服务
func UserPowerBusiness(appId string, userId string, agentType int) *userPowerBusiness {
	return &userPowerBusiness{
		appId,
		userId,
		agentType,
		service.AvailableService{AppId: appId, UserId: userId},
	}
}

// 直播权益判断
func (upb *userPowerBusiness) IsHaveAlivePower(resouceId string, resouceType string, needExpire bool) (string, bool) {
	if redis_gray.InGrayShopSpecial("is_switch_new_permission", upb.AppId) {
		result, err := upb.IsEncryAliveAccess(resouceId)
		if err != nil {
			logging.Error(err)
		}
		return "", result
	} else {
		// 参数动态配置
		resourceAvailable := service.ResourceAvailable{
			ResourceId:   resouceId,
			ResourceType: resouceType,
			NeedExpire:   needExpire,
			AgentType:    upb.AgentType,
		}
		return upb.service.IsResourceAvailable(resourceAvailable)
	}
}

// 专栏权益判断
func (upb *userPowerBusiness) IsHaveSpecialColumnPower(productId string) (string, bool) {
	productAvailable := service.ProductAvailable{
		ProductId: productId,
		AgentType: upb.AgentType,
	}
	return upb.service.IsProductAvailable(productAvailable)
}

// 内部课程权限判断
func (upb *userPowerBusiness) IsInsideAliveAccess(resouceId string) (bool, error) {
	return upb.service.IsResourceAccess(resouceId, false, 0)
}

// 加密、付费以及免费的课程权益判断
func (upb *userPowerBusiness) IsEncryAliveAccess(resouceId string) (bool, error) {
	return upb.service.IsResourceAccess(resouceId, false, 0)
}
