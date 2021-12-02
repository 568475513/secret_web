// Package course 直播关联父级列表信息
package course

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/gomodule/redigo/redis"

	"abs/models/alive"
	"abs/models/business"
	"abs/models/sub_business"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"
)

type ProductInfo struct {
	AppId     string
	AliveId   string
	ProductId string
	Channel   string
	BuzUri    string
	TargetUrl string
}

const (
	columnsIdsCacheKeyPre = "all_column:%s:%s"
	columnsIdsCacheTime   = "60"

	columnPhaseCacheKeyPre = "column_phase_update_num:%s:%s"
	columnPhaseCacheTime   = "60"

	columnsDetailsInfoCacheKeyPre = "columns_details_info:%s:%s"
	columnsDetailsInfoCacheTime   = "60"
)

// GetFromTargetUrl 从target_url提取元素
func (pi *ProductInfo) GetFromTargetUrl(key string) string {
	var (
		err          error
		value        string
		targetUrlObj = make(map[string]string)
	)

	err = json.Unmarshal([]byte(key), &targetUrlObj)
	if err != nil {
		logging.Error(fmt.Sprintf("GetFromTargetUrl json.Unmarshal fails: %s", err.Error()))
		return value
	}

	value, ok := targetUrlObj[key]
	if !ok {
		logging.Error(fmt.Sprintf("GetFromTargetUrl get %s from target_url fails", key))
	}
	return value
}

