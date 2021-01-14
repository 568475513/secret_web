package validator

// 指明控制器所属！
// v2/courseware.go/GetCourseWareRecords
type CourseWareRecordsRuleV2 struct {
	AliveId   string `json:"alive_id" form:"alive_id" binding:"required,startswith=l_"`
	Client    int    `json:"client" form:"client"`
	AliveTime int    `json:"alive_time" form:"alive_time"`
	PageSize  int    `json:"page_size" form:"page_size"`
}

// v2/courseware.go/GetCourseWareInfo
type CourseWareInfoRuleV2 struct {
	AliveId      string `json:"alive_id" form:"alive_id" binding:"required,startswith=l_"`
	CourseWareId string `json:"courseware_id" form:"courseware_id"`
}
