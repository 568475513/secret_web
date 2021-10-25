package marketing

import (
	"abs/pkg/cache/redis_alive"
	"abs/pkg/cache/redis_gray"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"

	"abs/models/business"
	e "abs/pkg/enums"
	"abs/pkg/logging"
)

type InviteBusiness struct {
	AppId  string
	UserId string
}

type InviteUserInfo struct {
	ShareUserId  string
	ResourceId   string
	ResourceType int
	PaymentType  int
	ProductId    string
}

// 冗余一点方便以后扩展
type InviteRelation struct {
	InviteUserInfo
	ShareType int
}

// 0 = 音频分享 1 = 日签分享 2 = 专栏分享 4-邀请卡分享
const (
	Audio        int = 0
	DaySignature int = 1
	SpecialClumn int = 2
	InvitaCard   int = 4
	NotKnow          = 5 // 5不晓得是什么鬼
	MaxTime int64 = 9999999999
)

// 更新邀请关系及邀请数量 邀请卡与推广合并后的新逻辑
func (businesss *InviteBusiness) AddInviteCountUtilsNew(inviteUserInfo InviteUserInfo) {
	if inviteUserInfo.PaymentType == e.PaymentTypeSingle {
		inviteUserInfo.ProductId = ""
	}
	if inviteUserInfo.ShareUserId == "undefined" {
		inviteUserInfo.ShareUserId = ""
	}
	inviteUserInfoPo := businesss.transformInviteUserPo(inviteUserInfo)
	if inviteUserInfo.ShareUserId == "" || businesss.UserId == "" || businesss.UserId == inviteUserInfo.ShareUserId {
		return
	}

	inviteRelation := InviteRelation{ShareType: NotKnow, InviteUserInfo: inviteUserInfo}
	inviteRelationPo := businesss.transformInviteRelationPo(inviteRelation)
	_, err := business.GetInviteUserByInvitedUser(inviteRelationPo)
	if err == nil {
		return
	}
	if err != gorm.ErrRecordNotFound {
		logging.Error(err)
		return
	}
	count := business.UpdateInviteUserByInviteCount(inviteUserInfoPo)
	if count <= 0 {
		// 通过错误来避免一次数据库查询
		inviteUserInfoPo.InviteCount = 1
		if err := business.SetInviteUser(inviteUserInfoPo); err != nil {
			logging.Error(err)
			return
		}
	}
	if err = business.SetInviteRelation(inviteRelationPo); err != nil {
		logging.Error(err)
	}

	//更新或生成排行榜
	businesss.updateOrCreateRanking(inviteUserInfo)
}

// 更新或生成排行榜
func (businesss *InviteBusiness) updateOrCreateRanking(inviteUserInfo InviteUserInfo) bool {
	if !redis_gray.InGrayShopNew("abs:alive_ic_gr:invite_gray_shop", businesss.AppId) {
		return false
	}

	conn, _ := redis_alive.GetBusinessConn()
	defer conn.Close()

	cacheKey := businesss.getRankingKeyName(businesss.AppId, inviteUserInfo.ResourceId)
	exists, err := redis.Bool(conn.Do("EXISTS", cacheKey))
	if err != nil {
		exists = false
	}
	if !exists {
		// 排行榜不存在，直接生成排行榜
		businesss.initRanking(conn, businesss.AppId, inviteUserInfo.ResourceId)
	}else {
		// 排行榜存在，更新指定用户的排名
		inviteUser,err := business.GetInviteUserByShareUser(businesss.AppId, inviteUserInfo.ShareUserId, inviteUserInfo.ResourceId)
		if err != nil {
			logging.Error(err)
			return false
		}
		err = conn.Send("ZADD", cacheKey, businesss.getRankingScore(inviteUser.InviteCount, inviteUser.CreatedAt), inviteUserInfo.ShareUserId)
		if err != nil {
			logging.Error(err)
		}
		err = conn.Send("EXPIRE", cacheKey, 3600*24*3)
		if err != nil {
			logging.Error(err)
		}
		err = conn.Flush()
		if err != nil {
			logging.Error(err)
		}
	}
	return true
}

// 初始化直播间排行榜
func (businesss *InviteBusiness) initRanking(conn redis.Conn, appId string, resourceId string) bool {
	inviteUsers, err := business.GetInviteUsersByResourceId(appId, resourceId)
	if err != nil {
		logging.Error(err)
		return false
	}
	cacheKey := businesss.getRankingKeyName(appId, resourceId)

	sTime := time.Now().Unix()

	count := 0
	for _, val := range inviteUsers {
		count++
		err = conn.Send("ZADD", cacheKey, businesss.getRankingScore(val.InviteCount, val.CreatedAt), val.ShareUserId)
		if err != nil {
			logging.Error(err)
		}
		if count % 200 == 0 {
			err = conn.Flush()
			if err != nil {
				logging.Error(err)
			}
		}
	}
	err = conn.Send("EXPIRE", cacheKey, 3600*24*3)
	if err != nil {
		logging.Error(err)
	}
	err = conn.Flush()
	if err != nil {
		logging.Error(err)
	}
	logging.Info(fmt.Sprintf("生成排行榜耗时%s秒", time.Now().Unix() - sTime))
	return true
}

// 获取排行分数
func (businesss *InviteBusiness) getRankingScore(inviteCount int, createdAt time.Time) string {
	local, _ := time.LoadLocation("Asia/Shanghai")
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Unix(createdAt.Unix(), 0).Format("2006-01-02 15:04:05"), local)
	var bt bytes.Buffer
	bt.WriteString(strconv.Itoa(inviteCount))
	bt.WriteString(strconv.FormatInt(MaxTime - stamp.Unix(),10))
	return bt.String()
}

// 获取排行榜缓存key
func (businesss *InviteBusiness) getRankingKeyName(appId string, resourceId string) string {
	return fmt.Sprintf("abs:alive_ic_grkn_%s_%s", appId, resourceId)
}

// dto模型向底层po模型转换
func (businesss *InviteBusiness) transformInviteUserPo(inviteUserInfo InviteUserInfo) business.InviteUser {
	return business.InviteUser{
		AppId:        businesss.AppId,
		PaymentType:  inviteUserInfo.PaymentType,
		ResourceId:   transformGormString(inviteUserInfo.ResourceId),
		ResourceType: inviteUserInfo.ResourceType,
		ShareUserId:  inviteUserInfo.ShareUserId,
		ProductId:    inviteUserInfo.ProductId,
	}
}

// dto模型向底层po模型转换
func (businesss *InviteBusiness) transformInviteRelationPo(inviteUserInfo InviteRelation) business.InviteRelation {
	return business.InviteRelation{
		AppId:         businesss.AppId,
		PaymentType:   inviteUserInfo.PaymentType,
		ResourceId:    transformGormString(inviteUserInfo.ResourceId),
		ResourceType:  inviteUserInfo.ResourceType,
		ProductId:     inviteUserInfo.ProductId,
		ShareUserId:   inviteUserInfo.ShareUserId,
		ShareType:     inviteUserInfo.ShareType,
		InvitedUserId: businesss.UserId,
	}
}

// 好坑，先这样搞一波
func transformGormString(str string) (nullString sql.NullString) {
	nullString.String = str
	nullString.Valid = true
	return nullString
}