// GetAliveProductsInfo 获取直播关联父级信息
// 从业务后台get_resource_info铲过来的屎，看了代码记得洗眼睛，改了代码记得洗手洗键盘
func (pi *ProductInfo) GetAliveProductsInfo(paymentType int) (result []map[string]interface{}) {
	//申明局部变量
	var (
		err           error
		contentAppId  string
		cacheKey      string
		pRelationIds  []string                    //该直播关联的所有上级的id
		termIds       []string                    //该直播关联的所有训练营id
		pRelationList []*business.ProductRelation //该直播所有上级的关联关系
		pDetailsInfos []*business.PayProducts     //该直播关联的所有上级详情
	)

	//获取redis连接
	conn, err := redis_alive.GetSubBusinessConn()
	if err != err {
		logging.Error(fmt.Sprintf("GetAliveProductsInfo redis conn fails: %s", err.Error()))
	}
	defer conn.Close()

	//先从缓存查关联父级数据
	contentAppId = pi.GetFromTargetUrl("content_app_id")
	if contentAppId == "" {
		cacheKey = fmt.Sprintf(columnsIdsCacheKeyPre, pi.AppId, pi.AliveId)
	} else {
		cacheKey = fmt.Sprintf(columnsIdsCacheKeyPre, contentAppId, pi.AliveId)
	}
	cacheData, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil && err != redis.ErrNil {
		logging.Error(fmt.Sprintf("GetAliveProductsInfo redis get %s fails: %s", cacheKey, err.Error()))
	} else {
		err = json.Unmarshal(cacheData, &pRelationList)
		if err != nil {
			logging.Error(fmt.Sprintf("GetAliveProductsInfo json.Unmarshal fails: %s", err.Error()))
		}
	}

	//无数据则查库 todo::注意缓存穿透问题
	if pRelationIds == nil {
		if contentAppId == "" {
			pRelationList, err = business.GetResRelation(pi.AppId, pi.AliveId, []string{"*"})
		} else {
			pRelationList, err = business.GetResRelation(contentAppId, pi.AliveId, []string{"*"})
		}
		//缓存结果
		if err == nil {
			cacheData, err = json.Marshal(pRelationList)
			if err != nil {
				logging.Error(fmt.Sprintf("GetAliveProductsInfo json.Marshal fails: %s", err.Error()))
			}
			_, err = conn.Do("SETEX", cacheKey, columnsIdsCacheTime, string(cacheData))
			if err != nil {
				logging.Error(fmt.Sprintf("GetAliveProductsInfo reids SETEX %s fails: %s", cacheKey, err.Error()))
			}
		}
	}

	//提取父级id
	pRelationIds = pi.getProductIds(pRelationList)
	if len(pRelationIds) == 0 {
		return
	}

	//查询父级信息
	if contentAppId == "" {
		pDetailsInfos = pi.getProductsDetailsInfo(pRelationIds, pi.AppId)
		//查询训练营的信息
		for _, item := range pRelationList {
			if item.ProductType == enums.ResourceTypeLive {
				termIds = append(termIds, item.ProductId)
			}
			if len(termIds) != 0 {
				cs := service.CampService{AppId: pi.AppId}
				termInfos, err := cs.GetCampTermInfo(termIds, []string{
					"app_id",
					"id",
					"title",
					"price",
					"distribute_percent",
					"first_distribute_percent",
					"join_count",
					"img_url",
					"display_state",
					"recycle_bin_state",
					"img_url_compressed",
					"created_at"})
				if err != nil {
					logging.Error(fmt.Sprintf("GetAliveProductsInfo GetCampTermInfo fails: %s", err.Error()))
				} else {
					//todo::注意训练营和专栏部分字段名称不一致问题
					pDetailsInfos = append(pDetailsInfos, termInfos...)
				}
			}
		}
	} else {
		pDetailsInfos = pi.getProductsDetailsInfo(pRelationIds, contentAppId)
	}

	//遍历父级列表数据，一些逻辑处理 todo::确认此处的item变更是否生效
	var termSlice []map[string]interface{}
	if len(pDetailsInfos) != 0 {
		pfRelationList := pi.formatRelationList(pRelationList)
		for _, item := range pDetailsInfos {
			pfRelation, ok := pfRelationList[item.Id]
			if ok {
				item.CreatedAt = pfRelation.CreatedAt
				item.ResourceId = pfRelation.ResourceId
				item.IsTry = pfRelation.IsTry
			} else {
				item.ResourceId = pi.AliveId
				item.IsTry = 0
			}

			var pType int
			if item.IsMember == 1 && item.MemberType == 1 {
				pType = enums.ResourceTypeActivity
			} else if item.IsMember == 0 && item.MemberType == 1 {
				pType = enums.ResourceTypeTopic
			} else if item.IsMember == 4 {
				pType = enums.ResourceTypeCamp
			} else {
				pType = enums.ResourceTypeCamp
			}

			//todo::需要注意默认值问题
			if item.RecycleBinState != 1 && item.State != 1 {
				columnInfo := make(map[string]interface{})
				columnInfo["app_id"] = item.AppId
				columnInfo["id"] = item.Id
				columnInfo["title"] = item.Name
				columnInfo["img_url"] = item.ImgUrl
				if pType == enums.ResourceTypeCamp {
					//营期不显示这些字段
					columnInfo["update_num"] = 0
					columnInfo["resource_count"] = 0
				} else {
					columnInfo["update_num"] = pi.GetUpdatePhase(item.Id)
					columnInfo["resource_count"] = item.ResourceCount
				}
				columnInfo["is_member"] = item.IsMember
				columnInfo["member_type"] = item.MemberType
				columnInfo["purchase_count"] = item.PurchaseCount
				columnInfo["sell_type"] = item.SellType
				result = append(result, columnInfo)
			} else if len(termSlice) == 0 {
				//确保就算上级全部下架或隐藏，也留一个临时上级，避免无处跳转
				termInfo := make(map[string]interface{})
				termInfo["app_id"] = item.AppId
				termInfo["id"] = item.Id
				termInfo["title"] = item.Name
				termInfo["img_url"] = item.ImgUrl
				if pType == enums.ResourceTypeCamp {
					//营期不显示这些字段
					termInfo["update_num"] = 0
					termInfo["resource_count"] = 0
				} else {
					termInfo["update_num"] = pi.GetUpdatePhase(item.Id)
					termInfo["resource_count"] = item.ResourceCount
				}
				termInfo["is_member"] = item.IsMember
				termInfo["member_type"] = item.MemberType
				termInfo["purchase_count"] = item.PurchaseCount
				termInfo["sell_type"] = item.SellType
				termSlice = append(termSlice, termInfo)
			}
		}
	}
	if len(result) == 0 && paymentType == 3 {
		//原注释："非单卖资源现在无上级？？要爆炸"
		//我也不知道这段逻辑是干嘛，我也不敢动它
		result = append(result, termSlice...)
	}

	return
}

// GetMoreInfo 拼接展示多个父级课程的中间页url
func (pi *ProductInfo) GetMoreInfo(productList []map[string]interface{}, alive *alive.Alive) (url string) {
	if len(productList) < 2 {
		return
	}
	path := util.ContentUrl(util.ContentParam{
		ResourceType: enums.ResourceTypeLive,
		ResourceId:   alive.Id,
		ChannelId:    pi.Channel,
	})
	url = util.UrlWrapper(path, pi.BuzUri, pi.AppId)
	return url
}

//格式化父级关联关系数据
func (pi *ProductInfo) formatRelationList(pr []*business.ProductRelation) map[string]*business.ProductRelation {
	var result = make(map[string]*business.ProductRelation)
	if len(pr) == 0 {
		return result
	}
	for _, item := range pr {
		if item.ProductId != "" {
			result[item.ProductId] = item
		}
	}
	return result
}

//提取product_id
func (pi *ProductInfo) getProductIds(data []*business.ProductRelation) (result []string) {
	if len(data) == 0 {
		return
	}
	for _, item := range data {
		result = append(result, item.ProductId)
	}
	return
}

