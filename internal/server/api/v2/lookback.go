package v2

import (
	"abs/pkg/cache/redis_gray"
	"abs/pkg/logging"
	"abs/service"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/url"
	_ "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"abs/internal/server/repository/course"
	"abs/internal/server/repository/material"
	"abs/internal/server/rules/validator"
	"abs/pkg/app"
	"abs/pkg/enums"
)

/**
 * 获取直播回放链接接口
 */
func GetLookBack(c *gin.Context) {
	var (
		err error
		req validator.LookBackRuleV2
	)

	// 参数校验
	AppId := app.GetAppId(c)
	if err = app.ParseQueryRequest(c, &req); err != nil {
		return
	}
	// req.AliveId = c.Query("alive_id")
	// req.Client, err = strconv.Atoi(c.Query("client"))
	if req.Client == 0 {
		// 默认公众号
		req.Client = 1
	}
	// if AppId == "" || req.AliveId == "" {
	// 	app.FailWithMessage("内容已被删除", enums.Code_Db_Not_Find, c)
	// 	return
	// }

	//获取直播数据
	aliveRep := course.AliveInfo{AppId: AppId, AliveId: req.AliveId}
	aliveInfo, err := aliveRep.GetAliveInfo()
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("GetAliveInfo错误:%s", err.Error()), enums.ERROR, c)
		return
	}

	//获取直播状态
	aliveState := aliveRep.GetAliveLookBackStates(aliveInfo)

	//获取直播回放链接
	lookBackRep := material.LookBack{AppId: AppId, AliveId: req.AliveId}
	data := lookBackRep.GetLookBackUrl(aliveInfo, aliveState, req.Client)

	app.OkWithData(data, c)
}

func defenceDownload(appId, url string) string {
	if redis_gray.InGrayShopSpecialHit("lookback_video_encrypt_gray", appId) {
		conInfo := service.ConfHubServer{AppId: appId, WxAppType: 1}
		result, err := conInfo.GetConf([]string{"base", "safe"})
		if err != nil {
			logging.Error(err)
			return url
		}
		var whref string
		if result.Safe["is_video_encrypt"].(bool) {
			whref = "*.xiaoe-tech.com,*.xiaoeknow.com"
		} else {
			whref = ""
		}
		t := time.Now().AddDate(0, 0, 1).Unix()
		exper := "0"
		replaceUrl := GetSignByVideoUrl(url, whref, strconv.FormatInt(t, 16), exper)
		return replaceUrl
	}
	return url
}

func GetSignByVideoUrl(urlPath, whref, t, exper string) string {
	randStr := GetRandomLen(12)
	key := os.Getenv("QCLOUD_VOD_ENCRYPT_KEY")
	u, _ := url.Parse(urlPath)
	dir := filepath.Dir(u.Path)
	sign := md5.Sum([]byte(key + dir + t + exper + randStr + whref))
	whrefEn := url.QueryEscape(whref)
	baseUrl := os.Getenv("QCLOUD_VOD_MAIN_URL")
	keyUrl := os.Getenv("QCLOUD_VOD_ENCRYPT_KEY_URL2")
	replaceUrl := strings.Replace(urlPath, baseUrl, keyUrl, 1)
	if strings.Contains(replaceUrl, "http") {
		replaceUrl = strings.Replace(replaceUrl, "http", "https", 1)
	}
	signStr := (*string)(unsafe.Pointer(&sign))
	if whrefEn == "" {
		return replaceUrl + "?t=" + t + "&exper=" + exper + "&us=" + randStr + "&sign=" + *signStr
	}
	return replaceUrl + "?t=" + t + "&exper=" + exper + "&us=" + randStr + "&whref=" + whrefEn + "&sign=" + *signStr
}

func GetRandomLen(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
		fmt.Println(string(byte(b)))
	}
	return string(bytes)
}
