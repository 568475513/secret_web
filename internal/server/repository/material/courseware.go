package material

import (
	"abs/models/alive"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strings"
)

/**
 * 课件基本结构体
 */
type CourseWare struct {
	AppId   string
	AliveId string
}

const (
	aliveCourseWareInfoKey          = "aliveCoursewareInfo:%s:%s:%s:%s"
	aliveCourseWareInfoByAliveIdKey = "aliveCoursewareInfoByAliveId:%s:%s:%s"
	aliveRecordInfoKey              = "aliveRecordInfo:%s:%s"

	// 缓存时间控制(秒)
	// 获取课件使用记录
	courseWareRecordsCacheTime = "30"
	// -
	courseWareInfoCacheTime = "30"
)

/**
 * 获取课件详情（通过courseWareId）
 */
func (c *CourseWare) GetCourseWareInfo(courseWareId string) (*alive.CourseWare, error) {

	var (
		err                 error
		cacheCourseWareInfo []*alive.CourseWare
	)

	//从缓存中获取
	courseWareIdArr := []string{courseWareId}
	cacheCourseWareInfo, err = c.GetCourseWareInfoCache(courseWareIdArr, []string{
		"id", "page_count", "courseware_image", "current_preview_page"})
	if err != nil {
		return nil, err
	}

	//判断是否有数据
	if len(cacheCourseWareInfo) == 0 {
		return nil, err
	}

	//定义课件数组字段
	var coursewareImage []map[string]interface{}
	err = json.Unmarshal([]byte(cacheCourseWareInfo[0].CoursewareImage), &coursewareImage)
	if err != nil {
		return nil, err
	}

	//赋值到新字段
	cacheCourseWareInfo[0].CourseImageArray = coursewareImage

	return cacheCourseWareInfo[0], nil
}

/**
 * 获取课件详情（通过AliveId）
 */
func (c *CourseWare) GetCourseWareInfoByAliveId() (*alive.CourseWare, error) {

	var (
		err                 error
		cacheCourseWareInfo *alive.CourseWare
	)

	//从缓存中获取
	cacheCourseWareInfo, err = c.GetCourseWareInfoCacheByAliveId([]string{
		"id", "page_count", "courseware_image", "current_preview_page"})
	if err != nil {
		return nil, err
	}

	//判断是否有数据
	if cacheCourseWareInfo.Id.String == "" {
		return nil, err
	}

	//定义课件数组字段
	var coursewareImage []map[string]interface{}
	err = json.Unmarshal([]byte(cacheCourseWareInfo.CoursewareImage), &coursewareImage)
	if err != nil {
		return nil, err
	}

	//赋值到新字段
	cacheCourseWareInfo.CourseImageArray = coursewareImage

	return cacheCourseWareInfo, nil
}

/**
 * 获取课件使用记录
 */
