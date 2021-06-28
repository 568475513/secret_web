package alive

type AliveImMiddler struct {
	Model

	AppId     string `json:"app_id"`
	AliveId   string `json:"alive_id"`
	OldRoomId string `json:"old_room_id"`
	NewRoomId string `json:"new_room_id"`
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
	var af AliveForbid
	return db.Model(&af).Where("app_id=? and room_id=? ", appId, roomId).
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
