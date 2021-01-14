package course

import (
	"fmt"
	"strings"
	"sync"

	// "github.com/goinggo/mapstructure"
	"github.com/gomodule/redigo/redis"

	"abs/models/business"
	"abs/models/sub_business"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"
)

type Product struct {
	AppId      string
	ResourceId string
}

const (
	// Redis key
	proResRelation = "pro_res_relation_%s_%s"
	activityTags   = "%s:%s:%s"

	// 设置缓存过期时间
	// 资源（直播）关联关系
	resourceRelationCacheTime = "60"

	// 类型类的全局变量
	LIMIT_ACCOUNT_ACTIVITY_TYPE = "3"
	SECKILL_ACTIVITY_TYPE       = "8"
)

// 获取资源的专栏关联信息
func (p *Product) GetResourceRelation() ([]*business.ProResRelation, error) {
	var cacheProResRelation []*business.ProResRelation
	conn, _ := redis_alive.GetLiveInteractConn()
	defer conn.Close()

	cacheKey := fmt.Sprintf(proResRelation, p.AppId, p.ResourceId)
	info, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err == nil {
		util.JsonDecode(info, &cacheProResRelation)
		return cacheProResRelation, nil
	}

	cacheProResRelation, err = business.GetResourceProducts(p.AppId, p.ResourceId)
	if err != nil {
		return cacheProResRelation, err
	}

	if value, err := util.JsonEncode(cacheProResRelation); err == nil {
		if _, err = conn.Do("SET", cacheKey, value, "EX", resourceRelationCacheTime); err != nil {
			logging.Error(err)
		}
	}

	return cacheProResRelation, nil
}

// 获取所属的栏目
func (p *Product) GetParentColumns(relations []*business.ProResRelation) (products []*business.PayProducts, err error) {
	// 这一步需要吗？看逻辑不需要啊，反正都要查专栏信息，先跳过趴...
	if false {
		// 也要实例化一个锁
		wg, lock := sync.WaitGroup{}, sync.Mutex{}
		for k, v := range relations {
			wg.Add(1)
			go func(k int, v *business.ProResRelation) { // 注意：循环体协程需要这样带值！！
				defer wg.Done()
				product, err := business.GetPayProductState(v.AppId, v.ResourceId)
				if err != nil || product.Id == "" {
					// 多协程操作同一个map或者切片要加锁！！！否则会报不是协程安全
					lock.Lock()
					// 切片删除性能低下，如果是大数据请使用链表实现！！！
					relations = append(relations[:k], relations[k+1:]...)
					lock.Unlock()
				}
			}(k, v)
		}
		wg.Wait()
	}

	if len(relations) > 0 {
		rids := []string{}
		for _, v := range relations {
			rids = append(rids, v.ProductId)
		}

		products, err = business.GetPayProductByIds(p.AppId, strings.Join(rids, ","))
		if err != nil {
			logging.Error(err)
			return
		}
		for _, val := range products {
			val.RatePrice = -1
			if val.IsMember == 1 {
				val.SrcType = business.SINGLE_GOODS_MEMBER
			} else {
				val.SrcType = business.SINGLE_GOODS_PACKAGE
			}
		}
	}
	return
}

// 获取营期列表
func (p *Product) GetCampTermListByIds(relations []*business.ProResRelation) ([]*business.PayProducts, error) {
	terms, ids := []*business.PayProducts{}, []string{}
	for _, v := range relations {
		if strings.HasPrefix(v.ProductId, "term") {
			ids = append(ids, v.ProductId)
		}
	}
	if len(ids) == 0 {
		return terms, nil
	}
	//查询字段
	selectFields := []string{"app_id", "id", "img_url", "img_url_compressed", "title", "summary", "join_count",
		"price", "display_state", "distribute_percent", "first_distribute_percent", "lesson_start_at", "lesson_stop_at", "recycle_bin_state"}
	// 初始化营期请求服务
	campReq := service.CampService{AppId: p.AppId}
	terms, err := campReq.GetCampTermInfo(ids, selectFields)
	if err != nil {
		logging.Error(err)
	}
	return terms, nil

	// 老方法，注意废弃！！！
	// if termResult["code"].(float64) == 0 {
	// 	rdata := termResult["data"].(map[string]interface{})
	// 	var item business.PayProducts
	// 	for _, val := range rdata["terms"].([]interface{}) {
	// 		v := val.(map[string]interface{})
	// 		v["app_id"] = p.AppId
	// 		err = mapstructure.Decode(v, &item)
	// 		// 这里注意下！！！有过滤
	// 		if err != nil {
	// 			logging.Error(err)
	// 			continue
	// 		}
	// 		terms = append(terms, &item)
	// 	}
	// }
}

// 小程序部分代码此期不用
// 过滤父级的专栏关联信息
func (p *Product) FilterParentColumns(resources []*business.PayProducts, client, userAgent string) []*business.PayProducts {
	if len(resources) > 0 {
		wg, lock := sync.WaitGroup{}, sync.Mutex{}
		for k, v := range resources {
			if v.RecycleBinState == 1 {
				resources = append(resources[:k], resources[k+1:]...)
			} else {
				if util.GetMiniProgramVersion(client, userAgent) == 2 {
					v.Price = 0 // 不显示价格
					pResType := 6
					if v.MemberType == 1 && v.IsMember == 1 {
						// 会员 在ios小程序不显示在列表
						resources = append(resources[:k], resources[k+1:]...)
						continue
					} else if v.MemberType == 2 && v.IsMember == 1 {
						// 大专栏
						pResType = 8
					}
					wg.Add(1)
					go func(v *business.PayProducts, k int, pType int) {
						specInfo, err := sub_business.GetSpecInfo(v.AppId, v.Id, pType)
						if err == nil && specInfo.Id != 0 {
							// 勾选了不显示的资源
							if specInfo.State == 1 {
								lock.Lock()
								resources = append(resources[:k], resources[k+1:]...)
								lock.Unlock()
							} else {
								v.PurchaseCount = 0
								if specInfo.ImgUrl != "" {
									v.ImgUrl.String = specInfo.ImgUrl
								}
								if specInfo.ImgUrlCompressed != "" {
									v.ImgUrlCompressed.String = specInfo.ImgUrlCompressed
								}
								if specInfo.Title != "" {
									v.Name.String = specInfo.Title
								}
								if specInfo.OrgSummary != "" {
									v.Summary.String = specInfo.OrgSummary
								}
							}
						}
					}(v, k, pResType)
				}
			}
		}
		wg.Wait()
	}
	return resources
}
