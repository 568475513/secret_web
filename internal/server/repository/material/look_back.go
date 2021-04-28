package material

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"abs/models/alive"
	"abs/models/business"
	"abs/models/sub_business"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/cache/redis_gray"
	e "abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"
)

const (
	aliveLookBackKey     = "alive_look_back_new:%s:%s"
	LookBackEncryptkey   = "alive_lookback_encryption"
	aliveOnlyDrmMp4Key   = "alive_only_drm_mp4"
	xiaoeVideoEncryptKey = "xiaoe_video_encrypt_whirly_is_a_liangzai"
	aliveLookbackTimeKey = "alive_lookback_time_%s_%s"

	// 缓存时间控制(秒)
	// 直播剪辑表的数据
	lookBackFileCacheTime = "10"
	// 课程设置的回放过期时间
	lookbackExpireCacheTime = "60"
)

type LookBack struct {
	AppId   string
	AliveId string
}

type lookBackTime struct {
	ExpireTime string `json:"expire_time"`
	Expire     string `json:"expire"`
	ExpireType int    `json:"expire_type"`
}

// 获取直播结束后的回放视频链接等内容
func (lb *LookBack) GetLookBackUrl(aliveInfo *alive.Alive, aliveState, appType int) map[string]string {
	data := make(map[string]string)
	aliveVideoUrl := ""
	miniAliveVideoUrl := ""
	aliveVideoMp4Url := ""
	aliveReviewUrl := "/" + lb.AliveId + ".m3u8"
	aliveVideoUrlEncrypt := ""

	if aliveState == 3 { //直播已结束
		aliveVideoUrl = fmt.Sprintf("https://%s/%s.m3u8", util.GetH5Domain(lb.AppId, false), lb.AliveId)
		miniAliveVideoUrl = aliveVideoUrl
		if aliveInfo.IsLookback == 0 { //没有开启回放
			aliveVideoUrl = ""
			miniAliveVideoUrl = ""
		} else {
			if aliveInfo.AliveType == e.AliveTypeVideo { //语音直播
				videoMiddleTranscode, err := business.GetVideoMiddleTranscode(aliveInfo.FileId)
				if err != nil {
					aliveVideoUrl, miniAliveVideoUrl = "", ""
				}
				if videoMiddleTranscode != nil && videoMiddleTranscode.VideoHls != "" {
					aliveVideoUrl = videoMiddleTranscode.VideoHls
					miniAliveVideoUrl = videoMiddleTranscode.VideoHls
				}
			} else if aliveInfo.AliveType == e.AliveTypePush || aliveInfo.AliveType == e.AliveOldTypePush { //视频直播
				lookBackFile, _ := lb.GetLookBackFile(lb.AppId, lb.AliveId)
				if aliveInfo.CreateMode == 1 { // 转播课程，判断直播方的回放权限
					originAliveInfo, _ := alive.GetAliveInfoByChannelId(aliveInfo.ChannelId, []string{"app_id", "id", "is_lookback"})
					if originAliveInfo != nil {
						permission, err := sub_business.GetCloneResApply(originAliveInfo.AppId, originAliveInfo.Id, aliveInfo.AppId, []string{"lookback_permission"})
						if err != nil {
							logging.Error(err)
						} else if permission != nil && originAliveInfo.IsLookback == 1 && permission.LookbackPermission == 1 { // 原视频有开启回放且开了权限
							lookBackFile, _ = lb.GetLookBackFile(originAliveInfo.AppId, originAliveInfo.Id)
						}
					}
				}

				// 如果存在回看文件的记录 直播推流才有数据并且转码拼接成功
				if lookBackFile != nil && lookBackFile.AliveId != "" {
					aliveVideoUrl = lookBackFile.LookbackM3u8
					miniAliveVideoUrl = lookBackFile.LookbackM3u8
					aliveVideoMp4Url = lookBackFile.LookbackMp4
				} else { //没有走原逻辑
					aliveVideoUrlOrigin, miniAliveVideoUrlOrigin, aliveVideoMp4UrlOrigin, _ := lb.GetAliveComposeLookBack(aliveInfo)
					if aliveVideoUrlOrigin != "" {
						aliveVideoUrl = aliveVideoUrlOrigin
					}
					if miniAliveVideoUrlOrigin != "" {
						miniAliveVideoUrl = miniAliveVideoUrlOrigin
					}
					aliveVideoMp4Url = aliveVideoMp4UrlOrigin
				}
			}
		}

		chanceMp4 := redis_gray.InGrayShop(aliveOnlyDrmMp4Key, lb.AppId)
		if chanceMp4 {
			if strings.Index(aliveVideoUrl, "/drm/") == -1 {
				aliveVideoUrl = ""
			}
			if strings.Index(miniAliveVideoUrl, "/drm/") == -1 {
				miniAliveVideoUrl = ""
			}
		}

		//如果在o端名单内则将字段加密
		if redis_gray.InGrayShop(LookBackEncryptkey, lb.AppId) {
			//加上app_type判断，appType == 1 公众号
			if appType == 1 && aliveVideoUrl != "" && strings.Index(aliveVideoUrl, ".mp4") == -1 &&
				(aliveInfo.AliveType == 2 || aliveInfo.AliveType == 4) {
				miniAliveVideoUrl = util.VideoEncrypt(aliveVideoUrl)
				aliveReviewUrl = miniAliveVideoUrl
				aliveVideoUrlEncrypt = util.VideoEncrypt("https://" + util.GetH5Domain(lb.AppId, false) +
					"/video_encrypt/index?m3u8=" + url.QueryEscape(util.EncryptEncode(aliveVideoUrl, xiaoeVideoEncryptKey)))
				aliveVideoUrl = miniAliveVideoUrl
			}
		}
	}

	data["aliveVideoUrl"] = aliveVideoUrl
	data["miniAliveVideoUrl"] = miniAliveVideoUrl
	data["aliveReviewUrl"] = aliveReviewUrl
	data["aliveVideoUrlEncrypt"] = aliveVideoUrlEncrypt
	data["aliveVideoMp4Url"] = aliveVideoMp4Url
	data = lb.ReplaceLookBackUrl(data)

	return data
}