//批量查询父级详情信息
func (pi *ProductInfo) getProductsDetailsInfo(productIds []string, contentAppId string) []*business.PayProducts {
	var (
		cacheKey string
		result   []*business.PayProducts
		s        = []string{
			"app_id",
			"id",
			"name",
			"price",
			"resource_count",
			"distribute_percent",
			"first_distribute_percent",
			"is_member",
			"member_type",
			"purchase_count",
			"img_url",
			"state",
			"recycle_bin_state",
			"img_url_compressed",
			"created_at",
			"sell_type"}
	)

	//获取redis连接
	conn, err := redis_alive.GetSubBusinessConn()
	if err != nil {
		logging.Error(fmt.Sprintf("getProductsDetailsInfo conn fails: %s", err.Error()))
	}
	defer conn.Close()

	//先从redis获取数据
	if pi.AppId != contentAppId {
		cacheKey = fmt.Sprintf(columnsDetailsInfoCacheKeyPre, pi.AppId, pi.AliveId)
	} else {
		cacheKey = fmt.Sprintf(columnsDetailsInfoCacheKeyPre, contentAppId, pi.AliveId)
	}
	cacheData, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err == nil {
		err = json.Unmarshal(cacheData, &result)
		if err != nil {
			logging.Error(fmt.Sprintf("getProductsDetailsInfo json.Unmarshal fails: %s", err.Error()))
		}
		return result
	} else if err != redis.ErrNil {
		logging.Error(fmt.Sprintf("getProductsDetailsInfo redis get %s fails: %s", cacheKey, err.Error()))
	}

	//无缓存则查数据库
	if pi.AppId != contentAppId {
		var resourceIds []string
		//分销商品
		channelData, err := sub_business.GetChannelRepertoryList(contentAppId, pi.AppId, productIds, []string{"resource_id"})
		if err != nil {
			logging.Error(fmt.Sprintf("getProductsDetailsInfo GetChannelRepertoryList fails: %s", err.Error()))
			return result
		}
		//提取resource_id
		for _, item := range channelData {
			resourceIds = append(resourceIds, item.ResourceId)
		}
		//查询具体信息
		result, err = business.GetInfoBatch(contentAppId, resourceIds, s)
		if err != nil {
			logging.Error(fmt.Sprintf("getProductsDetailsInfo GetInfoBatch fails: %s", err.Error()))
			return result
		}
	} else {
		//查询具体信息
		result, err = business.GetInfoBatch(contentAppId, productIds, s)
		if err != nil {
			logging.Error(fmt.Sprintf("getProductsDetailsInfo GetInfoBatch fails: %s", err.Error()))
			return result
		}
	}

	//写入缓存
	cacheData, err = json.Marshal(result)
	if err != nil {
		logging.Error(fmt.Sprintf("getProductsDetailsInfo json.Marshal fails: %s", err.Error()))
	}
	_, err = conn.Do("SETEX", cacheKey, columnsDetailsInfoCacheTime, cacheData)
	if err != nil {
		logging.Error(fmt.Sprintf("getProductsDetailsInfo redis SETEX %s fails: %s", cacheKey, err.Error()))
	}

	return result
}

