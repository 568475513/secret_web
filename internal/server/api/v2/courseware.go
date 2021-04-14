package v2

import (
	"encoding/json"
	"fmt"
	// "strconv"

	"github.com/gin-gonic/gin"

	"abs/internal/server/repository/course"
	"abs/internal/server/repository/material"
	"abs/internal/server/rules/validator"
	"abs/pkg/app"
	"abs/pkg/enums"
)

/**
 * 课件详情统一返回结构
 */
type CourseWareInfo struct {
	Id                 string                   `json:"id"`
	PageCount          int                      `json:"page_count"`
	CoursewareImage    []map[string]interface{} `json:"courseware_image"`
	CurrentPreviewPage int                      `json:"current_preview_page"`
}

/**
 * 课件使用记录返回结构
 */
type CourseWareRecords struct {
	AliveId            string `json:"alive_id"`
	AliveTime          int    `json:"alive_time"`
	CourseUseTime      int    `json:"course_use_time"`
	UserId             string `json:"user_id"`
	CurrentPreviewPage int    `json:"current_preview_page"`
	CurrentImageUrl    string `json:"current_image_url"`
	CoursewareId       string `json:"courseware_id"`
}

// 获取课件使用记录接口
func GetCourseWareRecords(c *gin.Context) {
	var (
		err        error
		courseWare CourseWareRecords
		returnData []CourseWareRecords
		req        validator.CourseWareRecordsRuleV2
	)

	// 参数校验
	AppId := app.GetAppId(c)
	if err = app.ParseQueryRequest(c, &req); err != nil {
		return
	}

	if req.PageSize == 0 { //默认100条数据
		req.PageSize = 100
	}

	if AppId == "" || req.AliveId == "" {
		app.FailWithMessage("缺失必要参数", enums.Code_Db_Not_Find, c)
		return
	}

	//获取课件使用记录
	courseWareRep := material.CourseWare{AppId: AppId, AliveId: req.AliveId}
	data, err := courseWareRep.GetCourseWareRecords(c.GetInt("client"), req.AliveTime, req.PageSize)
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("获取课件使用记录错误:%s", err.Error()), enums.ERROR, c)
		return
	}

	// 替换错误链接
	courseWareRep.ReplaceCourseLinkArrStr(data)

	//统一返回结构
	for _, v := range data {
		courseWare = CourseWareRecords {
			AliveId:            v.AliveId,
			AliveTime:          v.AliveTime,
			CourseUseTime:      v.CourseUseTime,
			UserId:             v.UserId.String,
			CurrentPreviewPage: v.CurrentPreviewPage,
			CurrentImageUrl:    v.CurrentImageUrl.String,
			CoursewareId:       v.CoursewareId.String,
		}
		returnData = append(returnData, courseWare)
	}

	//返回
	app.OkWithData(returnData, c)
}

/**
 * 获取课件详情接口
 */
func GetCourseWareInfo(c *gin.Context) {
	var (
		err  error
		req  validator.CourseWareInfoRuleV2
		data CourseWareInfo
	)

	// 参数校验
	AppId := app.GetAppId(c)
	if err = app.ParseQueryRequest(c, &req); err != nil {
		return
	}
	// req.AliveId = c.Query("alive_id")
	// req.CourseWareId = c.Query("courseware_id")

	if AppId == "" || req.AliveId == "" {
		app.FailWithMessage("内容已被删除1", enums.Code_Db_Not_Find, c)
		return
	}

	//先获取直播表中的ppt_imgs
	aliveReq := course.AliveInfo{AppId: AppId, AliveId: req.AliveId}
	aliveInfo, err := aliveReq.GetAliveInfo()
	if err != nil {
		app.FailWithMessage(fmt.Sprintf("获取直播基础信息错误:%s", err.Error()), enums.ERROR, c)
		return
	}

	//课件id与ppt_img都不存在
	courseWareRep := material.CourseWare{AppId: AppId, AliveId: req.AliveId}
	if req.CourseWareId == "" && aliveInfo.PptImgs.String == "" {
		courseWareInfo, err := courseWareRep.GetCourseWareInfoByAliveId()
		if err != nil {
			app.FailWithMessage(err.Error(), enums.ERROR, c)
			return
		}
		if courseWareInfo.CoursewareImage != "" {
			data.Id = courseWareInfo.Id.String
			data.CoursewareImage = courseWareInfo.CourseImageArray
			data.PageCount = courseWareInfo.PageCount
			data.CurrentPreviewPage = courseWareInfo.CurrentPreviewPage
		}
		app.OkWithData(data, c)
		return
	}

	//有ppt_imgs，优先拿这个
	if aliveInfo.PptImgs.String != "" {
		var coursewareImage []map[string]interface{}
		err = json.Unmarshal([]byte(aliveInfo.PptImgs.String), &coursewareImage)
		if err != nil {
			app.FailWithMessage(fmt.Sprintf("获取直播基础信息错误:%s", err.Error()), enums.ERROR, c)
			return
		}

		//更改新字段返回
		var imageReturnSlice []map[string]interface{}
		for k, v := range coursewareImage {
			imageReturnSlice = append(imageReturnSlice, map[string]interface{}{
				"index":              k,
				"pic_url":            courseWareRep.ReplaceCourseLinkStr(v["pic_url"].(string)),
				"server_id":          v["server_id"],
				"pic_compressed_url": courseWareRep.ReplaceCourseLinkStr(v["pic_url_compressed"].(string)), // v["pic_url_compressed"]
			})
		}
		//赋值到新字段
		data.CoursewareImage = imageReturnSlice
		data.PageCount = len(imageReturnSlice)
		app.OkWithData(data, c)
		return
	}

	//获取课件详情
	courseWareInfo, err := courseWareRep.GetCourseWareInfo(req.CourseWareId)
	if err != nil {
		app.FailWithMessage("获取课件信息错误", enums.ERROR, c)
		return
	}
	if courseWareInfo == nil {
		app.FailWithMessage("获取课件信息错误, courseWareInfo为空", enums.ERROR, c)
		return
	}

	if courseWareInfo.CoursewareImage != "" {
		data.Id = courseWareInfo.Id.String
		data.CoursewareImage = courseWareInfo.CourseImageArray
		data.PageCount = courseWareInfo.PageCount
		data.CurrentPreviewPage = courseWareInfo.CurrentPreviewPage
	}

	//返回
	app.OkWithData(data, c)
}
