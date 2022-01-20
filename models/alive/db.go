package alive

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 数据库实例
var db *gorm.DB

type Model struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// https://blog.csdn.net/qq_39384184/article/details/103954821 不懂的可以看这篇文章
const (
	DataBase = "db_ex_alive"
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	maxOpenConns = 320
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	maxIdleConns = 32
	// 可以重用连接的最长时间[5分钟先]
	maxLifetime = 180
)

// 初始化数据库连接
func Init() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_ALIVE_RW_USERNAME"),
		os.Getenv("DB_ALIVE_RW_PASSWORD"),
		os.Getenv("DB_ALIVE_RW_HOST"),
		os.Getenv("DB_ALIVE_RW_PORT"),
		DataBase))

	if err != nil {
		log.Fatalf("Alive models.Init err: %v", err)
	}

	db.SingularTable(true)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetConnMaxLifetime(time.Second * maxLifetime)

	// 日志[生产必须关闭！]
	if os.Getenv("RUNMODE") == "debug" {
		db.LogMode(true)
	}
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}

// updateTimeStampForUpdateCallback will set `UpdatedAt` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdatedAt", time.Now().Format("2006-01-02 15:04:05"))
	}
}
