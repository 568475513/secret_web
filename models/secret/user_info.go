package user

import (
	"abs/pkg/enums"
	"time"

	"github.com/jinzhu/gorm"
)

// 用户信息结构体
type SecretUser struct {
	Model

	Id                   string    `json:"id"`
	UserId               string    `json:"user_id"`
	RegisterId           string    `json:"register_id"`
	UserIp               string    `json:"user_ip"`
	UserDnsPreventDomain string    `json:"user_dns_prevent_domain"`
	UserPrice            float64   `json:"user_price"`
	UserName             string    `json:"user_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type UId struct {
	UserId               string  `json:"user_id"`
	UserName             string  `json:"user_name"`
	UserDnsPreventDomain string  `json:"user_dns_prevent_domain"`
	RegisterId           string  `json:"register_id"`
	UserPrice            float64 `json:"user_price"`
}

type UserConf struct {
	UserId        string `json:"user_id"`
	IsBusMonitor  int    `json:"is_bus_monitor"`
	IsLargeData   int    `json:"is_large_data"`
	IsSpy         int    `json:"is_spy"`
	IsCollectInfo int    `json:"is_collect_info"`
	ExpiredAt     string `json:"expired_at"`
}

type UConf struct {
	IsBusMonitor  int `json:"is_bus_monitor"`
	IsLargeData   int `json:"is_large_data"`
	IsSpy         int `json:"is_spy"`
	IsCollectInfo int `json:"is_collect_info"`
}

// 获取用户信息
func GetUserInfo(userId, userIp string) (su *SecretUser, err error) {

	var ui SecretUser
	if userId != "" {
		err = db.Table("t_secret_user").Where("user_id=?", userId).First(&ui).Error
	} else if userIp != "" {
		err = db.Table("t_secret_user").Where("user_ip=?", userIp).First(&ui).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &ui, nil
}

//注册用户信息
func RegisterUser(userId, userDPD, registerId string, price float64) (err error) {
	var ui UId
	ui.UserId = userId
	ui.UserDnsPreventDomain = userDPD
	ui.RegisterId = registerId
	ui.UserPrice = price
	err = db.Table("t_secret_user").Create(ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

//注册用户信息2.0
func RegisterUserV2(userId, userDPD, registerId string) (err error) {
	var (
		ui  UId
		uid UserConf
	)
	ui.UserId = userId
	ui.UserDnsPreventDomain = userDPD
	ui.RegisterId = registerId
	err = db.Table("t_secret_user").Create(ui).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	uid.UserId = userId
	uid.ExpiredAt = time.Now().Format("2006-01-02 15:04:05")
	err = db.Table("t_secret_user_config").Create(uid).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	return
}

//获取所有用户配置列表
func GetAllUserConfigList() (UcList map[string][]string, err error) {

	var uc UserConf
	rs, err := db.Table("t_secret_user_config").Select("user_id, is_bus_monitor, is_large_data, is_spy, is_collect_info").Where("expired_at > ?", time.Now()).Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	UcList = make(map[string][]string)
	if rs != nil {
		for rs.Next() {
			rs.Scan(&uc.UserId, &uc.IsBusMonitor, &uc.IsLargeData, &uc.IsSpy, &uc.IsCollectInfo)
			if uc.IsSpy == 1 {
				UcList[uc.UserId] = append(UcList[uc.UserId], enums.IsSpy)
			}
			if uc.IsCollectInfo == 1 {
				UcList[uc.UserId] = append(UcList[uc.UserId], enums.IsCollectInfo)
			}
			if uc.IsLargeData == 1 {
				UcList[uc.UserId] = append(UcList[uc.UserId], enums.IsLargeData)
			}
			if uc.IsBusMonitor == 1 {
				UcList[uc.UserId] = append(UcList[uc.UserId], enums.IsBusMonitor)
			}
		}
	}
	return
}

//更新用户积分
func UpdateUserPrice(userId, userIp string, price float64) (err error) {

	if userId != "" {
		err = db.Table("t_secret_user").Where("user_id=?", userId).Update("user_price", price).Error
	} else if userIp != "" {
		err = db.Table("t_secret_user").Where("user_ip=?", userIp).Update("user_price", price).Error
	}
	return
}

//更新用户registerId
func UpdateUserRegisterId(userId, registerID string) (err error) {

	err = db.Table("t_secret_user").Where("user_id=?", userId).Update("register_id", registerID).Error
	return
}

//获取用户id
func GetUserId() (err error, ids []string) {

	var ui UId
	rs, err := db.Table("t_secret_user").Select("user_id").Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if rs != nil {
		for rs.Next() {
			rs.Scan(&ui.UserId)
			ids = append(ids, ui.UserId)
		}
	}
	return
}

//获取用户id和极光推送id
func GetUserIdAndRegisterID() (err error, ids map[string]UId) {

	var ui UId
	ids = make(map[string]UId)
	rs, err := db.Table("t_secret_user").Select("user_id, register_id, user_price").Where("register_id != ''").Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if rs != nil {
		for rs.Next() {
			rs.Scan(&ui.UserId, &ui.RegisterId, &ui.UserPrice)
			ids[ui.UserId] = ui
		}
	}
	return
}