//替换url为素材中心的url
func (lb *LookBack) ReplaceLookBackUrl(data map[string]string) map[string]string {
	//过滤数据，只去素材中心查满足以下正则匹配的url
	var requestParam []string
	for _, value := range data {
		match, _ := regexp.MatchString("https?:\\\\?\\/\\\\?\\/([0-9a-z\\-_]+?\\.[a-z]+|([0-9a-z\\-_]+)\\.vod2?)\\.myqcloud\\.com[^\"\\'\\s]+\\.(mp3|mp4|m3u8|epub|opf|pdf|m4a)", value)
		if match == true {
			requestParam = append(requestParam, value)
		}
	}
	responseData, err := service.WashingData(lb.AppId, requestParam)
	if err != nil {
		logging.Error(err)
	}

	//替换url
	if len(responseData.FilterData) > 0 {
		for key, value := range data {
			url, ok := responseData.FilterData[value]
			if ok {
				data[key] = url
			}
		}
	}

	return data
}

/**
 * 获取直播剪辑表的数据
 */
func (lb *LookBack) GetLookBackFile(appId string, aliveId string) (*alive.AliveLookBack, error) {
	var (
		err                error
		cacheAliveLookBack *alive.AliveLookBack
	)

	conn, err := redis_alive.GetLiveBusinessConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cacheKey := fmt.Sprintf(aliveLookBackKey, lb.AppId, lb.AliveId)
	info, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil {
		logging.Warn(err)
	} else {
		json.Unmarshal(info, &cacheAliveLookBack)
		return cacheAliveLookBack, nil
	}

	cacheAliveLookBack, err = alive.GetAliveLookBackFile(appId, aliveId, []string{
		"app_id",
		"alive_id",
		"lookback_file_id",
		"region_file_id",
		"lookback_mp4",
		"lookback_m3u8",
		"file_name",
		"transcode_state",
		"state",
		"origin_type"})
	if err != nil {
		return nil, err
	}

	value, err := json.Marshal(cacheAliveLookBack)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("SET", cacheKey, value, "EX", lookBackFileCacheTime)
	if err != nil {
		return nil, err
	}

	return cacheAliveLookBack, nil
}

/**
 * 获取t_alive_concat_hls_result表中的转码视频链接
 */
