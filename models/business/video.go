package business

type Video struct {
	Model
}

// CountVideo 资源计数
func CountVideo(appId string, resourceIds []string, startAt string) (total int, err error) {
	err = db.Table("t_video").
		Where("app_id = ? and id in (?) and start_at < ? and video_state = 0", appId, resourceIds, startAt).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
