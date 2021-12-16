package business

type ImgText struct {
	Model
}

// CountImgText 资源计数
func CountImgText(appId string, resourceIds []string, startAt string) (total int, err error) {
	err = db.Table("t_image_text").
		Where("app_id = ? and id in (?) and start_at < ? and display_state = 0", appId, resourceIds, startAt).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
