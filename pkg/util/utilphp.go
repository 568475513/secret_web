package util

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	e "abs/pkg/enums"
)

// 此文件方法皆为适配php或者abs特殊逻辑，以及接口返回的特殊转化方法！！！

// 转化字典里面的sqlNull[几乎废弃]
func MapNull2String(val map[string]interface{}) map[string]interface{} {
	for k, v := range val {
		switch v.(type) {
		case map[string]interface{}:
			if data, ok := v.(map[string]interface{})["String"]; ok {
				if valid, ok := v.(map[string]interface{})["Valid"]; ok && valid.(bool) {
					val[k] = nil
				} else {
					val[k] = data.(string)
				}
			}
		}
	}

	return val
}

// 兼容小程序样式的替换
func StyleReplace(val string) string {
	val = strings.Replace(val, "&quot;", "", -1)
	val = strings.Replace(val, "&lt;", "<", -1)
	val = strings.Replace(val, "&gt;", ">", -1)
	return val
}

// 兼容小程序崩溃字符的替换
func ErrorReplace(val string) string {
	val = strings.ReplaceAll(val, "\u2028", "")
	val = strings.ReplaceAll(val, "\u2029", "")
	return val
}

type ContentParam struct {
	Type         int    `json:"type"`
	ResourceType int    `json:"resource_type"`
	ResourceId   string `json:"resource_id"`
	ProductId    string `json:"product_id"`
	PaymentType  int    `json:"payment_type"`
	ChannelId    string `json:"channel_id"`
	AppId        string `json:"app_id"`
	Source       string `json:"source"`
	Scene        string `json:"scene"`
	ContentAppId string `json:"content_app_id"`
	ShareUserId  string `json:"share_user_id"`
	ShareType    int    `json:"share_type"`
	ShareAgent   string `json:"share_agent"`
	ShareFrom    string `json:"share_from"`
	ExtraData    int    `json:"extra_data"`
}

// 老的跳转拼接方法
// type int,
// resourceType int,
// resourceId string,
// productId string,
// channelId string,
// appId string,
// source string,
// content_app_id string,
func ContentUrl(param ContentParam) string {
	return fmt.Sprintf("/content_page/%s", SafeBase64Encode(param))
}

// 接上
func ParentColumnsUrl(param ContentParam) string {
	return fmt.Sprintf("/parent_columns/%s", SafeBase64Encode(param))
}

// 输入参数获取包装有参数的字符串
func PutParmToStr(tempParam map[string]interface{}) (strBase64 string, err error) {
	tempParamJson, err := JsonEncode(tempParam)
	if err != nil {
		return
	}

	return SafeBase64Encode(tempParamJson), nil
}

// 获取直播间地址
func GetAliveRoomUrl(resourceId string, productId string, channelId string, appId string, extraData int) string {
	params := map[string]interface{}{
		"type":          e.PaymentTypeReward,
		"resource_type": e.ResourceTypeLive,
		"resource_id":   resourceId,
		"product_id":    productId,
		"channel_id":    channelId,
		"app_id":        appId,
		"extra_data":    extraData,
	}
	return fmt.Sprintf("/content_page/%s", SafeBase64Encode(params))
}

// Base64封装
func SafeBase64Encode(v interface{}) string {
	var paramsJsonStr []byte
	switch v.(type) {
	case string:
		paramsJsonStr = []byte(v.(string))
	case []byte:
		paramsJsonStr = v.([]byte)
	default:
		paramsJsonStr, _ = JsonEncode(v)
	}
	baseStr := base64.StdEncoding.EncodeToString(paramsJsonStr)
	baseStr = strings.Replace(baseStr, "+", "-", -1)
	baseStr = strings.Replace(baseStr, "/", "_", -1)
	baseStr = strings.Replace(baseStr, "=", "", -1)
	return baseStr
}

// php项目老方法
func UrlSafeB64Decode(str string) []byte {
	str = strings.Replace(str, "-", "+", -1)
	str = strings.Replace(str, "_", "/", -1)
	mod4 := len(str) % 4
	if mod4 > 0 {
		str = str + Substr("====", mod4, 0)
	}
	resBytes, _ := base64.StdEncoding.DecodeString(str)
	return resBytes
}

// Substr substr()
func Substr(str string, start int, length int) string {
	if start < 0 || length < -1 {
		return str
	}
	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}
	end := int(start) + length
	if end > len(str) {
		end = len(str)
	}
	return str[start:end]
}

//取系统内url跳转时的全url(不包括跳转支付页)
func UrlWrapper(path, currentUrl, appId string) string {
	if Substr(path, 0, 7) == "http://" || Substr(path, 0, 8) == "https://" {
		return path
	}
	if strings.Contains(path, "javascript:") {
		return path
	}
	if Substr(path, 0, 1) == "/" {
		path = Substr(path, 1, utf8.RuneCountInString(path)-1)
	}
	//从当前链接提取域名和协议
	var u, err = url.Parse(currentUrl)
	if err != nil {
		return path
	}
	host := u.Hostname()
	protocol := "http://"
	if Substr(protocol, 0, 8) == "https://" {
		protocol = "https://"
	}
	urlType := ParseUrlType(currentUrl)
	if urlType == 1 || strings.Contains(path, "content_page") {
		return fmt.Sprintf("%s%s/%s", protocol, host, path)
	}
	return fmt.Sprintf("%s%s%s/%s", protocol, host, appId, path)
}

//解析url类型
//-1对应原abs此方法的null，可能是负载心跳
//0、解析异常
//1、老的所有号的url格式
//2、新的普通号url格式(普通页面)
//3、新的普通号url格式(支付页面)
//4、前端监控上报
//5、邮件模块请求
//6、其他普通请求<url里面不需要任何id信息>
func ParseUrlType(url string) int {
	var urlType int
	if strings.Contains(url, "/e_watcher/") {
		urlType = 4
	} else {
		if strings.Contains(url, "/e_mail/") {
			urlType = 5
		} else {
			if strings.Contains(url, "/charge_ask") {
				urlType = 6
			} else {
				if strings.Contains(url, "localhost") {
					urlType = -1
				} else {
					if strings.Contains(url, ".h5.") {
						urlType = 1
					} else {
						if strings.Contains(url, "content_page") {
							urlType = 3
						} else {
							urlType = 2
						}
					}
				}
			}
		}
	}
	return urlType
}

// 获取小程序版本 1-不是小程序 2-小程序苹果版 3-小程序安卓版
func GetMiniProgramVersion(client, userAgent string) (version int) {
	version = 1
	if client == strconv.Itoa(e.AGENT_TYPE_APP) && userAgent != "" {
		agent := strings.ToLower(userAgent)
		if strings.Contains(agent, "iphone") || strings.Contains(agent, "ios") {
			version = 2 // 苹果小程序
		} else {
			version = 3 // 安卓小程序
		}
	}
	return
}
