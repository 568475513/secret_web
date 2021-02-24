package user

import (
	"fmt"

	"github.com/gomodule/redigo/redis"

	"abs/models/user"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"
)

type userBusiness struct {
	AppId  string
	UserId string
}

// 用户黑名单信息
// type UserBlackInfo struct {
// 	PermissionVisit   int `json:"permission_visit"`   // 禁止访问
// 	PermissionComment int `json:"permission_comment"` // 禁止评论
// 	PermissionBuy     int `json:"permission_buy"`     // 禁止购买
// }

const (
	// Redis key
	cacheRedisUserInfoKey  = "alive_user_info:%s_%s"
	cacheRedisUserRoleKey  = "alive_user_role:%s_%s"
	cacheRedisUserCrowdKey = "alive_user_crowd:%s_%s"

	// 缓存时间
	// 请求黑名单信息
	UserBlackStateCacheTime = "30"

	// 用户类型
	Teacher uint8 = iota
	Student
)

// 工厂方法模拟构造方法，单纯得结构体太过灵活，无法对参数进行控制
func UserBusinessConstrct(appId string, userId string) *userBusiness {
	return &userBusiness{
		appId,
		userId,
	}
}

// 获取用户信息
func (business *userBusiness) GetUserInfo() (userInfo user.User, err error) {
	if business.AppId == "" && business.UserId == "" {
		return
	}
	conn, err := redis_alive.GetLiveBusinessConn()
	if err != nil {
		return
	}
	defer conn.Close()
	cacheKey := fmt.Sprintf(cacheRedisUserInfoKey, business.AppId, business.UserId)
	cacheData, err := redis.Bytes(conn.Do("GET", cacheKey))
	// 命中缓存就直接返回值
	if err == nil && cacheData != nil {
		util.JsonDecode(cacheData, &userInfo) // 错误判断先不加
		return
	}

	userService := service.UserService{AppId: business.AppId, UserId: business.UserId}
	userInfoResponse, err := userService.RequestUserInfo()
	if err != nil {
		logging.Error(err)
		return
	}

	// 指定默认头像
	userInfo = userInfoResponse.Data
	if userInfo.WxAvatar == "" {
		userInfo.WxAvatar = fmt.Sprintf("%s", "https://wechatavator-1252524126.cos.ap-shanghai.myqcloud.com/aaa/default.svg")
	}

	if jsonUserInfo, err := util.JsonEncode(userInfo); err == nil {
		if _, err = conn.Do("SET", cacheKey, jsonUserInfo, "EX", "150"); err != nil {
			logging.Warn(fmt.Sprintf("GetUserInfo Warn Redis Set Error: %s", err.Error()))
		}
	}
	return
}

// 获取用户请求黑名单信息
func (business *userBusiness) GetUserBlackStates() (blackInfo service.UserBlackInfo, err error) {
	conn, err := redis_alive.GetSubBusinessConn()
	defer conn.Close()

	cacheKey := fmt.Sprintf(cacheRedisUserCrowdKey, business.AppId, business.UserId)
	info, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err == nil {
		util.JsonDecode(info, &blackInfo)
		return
	}

	crowdReq := service.CrowdService{AppId: business.AppId, UserId: business.UserId}
	crowdResult, err := crowdReq.GetCrowdUserInfo()
	if err != nil {
		logging.Error(err)
		return
	}

	// 缓存一下，这个接口不稳定...
	if value, err := util.JsonEncode(crowdResult); err == nil {
		if _, err = conn.Do("SET", cacheKey, value, "EX", UserBlackStateCacheTime); err != nil {
			logging.Error(err)
		}
	}
	return crowdResult, nil
}