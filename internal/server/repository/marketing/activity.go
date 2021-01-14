package marketing

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"

	"abs/models/business"
	"abs/pkg/cache/redis_alive"
	e "abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
)

const (
	//资源活动信息的redis缓存key
	resourceActivityInfoRedisKey = "%s:%s:%d"
	//活动信息的redis缓存key
	activityInfoRedisKey = "%s:%s:%s"
	//活动类型-限时折扣
	activityTypeLimitDiscount = "3"
	//活动类型-秒杀
	activityTypeSeckill = "8"
)

type resourceActivityInfo struct {
	Stock       int
	NowStock    int    `json:"now_stock"`
	PriceType   string `json:"price_type"`
	PriceParams string `json:"price_params"`
	RatePrice   int    `json:"rate_price"`
}

type baseActivityInfo struct {
	StartAt       string
	EndAt         string
	ActivityLabel string
}

//获取专栏的活动标签
func GetActivityTags(resources []*business.PayProducts, priceType uint8, client string, appVersion string) []*business.PayProducts {
	if priceType != 1 && priceType != 2 {
		return resources
	}

	conn, err := redis_alive.GetLiveMarketingConn()
	if err != nil {
		logging.Error(err)
		return resources
	}
	defer conn.Close()

	for _, resource := range resources {
		resource.InActivity = 0
		resource.Tags = []string{}

		var (
			//该资源参与所有类型的活动
			allActivities map[string]interface{}
			//该资源参与的某一种类型活动
			resourceActivity map[string]resourceActivityInfo
			//该资源已经查询过的活动信息
			checkedActivities map[string]baseActivityInfo
			//该资源活动类型
			activityType string
			//是否需要检查库存
			isCheckStock = true
		)

		//获取该资源参与的所有活动类型
		info, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf(resourceActivityInfoRedisKey, resource.AppId, resource.Id, resource.SrcType)))
		if info == nil || len(info) < 2 {
			if err != nil {
				logging.Error(err)
			}
			continue
		} else {
			allActivities = make(map[string]interface{})
			for i := 0; i < len(info)-1; i++ {
				if i%2 == 0 {
					allActivities[string(info[i].([]byte))] = info[i+1]
				}
			}
		}

		//有多个类型活动取其中的一个类型
		if v, ok := allActivities[activityTypeLimitDiscount]; ok && v != nil {
			resourceActivity = make(map[string]resourceActivityInfo)
			err = json.Unmarshal(v.([]byte), &resourceActivity)
			activityType = activityTypeLimitDiscount
		} else if v, ok := allActivities[activityTypeSeckill]; ok && v != nil {
			resourceActivity = make(map[string]resourceActivityInfo)
			err = json.Unmarshal(v.([]byte), &resourceActivity)
			activityType = activityTypeSeckill
			isCheckStock = false
		}
		if err != nil {
			logging.Error(err)
			continue
		}

		//遍历该类型下的所有活动
		checkedActivities = make(map[string]baseActivityInfo)
		for activityId, activityInfo := range resourceActivity {
			if _, ok := checkedActivities[activityId]; !ok {
				//查询活动信息
				redisKey := fmt.Sprintf(activityInfoRedisKey, resource.AppId, activityId, activityType)
				startAt, _ := redis.String(conn.Do("HGET", redisKey, "start_at"))
				endAt, _ := redis.String(conn.Do("HGET", redisKey, "end_at"))
				activityLabel, _ := redis.String(conn.Do("HGET", redisKey, "activity_label"))
				checkedActivities[activityId] = baseActivityInfo{
					StartAt:       startAt,
					EndAt:         endAt,
					ActivityLabel: activityLabel,
				}
			}

			if checkedActivities[activityId].StartAt != "" && checkedActivities[activityId].EndAt != "" {
				startAt, err := time.Parse(util.TIME_LAYOUT, checkedActivities[activityId].StartAt)
				endAt, err := time.Parse(util.TIME_LAYOUT, checkedActivities[activityId].EndAt)
				if err != nil {
					logging.Error(err.Error() + " time.Parse Error")
				} else if startAt.Before(time.Now()) && endAt.After(time.Now()) {
					resource.InActivity = 1
					if isCheckStock && !(activityInfo.Stock == 0 || (activityInfo.Stock > 0 && activityInfo.NowStock > 0)) {
						resource.InActivity = 0
					}
					if resource.InActivity == 1 {
						if checkedActivities[activityId].ActivityLabel != "" {
							resource.Tags = []string{checkedActivities[activityId].ActivityLabel}
						} else if activityType == activityTypeLimitDiscount {
							resource.Tags = []string{"限时折扣"}
						} else if activityType == activityTypeSeckill {
							resource.Tags = []string{"秒杀"}
						}

						rate := 100
						if client == strconv.Itoa(e.AGENT_TYPE_APP) && appVersion == "55" {
							// 在55版本，首页显示折扣价(此版本首页价格显示的是后端传的price 与 piecePrice)
							// 在55版本，全部频道的价格，前段是当作分处理的,所以不能 /100，但是首页是当元处理的
							if priceType != 1 {
								rate = 1
							}
							resource.Price = activityInfo.RatePrice / rate
						} else {
							// 其他小程序和h5版本新增rate_price  低于55的版本显示price逻辑不变与详情页一个价格，高于56的版本有rate_price显示rate_price
							resource.RatePrice = activityInfo.RatePrice / 100
						}
					}
				}
			}
		}
	}
	
	return resources
}
