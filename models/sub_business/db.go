package sub_business

import (
	"abs/pkg/logging"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 数据库实例
var db *gorm.DB

const (
	DataBase = "db_ex_banned"
	// 其它库
	DatabaseShopConfig    = "db_ex_shop_config"
	DatabaseMicroPage     = "db_micro_page"
	DatabaseSvip          = "db_ex_svip"
	DatabaseContentMarket = "db_ex_content_market"
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	maxOpenConns = 35
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	maxIdleConns = 16
	// 可以重用连接的最长时间[5分钟先]
	maxLifetime = 180
)

// 初始化数据库连接
func Init() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_SUB_RW_USERNAME"),
		os.Getenv("DB_SUB_RW_PASSWORD"),
		os.Getenv("DB_SUB_RF_HOST"),
		os.Getenv("DB_SUB_RW_PORT"),
		DataBase))

	if err != nil {
		logging.Error(fmt.Sprintf("SubBusiness models.Init err: %v", err.Error()))
		//log.Fatalf("SubBusiness models.Init err: %v", err)
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetConnMaxLifetime(time.Second * maxLifetime)

	// 日志[生产必须关闭！]
	if os.Getenv("RUNMODE") == "debug" {
		db.LogMode(true)
	}
	// db.SetLogger(log.New(os.Stdout, "\r\n", 0))
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}

//"github.com/jmoiron/sqlx"
// _ "github.com/go-sql-driver/mysql"
// //初始化数据库连接 sqx
// func Init() {
// 	var err error
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
// 		os.Getenv("DB_SUB_RW_USERNAME"),
// 		os.Getenv("DB_SUB_RW_PASSWORD"),
// 		os.Getenv("DB_SUB_RW_HOST"),
// 		os.Getenv("DB_SUB_RW_PORT"),
// 		DataBase)

// 	// 也可以使用MustConnect连接不成功就panic
// 	db, err = sqlx.Connect("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Connect sub_business DB failed, err:%v\n", err)
// 	}

//  db.SetMaxOpenConns(maxOpenConns)
//  db.SetMaxIdleConns(maxIdleConns)
// }
