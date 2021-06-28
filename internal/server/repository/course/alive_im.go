package course

import (
	"abs/models/alive"
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

const (
	imCreateGroup = "https://console.tim.qq.com/v4/group_open_http_svc/create_group?sdkappid=%d&identifier=%s&usersig=%s&random=%d&contenttype=json"
)

// 过滤旧版群组
func (a *AliveInfo) GetAliveRommId(a2 *alive.Alive) string {
	redisConn, err := redis_im.GetLiveGroupActionConn()
	if err != nil {
		logging.Error(err)
	}
	defer redisConn.Close()
	room_id := a2.RoomId
	imGroupActiveCacheKey := fmt.Sprintf(imGroupActive, a.AliveId)
	rid, _ := redis.String(redisConn.Do("get", imGroupActiveCacheKey))
	if rid != "" {
		room_id = rid
	}
	if strings.Contains(room_id, "XET#") {
		a.hitJudgeActive(redisConn, room_id)
		return room_id
	}
	AliveImMiddler, err := alive.GetRoomIdByAliveId(a.AppId, a.AliveId)
	if err != nil {
		logging.Error(err)
		return room_id
	}
	var newRoomId string
	if AliveImMiddler.NewRoomId != "" {
		newRoomId = AliveImMiddler.NewRoomId
	} else {
		newRoomId = a.getRandRoomId(10)
	}
	res := a.hitJudgeActive(redisConn, newRoomId)
	if res {
		err = alive.UpdateTAliveRommId(a.AppId, a.AliveId, newRoomId)
		if err != nil {
			logging.Error(err)
		}
		err = alive.UpdateForbidRoomId(a.AppId, room_id, newRoomId)
		if err != nil {
			logging.Error(err)
		}
		aim := alive.AliveImMiddler{
			AppId:     a.AppId,
			AliveId:   a.AliveId,
			OldRoomId: room_id,
			NewRoomId: newRoomId,
		}
		err = alive.InsertImMiddle(aim)
		if err != nil {
			logging.Error(err)
		}
	}
	return room_id
}

//检测缓存
func (a *AliveInfo) hitJudgeActive(redisConn redis.Conn, room_id string) bool {

	imGroupActiveCacheKey := fmt.Sprintf(imGroupActive, a.AliveId)
	exists, _ := redis.Bool(redisConn.Do("exists", imGroupActiveCacheKey))
	if exists {
		return false
	}

	hitImActiveCacheKey := fmt.Sprintf(hitImActive, a.AliveId[len(a.AliveId)-1:])
	zScoreValue, _ := redisConn.Do("zscore", hitImActiveCacheKey, a.AliveId)
	expire := 86400 - (time.Now().Unix()+8*3600)%86400
	if zScoreValue != nil {
		redisConn.Do("setex", imGroupActiveCacheKey, expire, room_id)
		redisConn.Do("zadd", hitImActiveCacheKey, time.Now().Unix(), a.AliveId)
		return false
	}
	res := createGroup(room_id)
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
		"GroupId":       "12445544",
		"Name":          "TestGroup",
	}
	requestDataJson, _ := util.JsonEncode(requestData)
	var responseMap ImCreateRes
	request := service.Post(requestUrl)
	fmt.Println(requestData)
	request.SetParams(requestDataJson)
	request.SetTimeout(1000 * time.Millisecond)
	err := request.ToJSON(&responseMap)
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
