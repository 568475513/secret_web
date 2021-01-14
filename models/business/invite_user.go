package business

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type InviteUser struct {
	Model
	Id           int            `json:"id" gorm:"AUTO_INCREMENT"`
	AppId        string         `json:"app_id"`
	PaymentType  int            `json:"payment_type"`
	ResourceId   sql.NullString `json:"resource_id"`
	ResourceType int            `json:"resource_type"`
	ProductId    string         `json:"product_id"`
	ShareUserId  string         `json:"share_user_id"`
	InviteCount  int            `json:"invite_count"`
}

type InviteRelation struct {
	Model
	Id            int            `json:"id" gorm:"AUTO_INCREMENT"`
	AppId         string         `json:"app_id"`
	PaymentType   int            `json:"payment_type"`
	ResourceId    sql.NullString `json:"resource_id"`
	ResourceType  int            `json:"resource_type"`
	ProductId     string         `json:"product_id"`
	ShareUserId   string         `json:"share_user_id"`
	ShareType     int            `json:"share_type"`
	InvitedUserId string         `json:"invited_user_id"`
}

// 获取邀请达人前X名
func GetInviteUserList(appId string, paymentType int, resourceId string, resourceType int, productId string, userId string, limit int) ([]*InviteUser, error) {
	var iu []*InviteUser
	err := db.Select("share_user_id,invite_count").Where(
		"app_id=? and payment_type=? and resource_id=? and resource_type=? and product_id=? and invite_count>0",
		appId,
		paymentType,
		resourceId,
		resourceType,
		productId).Order("invite_count desc,created_at asc").Limit(limit).Find(&iu).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return iu, err
}

// 查询单个邀请用户的关联
func GetInviteUserByInvitedUser(inviteRelation InviteRelation) (*InviteRelation, error) {
	var iu InviteRelation
	err := db.Select("id").Where(
		"app_id=? and payment_type=? and resource_id=? and resource_type=? and product_id=? and invited_user_id=?",
		inviteRelation.AppId,
		inviteRelation.PaymentType,
		inviteRelation.ResourceId,
		inviteRelation.ResourceType,
		inviteRelation.ProductId,
		inviteRelation.InvitedUserId).First(&iu).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &iu, err
}

// 查询单个分享用户的关联
func GetInviteUserByShareUser(inviteUser InviteUser) ([]*InviteUser, error) {
	var iu []*InviteUser
	err := db.Select("share_user_id").Where(
		"app_id=? and share_user_id = ? and payment_type=? and resource_id=? and resource_type=? and product_id=?",
		inviteUser.AppId,
		inviteUser.ShareUserId,
		inviteUser.PaymentType,
		inviteUser.ResourceId,
		inviteUser.ResourceType,
		inviteUser.ProductId).Find(&iu).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return iu, err
}

// 添加一条InviteUser的数据
func SetInviteUser(insert InviteUser) error {
	if dbRw.NewRecord(insert) {
		return db.Create(&insert).Error
	} else {
		return nil
	}
}

// 添加一条InviteRelation的数据
func SetInviteRelation(insert InviteRelation) error {
	if dbRw.NewRecord(insert) {
		return db.Create(&insert).Error
	} else {
		return nil
	}
}

// 修改一条记录的邀请数
func UpdateInviteUserByInviteCount(insert InviteUser) int64 {
	// todo 错误一会抛出来
	return db.Table("t_invite_user"). // .Model(&insert)
						Where(
			"app_id=? and share_user_id = ? and payment_type=? and resource_id=? and resource_type=? and product_id=?",
			insert.AppId,
			insert.ShareUserId,
			insert.PaymentType,
			insert.ResourceId.String,
			insert.ResourceType,
			insert.ProductId).UpdateColumn("invite_count", gorm.Expr("invite_count + ?", 1)).RowsAffected
}
