package alive

import "github.com/jinzhu/gorm"

type AliveImMiddler struct {
	AppId     string `json:"app_id"`
	AliveId   string `json:"alive_id"`
	OldRoomId string `json:"old_room_id"`
	NewRoomId string `json:"new_room_id"`
}

type AliveImRecord struct {
	AppId   string `json:"app_id"`
	AliveId string `json:"alive_id"`
	RoomId  string `json:"room_id"`
	GroupId string `json:"group_id"`
}

// 更新t_alive room_id
func UpdateTAliveRommId(appId, aliveId, roomId string) error {
	var a Alive
	return db.Model(&a).Where("app_id=? and id=? ", appId, aliveId).
		Update("room_id", roomId).
		Limit(1).Error
}

// 更新禁言表room_id
func UpdateForbidRoomId(appId, roomId, newRoomId string) error {
	return db.Table("t_alive_forbid").Where("app_id=? and room_id=? ", appId, roomId).
		Update("room_id", newRoomId).
		Limit(1).Error
}

// 插入表t_alive_im_middle room_id
func InsertImMiddle(aim AliveImMiddler) error {
	result := db.Table("t_alive_im_middle").Create(&aim)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 插入表t_record_im_change room_id
func InsertImGroupIdRecord(aim AliveImRecord) error {
	result := db.Table("t_record_im_change").Create(&aim)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 通过t_alive_im_middle获取直播room_id
func GetRoomIdByAliveId(appId, aliveId, field string) (*AliveImMiddler, error) {
	var a AliveImMiddler
	err := db.Table("t_alive_im_middle").Select(field).Where("app_id=? and alive_id=? ", appId, aliveId).First(&a).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &a, nil
}