func (lb *LookBack) GetAliveComposeLookBack(aliveInfo *alive.Alive) (aVideoUrl string, miniAVideoUrl string, aVideoMp4Url string, err error) {

	aliveVideoUrl := ""
	miniAliveVideoUrl := ""
	aliveVideoMp4Url := ""
	if aliveInfo.ChannelId == "" {
		return aliveVideoUrl, miniAliveVideoUrl, aliveVideoMp4Url, nil
	}

	AliveHlsResult, err := alive.GetAliveHlsResult(aliveInfo.ChannelId, []string{
		"latest_m3u8_file_id",
		"concat_latest_file_id",
		"concat_m3u8_url",
		"transcode_state",
		"transcode_success_last_time",
		"concat_success_last_time",
		"transcode_m3u8_url",
		"concat_times",
		"transcode_times",
		"compose_latest_file_id",
		"concat_mp4_url",
		"is_use_concat_mp4",
		"concat_mp4_url",
		"is_drm",
		"drm_m3u8_url"})
	if err != nil {
		return aliveVideoUrl, miniAliveVideoUrl, aliveVideoMp4Url, err
	}

	if AliveHlsResult != nil { //查看备份临时表数据
		if AliveHlsResult.ConcatM3u8Url != "" && AliveHlsResult.LatestM3u8FileId == AliveHlsResult.ConcatLatestFileId {
			/*
			 * 为了保证优化后的回看时长有保证
			 * 1、转码的id和拼接回调的m3u8id保持一致   为防止腾讯云变更规则，去除该条
			 * 2、转码状态一定是成功的
			 * 3、最后的转码成功回调一定要大于拼接成功的回调
			 * 4、转码次数和拼接次数一样
			 */
			aliveVideoUrl = AliveHlsResult.ConcatM3u8Url
			miniAliveVideoUrl = AliveHlsResult.ConcatM3u8Url
			if AliveHlsResult.TranscodeM3u8Url != "" { //转码是否完成
				if AliveHlsResult.TranscodeState == 1 &&
					AliveHlsResult.TranscodeSuccessLastTime > AliveHlsResult.ConcatSuccessLastTime &&
					AliveHlsResult.ConcatTimes > AliveHlsResult.TranscodeTimes {
					aliveVideoUrl = AliveHlsResult.TranscodeM3u8Url
					miniAliveVideoUrl = AliveHlsResult.TranscodeM3u8Url
				}
			}
		}
	}

	aliveVideoMp4Url = AliveHlsResult.ConcatMp4Url

	// 新的直播拼接方式
	if AliveHlsResult.IsUseConcatMp4 == 1 && AliveHlsResult.ConcatMp4Url != "" &&
		AliveHlsResult.ComposeLatestFileId == AliveHlsResult.LatestM3u8FileId {
		aliveVideoUrl = AliveHlsResult.ConcatMp4Url
		miniAliveVideoUrl = AliveHlsResult.ConcatMp4Url
	}

	if AliveHlsResult.IsDrm == 1 && AliveHlsResult.DrmM3u8Url != "" {
		aliveVideoUrl = AliveHlsResult.DrmM3u8Url
		miniAliveVideoUrl = AliveHlsResult.DrmM3u8Url
	}

	return aliveVideoUrl, miniAliveVideoUrl, aliveVideoMp4Url, nil
}

// 获取课程设置的回放过期时间
func (lb *LookBack) GetLookbackExpire(isLookback int, lookbackTime string) (map[string]interface{}, error) {
	conn, _ := redis_alive.GetLiveBusinessConn()
	defer conn.Close()

	cacheKey := fmt.Sprintf(aliveLookbackTimeKey, lb.AppId, lb.AliveId)
	info, err := redis.Bytes(conn.Do("GET", cacheKey))
	redisLookbackInfo := make(map[string]interface{}) //回看过期信息
	if err == nil {
		json.Unmarshal(info, &redisLookbackInfo)
		return redisLookbackInfo, nil
	}

	defaultTime := map[string]interface{}{"expire_type": e.LookBackExpireTypeNever, "expire": -1}
	redisLookbackInfo["lookback_time"] = defaultTime

	if isLookback == 1 && len(lookbackTime) > 0 {
		midd := lookBackTime{ExpireType: e.LookBackExpireTypeNever, Expire: "-1"}
		json.Unmarshal([]byte(lookbackTime), &midd)
		if midd.Expire == "0" || midd.Expire == "-1" {
			redisLookbackInfo["lookback_time"] = defaultTime
		} else {
			expireTime, _ := time.Parse("2006-01-02", midd.Expire)
			// defaultTime["expire"] = expireTime.Unix() + 86399 // 当天最后一秒（老PHP写法）
			defaultTime["expire"] = time.Date(expireTime.Year(), expireTime.Month(), expireTime.Day(), 23, 59, 59, 0, expireTime.Location()).Unix() - 28800 // UTC -> CST
			defaultTime["expire_type"] = midd.ExpireType
			redisLookbackInfo["lookback_time"] = defaultTime
		}
	}

	//设置缓存
	info, _ = json.Marshal(redisLookbackInfo)
	_, err = conn.Do("SET", cacheKey, info, "EX", lookbackExpireCacheTime)
	if err != nil {
		return nil, err
	}

	return redisLookbackInfo, nil
}
