package alive

import (
	"github.com/jinzhu/gorm"
)

type AliveLookBack struct {
	Id             int    `json:"id"`
	AppId          string `json:"app_id"`
	AliveId        string `json:"alive_id"`
	LookbackFileId string `json:"lookback_file_id"`
	RegionFileId   string `json:"region_file_id"`
	LookbackMp4    string `json:"lookback_mp4"`
	LookbackM3u8   string `json:"lookback_m3u8"`
	FileName       string `json:"file_name"`
	TranscodeState uint8  `json:"transcode_state"`
	State          uint8  `json:"state"`
	OriginType     uint8  `json:"origin_type"`
}

type AliveConcatHlsResult struct {
	ChannelId                string `json:"channel_id"`
	LatestM3u8FileId         string `json:"latest_m3u8_file_id"`
	ConcatLatestFileId       string `json:"concat_latest_file_id"`
	ConcatM3u8Url            string `json:"concat_m3u8_url"`
	TranscodeState           uint8  `json:"transcode_state"`
	TranscodeSuccessLastTime string `json:"transcode_success_last_time"`
	ConcatSuccessLastTime    string `json:"concat_success_last_time"`
	TranscodeM3u8Url         string `json:"transcode_m3u8_url"`
	ConcatTimes              uint8  `json:"concat_times"`
	TranscodeTimes           uint8  `json:"transcode_times"`
	ComposeLatestFileId      string `json:"compose_latest_file_id"`
	ConcatMp4Url             string `json:"concat_mp4_url"`
	IsUseConcatMp4           uint8  `json:"is_use_concat_mp4"`
	IsDrm                    uint8  `json:"is_drm"`
	DrmM3u8Url               string `json:"drm_m3u8_url"`
}

func GetAliveLookBackFile(appId string, aliveId string, s []string) (*AliveLookBack, error) {
	var alf AliveLookBack
	err := db.Table("t_alive_lookback").Select(s).Where("app_id=? and alive_id=? and state=?", appId, aliveId, 1).First(&alf).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &alf, nil
}

func GetAliveHlsResult(channelId string, s []string) (*AliveConcatHlsResult, error) {
	var ahr AliveConcatHlsResult

	err := db.Select(s).Where("channel_id=?", channelId).First(&ahr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &ahr, nil
}