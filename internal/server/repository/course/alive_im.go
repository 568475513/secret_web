package course

import (
	"abs/models/alive"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/cache/redis_default"
	"abs/pkg/cache/redis_gray"
	"abs/pkg/cache/redis_im"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type ImCreateRes struct {
	ActionStatus string `json:"ActionStatus"`
	ErrorInfo    string `json:"ErrorInfo"`
	ErrorCode    int    `json:"ErrorCode"`
	GroupId      string `json:"GroupId"`
}

type ImInfoRes struct {
	ActionStatus string      `json:"ActionStatus"`
	ErrorInfo    string      `json:"ErrorInfo"`
	ErrorCode    int         `json:"ErrorCode"`
	GroupInfo    interface{} `json:"GroupInfo"`
}

const (
	imCreateGroup = "https://console.tim.qq.com/v4/group_open_http_svc/create_group?sdkappid=%d&identifier=%s&usersig=%s&random=%d&contenttype=json"

	imInfoGroup = "https://console.tim.qq.com/v4/group_open_http_svc/get_group_info?sdkappid=%d&identifier=%s&usersig=%s&random=%d&contenttype=json"

	aliveInfoKey3 = "alive:%s:%s"

	aliveInfoKey2 = "aliveInfoKey:%s%s"

	scriptSystemAliveInfoKey = "script_system_alive_info:%s:%s"
)

// 过滤旧版群组
func (a *AliveInfo) GetAliveRommId(a2 *alive.Alive) string {
	redisConn, err := redis_im.GetLiveGroupActionConn()
	if err != nil {
		logging.Error(err)
	}
	defer redisConn.Close()
	roomId := a2.RoomId
	logging.Info(roomId)
	//从新增的O端灰度获取 当前店铺是否在灰度名单
	if redis_gray.InGrayShopSpecialHit("alive_im_gray_encryption", a.AppId) {
		if strings.Contains(roomId, "XET#") {
			newRoomId := roomId
			AliveImMiddler, err := alive.GetRoomIdByAliveId(a.AppId, a.AliveId, "old_room_id")
			if err != nil {
				logging.Error(err)
				return newRoomId
			}
			logging.Info(AliveImMiddler)
			if AliveImMiddler.OldRoomId != "" {
				err = alive.UpdateTAliveRommId(a.AppId, a.AliveId, AliveImMiddler.OldRoomId)
				logging.Info(err)
				if err != nil {
					logging.Error(err)
					return newRoomId
				}
				go a.resetRoomIdCache()
				imGroupActiveCacheKey := fmt.Sprintf(imGroupActive, a.AliveId)
				redisConn.Do("del", imGroupActiveCacheKey)
				hitImActiveCacheKey := fmt.Sprintf(hitImActive, a.AliveId[len(a.AliveId)-1:])
				redisConn.Do("zrem", hitImActiveCacheKey, a.AliveId)

				err = alive.UpdateForbidRoomId(a.AppId, roomId, AliveImMiddler.OldRoomId)
				if err != nil {
					logging.Error(err)
				}
				return fitlerOldRoom(a.AppId, a.AliveId, AliveImMiddler.OldRoomId)
			}
			return newRoomId
		}
		return fitlerOldRoom(a.AppId, a.AliveId, roomId)
	}

	imGroupActiveCacheKey := fmt.Sprintf(imGroupActive, a.AliveId)
	rid, _ := redis.String(redisConn.Do("get", imGroupActiveCacheKey))
	if rid != "" {
		roomId = rid
	}
	if strings.Contains(roomId, "XET#") {
		a.hitJudgeActive(redisConn, roomId)
		return roomId
	}
	AliveImMiddler, err := alive.GetRoomIdByAliveId(a.AppId, a.AliveId, "new_room_id")
	if err != nil {
		logging.Error(err)
		return roomId
	}
	logging.Info(AliveImMiddler)
	var newRoomId string
	if AliveImMiddler.NewRoomId != "" {
		newRoomId = AliveImMiddler.NewRoomId
	} else {
		newRoomId = a.getRandRoomId(10)
	}
	logging.Info(newRoomId)
	res := a.hitJudgeActive(redisConn, newRoomId)
	logging.Info(res)
	if res {
		err = alive.UpdateTAliveRommId(a.AppId, a.AliveId, newRoomId)
		logging.Info(err)
		if err != nil {
			logging.Error(err)
			return roomId
		}

		err = alive.UpdateForbidRoomId(a.AppId, roomId, newRoomId)
		if err != nil {
			logging.Error(err)
			return roomId
		}
		aim := alive.AliveImMiddler{
			AppId:     a.AppId,
			AliveId:   a.AliveId,
			OldRoomId: roomId,
			NewRoomId: newRoomId,
		}
		go a.resetRoomIdCache()
		logging.Info(aim)
		err = alive.InsertImMiddle(aim)
		logging.Info(aim)
		if err != nil {
			logging.Error(err)
			return roomId
		}
		return newRoomId
	}
	return roomId
}

func fitlerOldRoom(appId string, aliveId string, roomId string) string {
	if strings.Contains(roomId, "@TGS#") {
		return hitJudgeActiveOld(appId, aliveId, roomId)
	}
	return roomId
}

func hitJudgeActiveOld(appId string, aliveId string, roomId string) string {
	redisConn, err := redis_im.GetLiveGroupActionConn()
	if err != nil {
		logging.Error(err)
		return roomId
	}
	defer redisConn.Close()
	imGroupActiveCacheKey := fmt.Sprintf(oldImGroupActive, aliveId)
	rid, _ := redis.String(redisConn.Do("get", imGroupActiveCacheKey))
	if rid != "" {
		return roomId
	}

	hitImActiveCacheKey := fmt.Sprintf(oldHitImActive, aliveId[len(aliveId)-1:])
	zScoreValue, _ := redisConn.Do("zscore", hitImActiveCacheKey, aliveId)
	expire := 86400 - (time.Now().Unix()+8*3600)%86400
	if zScoreValue != nil {
		redisConn.Do("setex", imGroupActiveCacheKey, expire, roomId)
		redisConn.Do("zadd", hitImActiveCacheKey, time.Now().Unix(), aliveId)
		return roomId
	}
	groupId, err := getGroupOldRoomId(appId, aliveId, roomId)
	if err == nil {
		redisConn.Do("setex", imGroupActiveCacheKey, expire, groupId)
		redisConn.Do("zadd", hitImActiveCacheKey, time.Now().Unix(), aliveId)
		return groupId
	}
	return groupId
}

// 这abs_go日志和异常警示消息体
func getGroupOldRoomId(appId string, aliveId string, roomId string) (string, error) {
	isExist := judgeRoomIdIsExist(roomId)
	if isExist {
		return roomId, nil
	}
	// 群组id已解散，需要建立一新的，并同步入库返回
	groupId, err := createOldGroup(roomId)
	if err != nil {
		// 没成，再来一次
		groupId, err = createOldGroup(roomId)
		if err != nil {
			return groupId, err
		}
	}
	isOk := changeRoomIdData(appId, aliveId, roomId, groupId)
	if isOk {
		return groupId, nil
	}
	return roomId, nil
}

func changeRoomIdData(appId string, aliveId string, roomId string, groupId string) bool {
	err := alive.UpdateTAliveRommId(appId, aliveId, groupId)
	logging.Info(err)
	if err != nil {
		logging.Error(err)
		return false
	}

	err = alive.UpdateForbidRoomId(appId, roomId, groupId)
	if err != nil {
		logging.Error(err)
		return false
	}
	return true
}

func judgeRoomIdIsExist(roomId string) bool {
	timeRestApi := service.TimeRestApi{
		SdkAppId:   os.Getenv("AliveVideoAppId"),
		Identifier: os.Getenv("AliveVideoAdminId"),
	}
	userSig, _ := timeRestApi.GenerateUserSig()
	random := getRandInt(4294967295)
	SdkAppId, _ := strconv.Atoi(os.Getenv("AliveVideoAppId"))
	requestUrl := fmt.Sprintf(imInfoGroup, SdkAppId, timeRestApi.Identifier, userSig, random)
	GroupIdList := []string{roomId}
	requestData := map[string][]string{
		"GroupIdList": GroupIdList,
	}
	logging.Info(userSig)
	logging.Info(requestUrl)
	requestDataJson, _ := util.JsonEncode(requestData)
	var responseMap ImInfoRes
	request := service.Post(requestUrl)
	request.SetParams(requestDataJson)
	request.SetTimeout(1000 * time.Millisecond)
	err := request.ToJSON(&responseMap)
	logging.Info(responseMap)
	if err != nil {
		return false
	}
	if responseMap.ErrorCode == 0 {
		return true
	}
	return false
}

func (a *AliveInfo) resetRoomIdCache() {

	defaultConn, _ := redis_default.GetLiveInfoConn()
	defer defaultConn.Close()
	scriptCacheKey := fmt.Sprintf(scriptSystemAliveInfoKey, a.AppId, a.AliveId)
	_, err := defaultConn.Do("del", scriptCacheKey)
	if err != nil {
		logging.Error(err)
	}

	aliveBusinessConn, _ := redis_alive.GetLiveBusinessConn()
	defer aliveBusinessConn.Close()
	cacheKey := fmt.Sprintf(aliveInfoKey, a.AppId, a.AliveId)
	_, err = aliveBusinessConn.Do("del", cacheKey)
	if err != nil {
		logging.Error(err)
	}
	cacheKey = fmt.Sprintf(aliveInfoKey2, a.AppId, a.AliveId)
	_, err = aliveBusinessConn.Do("del", cacheKey)
	if err != nil {
		logging.Error(err)
	}
	cacheKey = fmt.Sprintf(aliveInfoKey3, a.AppId, a.AliveId)
	_, err = aliveBusinessConn.Do("del", cacheKey)
	if err != nil {
		logging.Error(err)
	}

	return
}

//检测缓存
func (a *AliveInfo) hitJudgeActive(redisConn redis.Conn, room_id string) bool {

	imGroupActiveCacheKey := fmt.Sprintf(imGroupActive, a.AliveId)
	exists, _ := redis.Bool(redisConn.Do("exists", imGroupActiveCacheKey))
	if exists {
		return false
	}

	hitImActiveCacheKey := fmt.Sprintf(hitImActive, a.AliveId[len(a.AliveId)-1:])
	logging.Info(hitImActiveCacheKey)
	zScoreValue, _ := redisConn.Do("zscore", hitImActiveCacheKey, a.AliveId)
	expire := 86400 - (time.Now().Unix()+8*3600)%86400
	logging.Info(zScoreValue)
	if zScoreValue != nil {
		redisConn.Do("setex", imGroupActiveCacheKey, expire, room_id)
		redisConn.Do("zadd", hitImActiveCacheKey, time.Now().Unix(), a.AliveId)
		return false
	}
	res := createGroup(room_id)
	logging.Info(res)
	if res {
		redisConn.Do("setex", imGroupActiveCacheKey, expire, room_id)
		redisConn.Do("zadd", hitImActiveCacheKey, time.Now().Unix(), a.AliveId)
		return true
	}
	return false
}

//获取随机数
func getRandInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

// 创建群组
func createGroup(GroupId string) bool {
	SdkAppId, _ := strconv.Atoi(os.Getenv("WHITE_BOARD_SDK_APP_ID"))
	key := os.Getenv("WHITE_BOARD_SECRET_KEY")
	identifier := os.Getenv("WHITE_BOARD_ID")
	userSig, _ := tencentyun.GenUserSig(SdkAppId, key, identifier, 86400*180)
	random := getRandInt(4294967295)
	requestUrl := fmt.Sprintf(imCreateGroup, SdkAppId, identifier, userSig, random)
	requestData := map[string]string{
		"Owner_Account": identifier,
		"Type":          "AVChatRoom",
		"GroupId":       GroupId,
		"Name":          "TestGroup",
	}
	requestDataJson, _ := util.JsonEncode(requestData)
	var responseMap ImCreateRes
	request := service.Post(requestUrl)
	fmt.Println(requestData)
	request.SetParams(requestDataJson)
	request.SetTimeout(1000 * time.Millisecond)
	err := request.ToJSON(&responseMap)
	logging.Info(responseMap)
	ErrprCodes := map[int]string{
		0:     "无错误。",
		10021: "群组 ID 已被使用，请选择其他的群组 ID。",
		10025: "群组 ID 已被使用，并且操作者为群主，可以直接使用。",
	}
	if err != nil {
		return false
	}
	if _, ok := ErrprCodes[responseMap.ErrorCode]; ok {
		return true
	}
	return false
}

// 创建旧群组
func createOldGroup(GroupId string) (string, error) {
	timeRestApi := service.TimeRestApi{
		SdkAppId:   os.Getenv("AliveVideoAppId"),
		Identifier: os.Getenv("AliveVideoAdminId"),
	}
	userSig, _ := timeRestApi.GenerateUserSig()
	random := getRandInt(4294967295)
	SdkAppId, _ := strconv.Atoi(os.Getenv("WHITE_BOARD_SDK_APP_ID"))
	requestUrl := fmt.Sprintf(imCreateGroup, SdkAppId, timeRestApi.Identifier, userSig, random)
	requestData := map[string]string{
		"Owner_Account": timeRestApi.Identifier,
		"Type":          "AVChatRoom",
		"Name":          "TestGroup",
	}
	requestDataJson, _ := util.JsonEncode(requestData)
	var responseMap ImCreateRes
	request := service.Post(requestUrl)
	logging.Info(requestData)
	request.SetParams(requestDataJson)
	request.SetTimeout(1000 * time.Millisecond)
	err := request.ToJSON(&responseMap)
	logging.Info(responseMap)
	if err != nil {
		logging.Error(err)
		return GroupId, err
	}
	if responseMap.ErrorCode == 0 {
		return responseMap.GroupId, nil
	}
	return GroupId, err
}

// 获取随机不重复群组id
func (a *AliveInfo) getRandRoomId(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
		fmt.Println(string(byte(b)))
	}
	return "XET#" + string(bytes)
}
