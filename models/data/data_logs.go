package data

type GroupIm struct {
	AppId   string `json:"app_id"`
	AliveId string `json:"alive_id"`
	RoomId  string `json:"room_id"`
	GroupId string `json:"group_id"`
}

func InsertImGroupIdRecord(dataGroup GroupIm) error {
	result := db.Table("db_ex_statistics.old_im_room_change").Create(&dataGroup)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
