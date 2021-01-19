package app_conf

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	// "time"

	"github.com/gomodule/redigo/redis"

	"abs/models/business"
	"abs/models/sub_business"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"
)

type AppInfo struct {
	AppId string
	// BaseConf     *service.AppBaseConf
}

const (
	// Redis key
	shopInfoKey        = "app_conf_moudle_info:%s"
	shopInfoByTokenKey = "app_conf_moudle_info_token:%s"
	appMsgSwitch       = "app_msg_switch:%s"

	// 缓存时间控制(秒)
	// 短信预约总开关
	appConfSwitchCacheTime = "120"
	// 配置中心配置
	confHubInfoCacheTime = "60"
)

// 已废弃，请忽使用！！！
// 获取App店铺详情，此处是请求Model处理
func (a *AppInfo) GetAppInfo(needToken bool) (map[string]interface{}, error) {
	var (
		appInfo    *business.AppConf
		appModule  *business.AppModule
		shopConfig []*sub_business.ShopConfig
		cacheKey   string
	)
	cacheAppInfo := make(map[string]interface{})
	conn, err := redis_alive.GetLiveInteractConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if needToken {
		cacheKey = fmt.Sprintf(shopInfoByTokenKey, a.AppId)
	} else {
		cacheKey = fmt.Sprintf(shopInfoKey, a.AppId)
	}
	info, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil {
		logging.Error(err)
	} else {
		json.Unmarshal(info, &cacheAppInfo)
		return cacheAppInfo, nil
	}

	var wg = sync.WaitGroup{}
	eC := make(chan error)
	defer close(eC)
	wg.Add(3)
	go func() {
		appInfo, err = business.GetAppConfDetail(a.AppId)
		wg.Done()
		if err != nil {
			eC <- err
		}
	}()
	go func() {
		appModule, err = business.GetAppModule(a.AppId)
		wg.Done()
		if err != nil {
			eC <- err
		}
	}()
	go func() {
		shopConfig, err = sub_business.GetAppShopConfig(a.AppId)
		wg.Done()
		if err != nil {
			eC <- err
		}
	}()
	select {
	case err = <-eC:
		return nil, err
	default:
	}
	wg.Wait()
	// 封装为json适配老abs
	if err = util.StructJsonMap(appInfo, &cacheAppInfo); err != nil {
		return nil, err
	}
	if err = util.StructJsonMap(appModule, &cacheAppInfo); err != nil {
		return nil, err
	}

	// shopConfig 配置
	for _, v := range shopConfig {
		// t_shop_config  1绑定
		if v.Name == "is_force_phone" && v.Module == "h5_custom" {
			cacheAppInfo["is_force_phone"] = v.Value
		}
		// t_shop_config  1仅移动端访问
		if v.Name == "only_h5_play" && v.Module == "safe" {
			cacheAppInfo["only_h5_play"] = v.Value
		} else {
			cacheAppInfo["only_h5_play"] = 0
		}
		// t_shop_config 0-默认播放器 1-自研播放器 (默认0)
		if v.Name == "only_h5_play" && v.Module == "h5_custom" {
			cacheAppInfo["video_player_type"] = 1
		} else {
			cacheAppInfo["video_player_type"] = 0
		}
		// t_shop_config enable_web_rtc  1开启快直播功能
		if v.Name == "enable_web_rtc" && v.Module == "live" {
			cacheAppInfo["enable_web_rtc"] = v.Value
		}
	}

	value, err := json.Marshal(cacheAppInfo)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("SET", cacheKey, value, "EX", "30")
	if err != nil {
		return nil, err
	}

	return cacheAppInfo, nil
}

// 获取店铺配置的短信预约总开关
func (a *AppInfo) GetAppConfSwitchState() (info int, err error) {
	conn, _ := redis_alive.GetSubBusinessConn()
	defer conn.Close()

	cacheKey := fmt.Sprintf(appMsgSwitch, a.AppId)
	info, err = redis.Int(conn.Do("GET", cacheKey))
	if err == nil {
		return info, nil
	}

	ser := service.TemplateMsgService{AppId: a.AppId}
	data, err := ser.GetSwitchState()
	if err != nil {
		logging.Error(err)
		return
	}

	if _, err = conn.Do("SET", cacheKey, int(data["sms_state"].(float64)), "EX", appConfSwitchCacheTime); err != nil {
		logging.Error(err)
	}
	return int(data["sms_state"].(float64)), nil
}

