package validator

import "time"

// v2/alive_list.go/GetSubscribeAliveListByDate
type GetSubscribeAliveListByDateV2 struct {
	AppId            string    `form:"app_id" json:"app_id" binding:"required"`
	UniversalUnionId string    `form:"universal_union_id" json:"universal_union_id" binding:"required"`
	StartTime        time.Time `form:"start_time" json:"start_time" binding:"required" time_format:"2006-01-02 15:04:05"`
	EndTime          time.Time `form:"end_time" json:"end_time" binding:"required" time_format:"2006-01-02 15:04:05"`
}

// v2/alive_list.go/GetSubscribeAliveListByDate
type GetSubscribeAliveNumByDateV2 struct {
	AppId            string    `form:"app_id" json:"app_id" binding:"required"`
	UniversalUnionId string    `form:"universal_union_id" json:"universal_union_id" binding:"required"`
	StartTime        time.Time `form:"start_time" json:"start_time" binding:"required" time_format:"2006-01-02 15:04:05"`
	EndTime          time.Time `form:"end_time" json:"end_time" binding:"required" time_format:"2006-01-02 15:04:05"`
}

// v2/alive_list.go/GetSubscribeLivingAliveList
type GetSubscribeLivingAliveListV2 struct {
	AppId            string    `form:"app_id" json:"app_id" binding:"required"`
	AppIds           string    `form:"app_ids" json:"app_ids" binding:"required"`
	UniversalUnionId string    `form:"universal_union_id" json:"universal_union_id" binding:"required"`
	StartTime        time.Time `form:"start_time" json:"start_time" binding:"required" time_format:"2006-01-02 15:04:05"`
	EndTime          time.Time `form:"end_time" json:"end_time" binding:"required" time_format:"2006-01-02 15:04:05"`
}

// v2/alive_list.go/GetSubscribeLivingAliveList
type GetSubscribeUnStartAliveListV2 struct {
	AppId            string `form:"app_id" json:"app_id" binding:"required"`
	AppIds           string `form:"app_ids" json:"app_ids" binding:"required"`
	UniversalUnionId string `form:"universal_union_id" json:"universal_union_id" binding:"required"`
}
