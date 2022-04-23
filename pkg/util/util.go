package util

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
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
	DATE_LAYOUT = "2006-01-02"
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

// 是否是鹅直播店铺
func IsEliveApp(versionType int) bool {
	if versionType == e.VERSION_TYPE_ELIVE {
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

//时间字符串转换为time
func StringToTime(timeStr string, timeZone string) time.Time {
	loc, _ := time.LoadLocation(timeZone)
	theTime, _ := time.ParseInLocation(TIME_LAYOUT, timeStr, loc)
	return theTime
}

//[]string去重
func DuplicateRemovalByArrString(s []string) []string {
	var result []string // 存放结果
	for i := range s {
		flag := true
		for j := range result {
			if s[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, s[i])
		}
	}
	return result
}

func GetEncryptUserId(userId string) (encryptUserId string) {
	origData := []byte(userId)
	keyStr := fmt.Sprintf("%x", md5.Sum([]byte("xiaoeapp2021"))) // 密钥字符串
	//截取key前16位字符
	key16Str := keyStr[0 : len(keyStr)-16]
	key := []byte(key16Str) //加密字符串
	encrypted := AesEncryptECB(origData, key)
	encryptUserId = base64.StdEncoding.EncodeToString(encrypted)
	return
}

func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func GetPrice(price float64) (prices float64) {
	switch true {
	case price >= float64(80) && price <= float64(85):
		prices = price + 0.005
		break
	case price > float64(85) && price <= float64(90):
		prices = price + 0.002
		break
	case price > float64(90) && price <= float64(95):
		prices = price + 0.001
		break
	case price > float64(95) && price < float64(99):
		prices = price + 0.0002
		break
	default:
		prices = price
		break
	}
	return
}
