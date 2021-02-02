package alive

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"abs/pkg/provider/json"
)

type CourseWareRecords struct {
	AppId              string              `json:"app_id"`
	AliveId            string              `json:"alive_id"`
	AliveTime          int                 `json:"alive_time"`
	CourseUseTime      int                 `json:"course_use_time"`
	UserId             json.JSONNullString `json:"user_id"`
	CoursewareId       json.JSONNullString `json:"courseware_id"`
	CurrentPreviewPage int                 `json:"current_preview_page"`
	CurrentImageUrl    json.JSONNullString `json:"current_image_url"`
	CutFileId          int                 `json:"cut_file_id"`
	IsShow             uint8               `json:"is_show"`
	CreatedAt          json.JSONTime       `json:"created_at"`
	UpdatedAt          json.JSONTime       `json:"updated_at"`
}

type CourseWare struct {
	Id                 json.JSONNullString `json:"id"`
	AppId              string              `json:"app_id"`
	AliveId            string              `json:"alive_id"`
	FileName           string              `json:"filename"`
	FileUrl            string              `json:"file_url"`
	FileType           uint8               `json:"file_type"`
	UseState           uint8               `json:"use_state"`
	ChangeToImageState uint8               `json:"change_to_image_state"`
	PageCount          int                 `json:"page_count"`
	State              int                 `json:"state"`
	CurrentPreviewPage int                 `json:"current_preview_page"`
	CoursewareImage    string              `json:"courseware_image"`
	CourseImageArray   []map[string]interface{}
}

// 获取课件详情（通过courseWareId）
func GetCourseWareInfo(appId, aliveId string, courseWareId, s []string) ([]*CourseWare, error) {
	var cw []*CourseWare

	courseWareIds := strings.Join(courseWareId, ",") //用逗号,拼接

	err := db.Table("t_courseware").
		Select(s).
		Where("app_id=? and alive_id=? and id in (?)", appId, aliveId, courseWareIds).
		Find(&cw).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cw, nil
}

// 获取课件详情（通过aliveId）
func GetCourseWareInfoByAliveId(appId, aliveId string, s []string) (*CourseWare, error) {
	var cw CourseWare

	err := db.Table("t_courseware").
		Select(s).
		Where("app_id=? and alive_id=? and use_state=?", appId, aliveId, 1).
		First(&cw).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &cw, nil
}

// 通过直播时间获取课件记录
func GetCourseWareByAliveTime(appId, aliveId string, lookBackId, aliveTime, pageSize int, ascOrder bool) ([]*CourseWareRecords, error) {
	var cwr []*CourseWareRecords

	//字段
	s := []string{"app_id", "alive_id", "alive_time", "course_use_time", "user_id", "current_preview_page", "current_image_url", "courseware_id", "created_at", "updated_at"}
	where := "app_id = ? and alive_id = ? and is_show = ? and cut_file_id = ? and user_id is not null"
	conn := db.Table("t_alive_course_records").Select(s)

	if ascOrder {
		where += " and course_use_time >= ?"
		err := conn.Where(where, appId, aliveId, 1, lookBackId, aliveTime).
			Order("course_use_time ASC").
			Order("id ASC").
			Limit(pageSize).
			Find(&cwr).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	} else { //更加course_use_time从0到对应值内的课件
		where += " and courseware_id is not null and course_use_time >= ? and course_use_time <= ?"
		err := conn.Where(where, appId, aliveId, 1, lookBackId, 0, aliveTime).
			Order("course_use_time DESC").
			Order("id DESC").
			Limit(pageSize).
			Find(&cwr).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return cwr, nil
}

// 全部更新为不使用状态
func UpdateCourseWareAllNotUseState(appId, aliveId string) error {
	data := make(map[string]string)
	data["use_state"] = "0"
	data["updated_at"] = time.Now().Format("2006-01-02 15:06:04")

	err := db.Table("t_courseware").Where("app_id=? and alive_id=? and use_state=?", appId, aliveId, 1).
		Update(data).Error
	if err != nil {
		return err
	}

	return nil
}

// 更新课件使用状态
func UpdateCourseWareUseState(appId, aliveId, courseWareId string, data map[string]string) error {
	err := db.Table("t_courseware").
		Where("app_id=? and alive_id=? and id=?", appId, aliveId, courseWareId).
		Update(data).Error
	if err != nil {
		return err
	}

	return nil
}
