package business

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

// 读写实例
var dbRw *gorm.DB

type Model struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

const (
	// 默认链接库，有些模型里面需要设置库的
	DataBase = "db_ex_business"
	// 其它库
	DatabaseConfig = "db_ex_config"
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	maxOpenConns = 36
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	maxIdleConns = 32
	// 可以重用连接的最长时间[5分钟先]
	// 防止closing bad idle connection: EOF（ 在 MySQL Server 主动断开连接之前，MySQL Client 的连接池中的连接被关闭掉），具体值要问DBA
	// 数据库端设置生存时间30s
	maxLifetime = 28
)

// 初始化数据库连接
func Init() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_CORE_RW_USERNAME"),
		os.Getenv("DB_CORE_RW_PASSWORD"),
		os.Getenv("DB_CORE_RF_HOST"),
		os.Getenv("DB_CORE_RW_PORT"),
		DataBase))

	if err != nil {
		logging.Error(fmt.Sprintf("Business models.Init err: %v", err.Error()))
		//log.Fatalf("Business models.Init err: %v", err)
	}

	db.SingularTable(true)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetConnMaxLifetime(time.Second * maxLifetime)

	// 日志[生产必须关闭！]
	if os.Getenv("RUNMODE") == "debug" {
		db.LogMode(true)
		// db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	}
}

// 初始化读写数据库连接
func InitRw() {
	var err error
	dbRw, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_CORE_RW_USERNAME"),
		os.Getenv("DB_CORE_RW_PASSWORD"),
		os.Getenv("DB_CORE_RW_HOST"),
		os.Getenv("DB_CORE_RW_PORT"),
		DataBase))

	if err != nil {
		logging.Error(fmt.Sprintf("Business model.InitRw err: %v", err.Error()))
		//log.Fatalf("Business model.InitRw err: %v", err)
	}

	dbRw.SingularTable(true)
	// 这个有问题，会不生效的
	// dbRw.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	dbRw.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	dbRw.DB().SetMaxIdleConns(maxIdleConns)
	dbRw.DB().SetMaxOpenConns(maxOpenConns)
	dbRw.DB().SetConnMaxLifetime(time.Second * maxLifetime)

	// 日志[生产必须关闭！]
	if os.Getenv("RUNMODE") == "debug" {
		dbRw.LogMode(true)
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

// updateTimeStampForCreateCallback will set `CreatedAT`, `UpdatedAt` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedAt"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdatedAt"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}
