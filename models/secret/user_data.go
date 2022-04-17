package user

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"time"
)

//用户数据结构体
type UserData struct {
	Model

	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	DomainType int       `json:"domain_type"`
	Domain     string    `json:"domain"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Prevent struct {
	PreventName int `json:"prevent_name"`
	PreventNum  int `json:"prevent_num"`
}

type PreventDomain struct {
	PreventDomain string `json:"prevent_domain"`
	PreventNums   int    `json:"prevent_nums"`
}

type PreventDetail struct {
	PreventDomain string    `json:"domain"`
	DomainTag     string    `json:"domain_tag"`
	DomainSource  string    `json:"domain_source"`
	CreatedAt     time.Time `json:"created_at"`
}

type PreventInfo struct {
	UserId       string `json:"user_id"`
	DomainType   int    `json:"domain_type"`
	DomainTag    string `json:"domain_tag"`
	DomainSource string `json:"domain_source"`
	Domain       string `json:"domain"`
	UserIp       string `json:"user_ip"`
}

type UserWeekData struct {
	DomainType int    `json:"domain_type"`
	Count      int    `json:"count"`
	DomainName string `json:"domain_name"`
}

//获取用户拦截类型数
func GetPreventCountByUserId(userId, userIp string) (tcs []Prevent, err error) {

	var (
		tc Prevent
		rs *sql.Rows
	)
	if userId != "" {
		rs, err = db.Table("t_secret_user_data").Select("domain_type , count(id) as count").Where("user_id = ? ", userId).Group("domain_type").Rows()
	} else if userIp != "" {
		rs, err = db.Table("t_secret_user_data").Select("domain_type , count(id) as count").Where("user_ip = ? ", userIp).Group("domain_type").Rows()
	}
	if err != nil || rs == nil {
		return nil, nil
	}
	for rs.Next() {
		rs.Scan(&tc.PreventName, &tc.PreventNum)
		tcs = append(tcs, tc)
	}
	return
}

//获取用户类型详细数据
func GetPreventDetailByUserId(userId, userIp, dt string, page, page_size int) (ps []PreventDetail, err error) {

	var (
		p  PreventDetail
		rs *sql.Rows
	)
	if userId != "" {
		rs, err = db.Table("t_secret_user_data").Select("domain, created_at, domain_tag, domain_source").Where("user_id = ? and domain_type=? ", userId, dt).Limit(page_size).Offset((page - 1) * page_size).Order("created_at desc").Rows()
	} else if userIp != "" {
		rs, err = db.Table("t_secret_user_data").Select("domain, created_at, domain_tag, domain_source").Where("user_ip = ? and domain_type=? ", userIp, dt).Limit(page_size).Offset((page - 1) * page_size).Order("created_at desc").Rows()
	}

	if err != nil && err != gorm.ErrRecordNotFound || rs == nil {
		return ps, nil
	}
	for rs.Next() {
		rs.Scan(&p.PreventDomain, &p.CreatedAt, &p.DomainTag, &p.DomainSource)
		ps = append(ps, p)
	}
	return ps, nil
}

//记录用户拦截信息
func InsertPreventInfo(userId, userIp, domain, domainTag, domainSource string, domainType int) (err error) {

	p := PreventInfo{UserId: userId, UserIp: userIp, Domain: domain, DomainType: domainType, DomainTag: domainTag, DomainSource: domainSource}
	err = db.Table("t_secret_user_data").Create(p).Error
	return err
}

// 根据id查询用户区间数据
func SelectUserDataTime(userId string) (rs map[string][]UserWeekData, err error) {

	var uw UserWeekData
	rs = map[string][]UserWeekData{}
	t := time.Now().Add(-7 * time.Hour * 24).Format("2006-01-02")
	now := time.Now().Format("2006-01-02")
	re, err := db.Table("t_secret_user_data").Select("domain_type , count(id) as count").Where("user_id = ? and created_at >= ? and created_at < ? ", userId, t, now).Group("domain_type").Rows()
	if err != nil && err != gorm.ErrRecordNotFound || re == nil {
		return
	}
	for re.Next() {
		re.Scan(&uw.DomainType, &uw.Count)
		rs[userId] = append(rs[userId], UserWeekData{Count: uw.Count, DomainType: uw.DomainType})
	}
	return
}