// 【直播】统一获取配置中心配置信息
func (a *AppInfo) GetConfHubInfo() (baseConf *service.AppBaseConf, err error) {
	conn, err := redis_alive.GetSubBusinessConn()
	if err != nil {
		logging.Error(err)
	}
	defer conn.Close()

	cacheKey := fmt.Sprintf(shopInfoKey, a.AppId)
	info, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err == nil {
		confHubInfo := service.ConfHubInfo{}
		util.JsonDecode(info, &confHubInfo)
		baseConf = a.handleConfResult(confHubInfo)
		return
	}

	// 获取配置服务
	conInfo := service.ConfHubServer{AppId: a.AppId, WxAppType: 1}
	result, err := conInfo.GetConf([]string{"base", "version", "profit", "switches", "extra", "h5_custom", "safe", "pc", "domain"})
	if err != nil {
		return
	}
	// 组装基本配置
	baseConf = a.handleConfResult(result)

	// 缓存
	if confHubInfoBytes, err := util.JsonEncode(result); err == nil {
		if _, err = conn.Do("SET", cacheKey, confHubInfoBytes, "EX", confHubInfoCacheTime); err != nil {
			logging.Error(err)
		}
	}

	// 缓存失败不能影响返回
	return baseConf, nil
}

func (a *AppInfo) handleConfResult(result service.ConfHubInfo) *service.AppBaseConf {
	// 获取返回配置
	baseConfMap, switcheConfMap := result.Base, result.Switches
	// 添加基本配置
	baseConf := &service.AppBaseConf{
		ShopId:      baseConfMap["shop_id"].(string),
		ShopName:    baseConfMap["shop_name"].(string),
		ShopLogo:    baseConfMap["shop_logo"].(string),
		FooterLogo:  baseConfMap["footer_logo"].(string),
		Profile:     baseConfMap["profile"].(string),
		VersionType: int(result.Version["version_type"].(float64)),
		ExpireTime:  result.Version["expire_time"].(string),
		// 添加开关配置
		HasReward:           int(switcheConfMap["has_reward"].(float64)),
		HasInvite:           int(switcheConfMap["has_invite"].(float64)),
		AuthenticState:      int(switcheConfMap["authentic_state"].(float64)),
		IsShowResourcecount: int(switcheConfMap["is_show_resourcecount"].(float64)),
		RelateSellInfo:      int(result.H5Custom["relate_sell_info"].(float64)),
		VideoPlayerType:     int(result.H5Custom["video_player_type"].(float64)),
		// 新增在这里加
		PcCustomDomain: result.Domain["pc_custom_domain"].(string),
		IsEnable:       int(result.Pc["is_enable"].(float64)),
		// 是否只有H5观看
		OnlyH5Play: int(result.Safe["only_h5_play"].(float64)),
		// Profit数据
		Profit: result.Profit,
	}

	// 特殊的恶心返回兼容
	switch result.Extra["caption_define"].(type) {
	case string:
		baseConf.CaptionDefine = result.Extra["caption_define"].(string)
	default:
		baseConf.CaptionDefine = ""
	}

	return baseConf
}

// 获取云通信配置
func (a *AppInfo) GetCommunicationCloudInfo(identifier string) map[string]string {
	conf := map[string]string{
		"user_sign":    "",
		"sdk_app_id":   "",
		"account_type": os.Getenv("AccountType"),
	}
	timeRestApi := service.TimeRestApi{
		SdkAppId:   os.Getenv("AliveVideoAppId"),
		Identifier: identifier,
	}
	conf["sdk_app_id"] = timeRestApi.SdkAppId
	userSig, err := timeRestApi.GenerateUserSig()
	if err != nil {
		logging.Error(err)
		return conf
	}
	conf["user_sign"] = userSig
	return conf
}