// GetUpdatePhase 查询父级资源的更新期数
func (pi *ProductInfo) GetUpdatePhase(productId string) int {
	var (
		err                  error
		total                int
		cacheKey             string
		relationResourceList []*business.ProductRelation
	)

	//获取redis连接
	conn, err := redis_alive.GetSubBusinessConn()
	if err != err {
		logging.Error(fmt.Sprintf("GetAliveProductsInfo redis conn fails: %s", err.Error()))
	}
	defer conn.Close()

	//先从缓存里面查
	cacheKey = fmt.Sprintf(columnPhaseCacheKeyPre, pi.AppId, productId)
	total, err = redis.Int(conn.Do("GET", cacheKey))
	if err == nil {
		return total
	} else if err != redis.ErrNil {
		logging.Error(fmt.Sprintf("GetUpdatePhase redis get %s fails: %s", cacheKey, err.Error()))
	}

	//无缓存则查询该父级下所有关联关系
	relationResourceList, err = business.GetResByProductId(pi.AppId, productId, []string{"resource_id", "resource_type"})
	if err != nil {
		logging.Error(fmt.Sprintf("GetUpdatePhase GetResByProductId fails: %s", err.Error()))
		return total
	}

	if len(relationResourceList) > 0 {
		var (
			aliveNum   int
			audioNum   int
			videoNum   int
			eBookNum   int
			imgTextNum int
			aliveIds   []string
			audioIds   []string
			videoIds   []string
			eBookIds   []string
			imgTextIds []string
			nowTime    = time.Now().Format(util.TIME_LAYOUT)
		)
		//根据resource_type分拣
		for _, item := range relationResourceList {
			if item.ResourceType == enums.ResourceTypeLive {
				aliveIds = append(aliveIds, item.ResourceId)
			} else if item.ResourceType == enums.ResourceTypeAudio {
				audioIds = append(audioIds, item.ResourceId)
			} else if item.ResourceType == enums.ResourceTypeVideo {
				videoIds = append(videoIds, item.ResourceId)
			} else if item.ResourceType == enums.ResourceTypeEBook {
				eBookIds = append(eBookIds, item.ResourceId)
			} else if item.ResourceType == enums.ResourceTypeImageText {
				imgTextIds = append(imgTextIds, item.ResourceId)
			}
		}
		//分别查询各类资源数量
		if len(audioIds) > 0 {
			audioNum, err = business.CountAudio(pi.AppId, audioIds, nowTime)
			if err != nil {
				logging.Error(fmt.Sprintf("GetUpdatePhase CountAudio fails: %s", err.Error()))
			}
		}
		if len(videoIds) > 0 {
			videoNum, err = business.CountVideo(pi.AppId, videoIds, nowTime)
			if err != nil {
				logging.Error(fmt.Sprintf("GetUpdatePhase CountVideo fails: %s", err.Error()))
			}
		}
		if len(eBookIds) > 0 {
			eBookNum, err = business.CountEBook(pi.AppId, eBookIds, nowTime)
			if err != nil {
				logging.Error(fmt.Sprintf("GetUpdatePhase CountEBook fails: %s", err.Error()))
			}
		}
		if len(imgTextIds) > 0 {
			imgTextNum, err = business.CountImgText(pi.AppId, imgTextIds, nowTime)
			if err != nil {
				logging.Error(fmt.Sprintf("GetUpdatePhase CountImgText fails: %s", err.Error()))
			}
		}
		if len(aliveIds) > 0 {
			//分批查询，每次最多查200
			pageSize := 200
			loop := int(math.Ceil(float64(len(aliveIds)) / float64(pageSize)))
			for i := 0; i < loop; i++ {
				start := i * pageSize
				end := start + pageSize
				if end > len(aliveIds) {
					end = len(aliveIds)
				}
				batch := aliveIds[start:end]
				num, err := alive.CountAlive(pi.AliveId, batch, nowTime)
				if err != nil {
					logging.Error(fmt.Sprintf("GetUpdatePhase CountAlive fails: %s", err.Error()))
				}
				aliveNum += num
			}
		}
		total = aliveNum + audioNum + videoNum + imgTextNum + eBookNum
	}

	//缓存结果
	_, err = conn.Do("SETEX", cacheKey, columnPhaseCacheTime, total)
	if err != nil {
		logging.Error(fmt.Sprintf("GetUpdatePhase redis setex %s fails: %s", cacheKey, err.Error()))
	}

	return total
}

// DealProductsInfo 补充替换父级列表信息部分字段
func (pi *ProductInfo) DealProductsInfo(productList []map[string]interface{}, baseConf *service.AppBaseConf, client int, moduleProfit map[string]interface{}) []map[string]interface{} {
	if len(productList) == 0 {
		return productList
	}
	contentAppId := pi.GetFromTargetUrl("content_app_id")
	//todo::确认product有没有修改得到
	for _, product := range productList {
		//是否显示订阅数
		if baseConf.HideSubCount == 1 || client == 2 {
			profit, ok := moduleProfit["hide_sub_count_is_remind"].(int)
			if ok && (profit == 1 || profit == 0) {
				product["purchase_count"] = 0
			}
		}
		//是否显示期数
		if baseConf.IsShowResourcecount == 0 {
			profit, ok := moduleProfit["hide_resource_count"].(int)
			if ok && (profit == 1 || profit == 0) {
				product["update_num"] = 0
			}
		}
		var resourceType int
		if product["is_member"] == 0 {
			resourceType = enums.ResourceTypePackage
		} else if product["member_type"] == 1 {
			resourceType = enums.ResourceTypeActivity
		} else {
			resourceType = enums.ResourceTypeTopic
		}
		if product["id"].(string)[0:5] == "term" {
			resourceType = enums.ResourceTypeCamp
		}
		path := util.ContentUrl(util.ContentParam{
			Type:         enums.PaymentTypeProduct,
			ResourceType: resourceType,
			ResourceId:   "",
			ProductId:    product["id"].(string),
		})
		if contentAppId == "" {
			util.UrlWrapper(path, pi.BuzUri, pi.AppId)
		} else {
			util.UrlWrapper(path, pi.BuzUri, contentAppId)
		}
	}
	return productList
}
