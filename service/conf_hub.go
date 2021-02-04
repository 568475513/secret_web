package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"abs/pkg/app"
	e "abs/pkg/enums"
)

type ConfHubServer struct {
	AppId     string
	WxAppType int
}

const (
	// 配置接口的模块文档：http://doc.xiaoeknow.com/web/#/391?page_id=6610
	configModuleGet = "api/xe.shop.config.module.get/2.0.0"
	// 超时设置ms
	confGetConfTimeOut = 1000
)

type ConfHubInfoResp struct {
	app.Response
	Data ConfHubInfo `json:"data"`
}

// 返回模块控制
type ConfHubInfo struct {
	Base     map[string]interface{} `json:"base"`
	Version  map[string]interface{} `json:"version"`
	Profit   map[string]interface{} `json:"profit"`
	Switches map[string]interface{} `json:"switches"`
	Extra    map[string]interface{} `json:"extra"`
	H5Custom map[string]interface{} `json:"h5_custom"`
	Safe     map[string]interface{} `json:"safe"`
	Pc       map[string]interface{} `json:"pc"`
	Domain   map[string]interface{} `json:"domain"`
	Live     map[string]interface{} `json:"live"`
}

// 店铺基本配置
type AppBaseConf struct {
	ShopId     string `json:"app_id"`
	ShopName   string `json:"shop_name"`   // 商铺名
	ShopLogo   string `json:"shop_logo"`   // 商铺logo
	FooterLogo string `json:"footer_logo"` // 页脚logo
	Profile    string `json:"profile"`     // 简介
	// version 相关
	VersionType int    `json:"version_type"` // 版本
	ExpireTime  string `json:"expire_time"`  // 过期时间
	// extra 相关
	CaptionDefine string `json:"caption_define"` // 直播自定义文案
	// 开关相关
	HasReward           int `json:"has_reward"` // 是否有打赏功能
	HasInvite           int `json:"has_invite"` // 是否有邀请功能
	AuthenticState      int `json:"authentic_state"`
	IsShowResourcecount int `json:"is_show_resourcecount"`
	RelateSellInfo      int `json:"relate_sell_info"` // 是否显示关联售卖界面，默认1-显示，0-不显示
	OnlyH5Play          int `json:"only_h5_play"`
	VideoPlayerType     int `json:"video_player_type"` // 0-默认播放器 1-自研播放器 (默认0)
	EnableWebRtc        int `json:"enable_web_rtc"`    // 店铺快直播开关
	// profit
	Profit map[string]interface{} `json:"profit"`
	// pc
	IsEnable int `json:"is_enable"` //是否启用PC：0-否，1-是
	// domain
	PcCustomDomain string `json:"pc_custom_domain"` //pc店铺自定义域名(无schema前缀)
}

// 配置中心获取配置，传入fields中需要获取的字段。
func (c *ConfHubServer) GetConf(fields []string) (ConfHubInfo, error) {
	var result ConfHubInfoResp
	request := Post(fmt.Sprintf("%s%s", os.Getenv("LB_PF_CONFCENTER_IN"), configModuleGet))
	request.SetParams(map[string]interface{}{
		"shop_id":     c.AppId,
		"wx_app_type": c.WxAppType,
		"modules":     fields,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(confGetConfTimeOut * time.Millisecond)
	err := request.ToJSON(&result)
	if err != nil {
		return ConfHubInfo{}, err
	}
	if result.Code != e.SUCCESS {
		return ConfHubInfo{}, errors.New(fmt.Sprintf("请求配置中心获取配置错误：%s", result.Msg))
	}
	return result.Data, nil
}
