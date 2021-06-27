package util

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	e "abs/pkg/enums"
)

const (
	TIME_LAYOUT = "2006-01-02 15:04:05"
)

// 返回项目根目录
func rootPath() string {
	dir, _ := os.Getwd()
	return dir
}

// 获取runtime目录路径
func GetRuntimeDir() string {
	return rootPath() + "/runtime/"
}

// json方法之encode
func JsonEncode(v interface{}) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(v)
}

// json方法之decode
func JsonDecode(data []byte, v interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, v)
}

// 结构体转化为字典
func Struct2Map(obj interface{}) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Name] = field.Interface()
	}
	return data
}

// 结构体Json转化为字典
func StructJsonMap(obj interface{}, v *map[string]interface{}) error {
	jsoni := jsoniter.ConfigCompatibleWithStandardLibrary
	if obj == nil {
		return nil
	}
	buf, err := jsoni.MarshalIndent(obj, "", "    ") // 格式化编码
	if err != nil {
		return err
	}
	err = jsoni.Unmarshal([]byte(buf), v)
	if err != nil {
		return err
	}

	return nil
}

// 是否是app
func IsQyApp(versionType int) bool {
	// trainingType, _ := strconv.Atoi(versionType)
	if versionType == e.VERSION_TYPE_TRAINING_TRY || versionType == e.VERSION_TYPE_TRAINING_STD {
		return true
	}

	return false
}

// 获取H5域名
func GetH5Domain(appId string, contentPage bool) string {
	//新独立域名
	if os.Getenv("NEW_H5_DOMAIN") == "true" {
		return fmt.Sprintf("%s.%s", appId, os.Getenv("H5_DOMAIN"))
	}

	//支付域名
	if contentPage {
		return os.Getenv("H5_DOMAIN")
	}

	return fmt.Sprintf("%s/%s", os.Getenv("H5_DOMAIN"), appId)
}

/**
 * 回放视频加密方法
 */
func VideoEncrypt(str string) string {
	if str == "" {
		return ""
	}

	old := [4]string{"1", "2", "3", "4"}
	news := [4]string{"@", "#", "$", "%"}

	//base64编码
	str = SafeBase64Encode(str)
	for key, value := range old {
		str = strings.Replace(str, value, news[key], -1)
	}

	str += "__ba"

	return str
}

/**
 * 简单对称加密算法之加密
 * @param String str 需要加密的字串
 * @param String key   加密EKY
 * @return String
 */
func EncryptEncode(str, key string) string {
	strArr := strings.Split(SafeBase64Encode(str), "")
	strCount := len(strArr)

	for key, val := range strings.Split(key, "") {
		if key < strCount {
			strArr[key] += val
		}
	}

	old := [3]string{"=", "+", "/"}
	news := [3]string{"O0O0O", "o000o", "oo00o"}

	str = strings.Join(strArr, "")

	for k, v := range old {
		str = strings.Replace(str, v, news[k], -1)
	}

	return str
}

// 合并map
func MergeMap(m ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, v := range m {
		for k1, v1 := range v {
			result[k1] = v1
		}
	}
	return result
}

// 版本过期判断
func JudgeDate(versionType int, expireTime string) map[string]interface{} {
	result := map[string]interface{}{
		"type": versionType,
		"time": expireTime,
	}

	if expireTime != "0000-00-00 00:00:00" {
		expireParse, _ := time.Parse("2006-01-02 15:04:05", expireTime)
		expire := expireParse.Unix()
		if versionType == 4 {
			if time.Now().Unix() > expire {
				result["time"] = expireTime
				result["type"] = 1 //标准版已经过期
			} else if expire-time.Now().Unix() < 8*24*3600 {
				result["time"] = expireTime
				result["type"] = 2 //标准版即将过期
			}
		} else if versionType == 0 {
			if time.Now().Unix() > expire {
				result["time"] = expireTime
				result["type"] = 3 //试用版已经过期
			} else {
				result["time"] = expireTime
				result["type"] = 4 //试用版还未过期
			}
		}
	}

	return result
}
