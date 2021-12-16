package business

type EBook struct {
	Model
}

// CountEBook 资源计数
func CountEBook(appId string, resourceIds []string, startAt string) (total int, err error) {
	err = db.Table("t_ebook").
		Where("app_id = ? and id in (?) and start_at < ? and state = 0", appId, resourceIds, startAt).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
