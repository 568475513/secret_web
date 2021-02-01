package business

import (
	// "database/sql"
	// "time"

	"github.com/jinzhu/gorm"
)

type VideoMiddleTranscode struct {
	AppId                string  `json:"app_id"`
	FileId               string  `json:"file_id"`
	VideoUrl             string  `json:"video_url"`
	VideoAudioUrl        string  `json:"video_audio_url"`
	VideoMp4             string  `json:"video_mp4"`
	VideoMp4High         string  `json:"video_mp4_high"`
	VideoMp4Size         float64 `json:"video_mp4_size"`
	VideoMp4HighSize     float64 `json:"video_mp4_high_size"`
	VideoMp4Vbitrate     int     `json:"video_mp4_vbitrate"`
	VideoMp4HighVbitrate int     `json:"video_mp4_high_vbitrate"`
	VideoHls             string  `json:"video_hls"`
	VideoSize            float64 `json:"video_size"`
	VideoLength          int     `json:"video_length"`
	M3u8Url              string  `json:"m3u8url"`
	SourceType           uint8   `json:"source_type"`
}

// 设置表名 VideoMiddleTranscode
func (VideoMiddleTranscode) TableName() string {
	return DataBase + ".t_video_middle_transcode"
}

// 获取直播视频转码数据
func GetVideoMiddleTranscode(fileId string) (*VideoMiddleTranscode, error) {
	var ac VideoMiddleTranscode
	err := db.Select([]string{
		"app_id",
		"file_id",
		"video_url",
		"video_audio_url",
		"video_mp4",
		"video_mp4_high",
		"video_mp4_size",
		"video_mp4_high_size",
		"video_mp4_vbitrate",
		"video_mp4_high_vbitrate",
		"video_hls",
		"video_size",
		"video_length",
		"m3u8url",
		"source_type",
		"updated_at"}).Where("file_id=?", fileId).First(&ac).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ac, nil
}
