package business

type Audio struct {
	Model
}

// CountAudio 资源计数
func CountAudio(appId string, resourceIds []string, startAt string) (total int, err error) {
	err = db.Table("t_audio").
		Where("app_id = ? and id in (?) and start_at < ? and audio_state = 0", appId, resourceIds, startAt).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
