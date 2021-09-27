package user

import (
	"abs/pkg/cache/alive_static"
	//"abs/pkg/cache/redis_gray"
	"abs/pkg/logging"
	"abs/service"
	"fmt"
	"time"
)
const(
	staticAliveHashUser = "hash_static_alive_user_%s"
)
type userPowerBusiness struct {
	AppId     string
	UserId    string
	AgentType int
	service   service.AvailableService
	ResourceId   string
	ResourceType string
}

// 初始化用户权益服务
func UserPowerBusiness(appId string, userId string, ResourceId,ResourceType string,agentType int) *userPowerBusiness {
	return &userPowerBusiness{
		appId,
		userId,
		agentType,
		service.AvailableService{AppId: appId, UserId: userId},
		ResourceId,
		ResourceType,
	}
}

// 直播权益判断
func (upb *userPowerBusiness) IsHaveAlivePower(resouceId string, resouceType string, needExpire bool) (string, bool) {
	result, err := upb.IsEncryAliveAccess(resouceId)
	if err != nil {
		logging.Error(err)
	}
	return "", result
	/*if redis_gray.InGrayShopSpecial("is_switch_new_permission", upb.AppId) {
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
	}*/
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

// NewAliveAccess 新权益服务
func NewAliveAccess(upb userPowerBusiness) (bool, error) {
	return upb.service.IsResourceAccess(upb.ResourceId,false, 0)
}

// OldAliveAccess 旧权益服务(LB_CT_COPENCLOSE_IN)
//咨询bowuding(丁博武) 已不在维护且准备下架这个服务，所有权益服务迁移至 LB_PF_RIGHTS_IN
func  OldAliveAccess(upb userPowerBusiness)(available bool, err error) {
	resourceAvailable := service.ResourceAvailable{
		ResourceId:   upb.ResourceId,
		ResourceType: upb.ResourceType,
		NeedExpire:   false,
		AgentType:    upb.AgentType,
	}
	_,available,err =  upb.service.IsResourceAvailable(resourceAvailable)
	return
}

// CacheAliveAccess 今日已验证过权益,只针对非免费直播缓存
func CacheAliveAccess(upb userPowerBusiness) (bool, error){
	return alive_static.HEXISTS(fmt.Sprintf(staticAliveHashUser, time.Now().Format("2006-01-02")), upb.ResourceId+upb.UserId)
}

// AvaFunc 声明一个权益函数类型
type AvaFunc func(u userPowerBusiness)  (bool, error)


func serviceFunc(u userPowerBusiness, f AvaFunc) (bool, error) {
	return f(u) //通过调用f()实现任务
}


/**
直播权益查询
服务查询次序: 新权益服务 -> 当天Redis访问记录
权益服务 新权益 => LB_PF_RIGHTS_IN, redis => REDIS_RIGHT_IN
查询前检查当前服务是否出现异常 依次降级查询
 */
func (upb *userPowerBusiness)IsAliveAccess()  (available bool,err error) {
	//权益服务方法集合
	ServiceMethod := map[string]AvaFunc{
		"LB_PF_RIGHTS_IN" : NewAliveAccess,
		"REDIS_RIGHT_IN" : CacheAliveAccess,
	}
	//选择权益服务----------------------------------------------------------------//
	AvaService :=  []string{
		"LB_PF_RIGHTS_IN",
		"REDIS_RIGHT_IN",
	}

	for _, serviceName := range AvaService {
		available, err = serviceFunc(*upb, ServiceMethod[serviceName])
		if err != nil {
			//服务出现异常自动降级查询
			logging.Error(fmt.Sprintf("权益信息查询异常:%s:%s", err.Error(), serviceName))
			continue
		}else{
			break
		}
	}
	//----------------------------------------------------------------选择权益服务//
	return
}