func (c *CourseWare) GetCourseWareRecords(client int, aliveTime int, pageSize int) ([]*alive.CourseWareRecords, error) {

	var (
		err                    error
		cacheCourseWareRecords []*alive.CourseWareRecords
		lookBackFile           *alive.AliveLookBack
	)

	//获取剪辑表的id
	lookBackId := 0
	lookBackFile, err = alive.GetAliveLookBackFile(c.AppId, c.AliveId, []string{"id"})
	if err != nil {
		return nil, err
	}
	if lookBackFile != nil && lookBackFile.Id != 0 {
		lookBackId = lookBackFile.Id //剪辑id
	}

	if client != 1 { // 0-小程序 1-公众号
		cacheCourseWareRecords, err = alive.GetCourseWareByAliveTime(c.AppId, c.AliveId, lookBackId, 0, pageSize, true)
		if err != nil {
			return nil, err
		}
		if len(cacheCourseWareRecords) > 0 && cacheCourseWareRecords[0].AppId != "" {
			for k, v := range cacheCourseWareRecords {
				cacheCourseWareRecords[k].AliveTime = v.CourseUseTime
			}
		}
	} else { //公众号
		cacheKey := fmt.Sprintf(aliveRecordInfoKey, c.AppId, c.AliveId)
		conn, err := redis_alive.GetLiveBusinessConn()
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		if aliveTime == 0 { //如果是获取从alive_time = 0开始的课件，则从Redis里获取
			recordInfo, err := redis.Bytes(conn.Do("GET", cacheKey))
			if err != nil {
				logging.Error(err)
			} else {
				json.Unmarshal(recordInfo, &cacheCourseWareRecords)
			}
		}

		//缓存没数据从数据库中拿
		if len(cacheCourseWareRecords) == 0 || cacheCourseWareRecords[0].AppId == "" {
			cacheCourseWareRecords, err = alive.GetCourseWareByAliveTime(c.AppId, c.AliveId, lookBackId, aliveTime, pageSize, true)
			if err != nil {
				return nil, err
			}
			if aliveTime != 0 {
				//将alive_time插入到最前面
				preview, err := alive.GetCourseWareByAliveTime(c.AppId, c.AliveId, lookBackId, aliveTime, 1, false)
				if err != nil {
					return nil, err
				}
				if preview != nil && preview[0].AppId != "" { //将preview插入到头部
					cacheCourseWareRecords = append(preview, cacheCourseWareRecords...)
				}
			}

			//插入到redis
			if len(cacheCourseWareRecords) > 0 && cacheCourseWareRecords[0].AppId != "" {
				for k, v := range cacheCourseWareRecords { //保持两个时间一致
					cacheCourseWareRecords[k].AliveTime = v.CourseUseTime
				}

				if aliveTime == 0 { //保存缓存
					value, err := json.Marshal(cacheCourseWareRecords)
					if err != nil {
						return nil, err
					}

					_, err = conn.Do("SET", cacheKey, value, "EX", courseWareRecordsCacheTime)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return cacheCourseWareRecords, nil
}

/**
 * 走缓存（通过courseWareId）
 */
func (c *CourseWare) GetCourseWareInfoCache(courseWareId []string, s []string) ([]*alive.CourseWare, error) {

	var (
		cacheCourseWareInfo []*alive.CourseWare
	)

	//连接redis
	conn, err := redis_alive.GetLiveBusinessConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//拼接查询字段,作为key
	joinsKey := strings.Join(s, ":")

	cacheKey := fmt.Sprintf(aliveCourseWareInfoKey, c.AppId, c.AliveId, courseWareId, joinsKey)
	courseWareInfo, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil {
		logging.Error(err)
	} else {
		json.Unmarshal(courseWareInfo, &cacheCourseWareInfo)
	}

	if len(cacheCourseWareInfo) == 0 || cacheCourseWareInfo[0].Id.String == "" { //缓存中没有数据
		cacheCourseWareInfo, err = alive.GetCourseWareInfo(c.AppId, c.AliveId, courseWareId, s)
		if err != nil {
			return nil, err
		}

		if len(cacheCourseWareInfo) > 0 && cacheCourseWareInfo[0].Id.String != "" {
			value, err := json.Marshal(cacheCourseWareInfo)
			if err != nil {
				return nil, err
			}

			_, err = conn.Do("SET", cacheKey, value, "EX", "1")
			if err != nil {
				return nil, err
			}
		}
	}

	return cacheCourseWareInfo, nil
}

/**
 * 走缓存（通过aliveId）
 */
func (c *CourseWare) GetCourseWareInfoCacheByAliveId(s []string) (*alive.CourseWare, error) {

	var (
		cacheCourseWareInfo *alive.CourseWare
	)

	//连接redis
	conn, err := redis_alive.GetLiveBusinessConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//拼接查询字段,作为key
	joinsKey := strings.Join(s, ":")

	cacheKey := fmt.Sprintf(aliveCourseWareInfoByAliveIdKey, c.AppId, c.AliveId, joinsKey)
	cacheCourseWareInfo = &alive.CourseWare{}
	courseWareInfo, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil {
		logging.Error(err)
	} else {
		json.Unmarshal(courseWareInfo, &cacheCourseWareInfo)
	}

	if cacheCourseWareInfo.Id.String == "" { //缓存中没有数据
		cacheCourseWareInfo, err = alive.GetCourseWareInfoByAliveId(c.AppId, c.AliveId, s)
		if err != nil {
			return nil, err
		}

		if cacheCourseWareInfo.Id.String != "" {
			value, err := json.Marshal(cacheCourseWareInfo)
			if err != nil {
				return nil, err
			}

			_, err = conn.Do("SET", cacheKey, value, "EX", "1")
			if err != nil {
				return nil, err
			}
		}
	}

	return cacheCourseWareInfo, nil
}
