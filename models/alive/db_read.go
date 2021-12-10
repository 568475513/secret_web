package alive

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 只数据库实例
var readDb *gorm.DB


// InitReadDb 初始化只读 数据库连接
func InitReadDb() {
	var err error
	readDb, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_ALIVE_RW_USERNAME"),
		os.Getenv("DB_ALIVE_RW_PASSWORD"),
		os.Getenv("DB_ALIVE_RF2_HOST"),
		os.Getenv("DB_ALIVE_RF2_PORT"),
		DataBase))

	if err != nil {
		log.Fatalf("Alive models.Init err: %v", err)
	}

	readDb.SingularTable(true)
	readDb.DB().SetMaxIdleConns(maxIdleConns)
	readDb.DB().SetMaxOpenConns(maxOpenConns)
	readDb.DB().SetConnMaxLifetime(time.Second * maxLifetime)

	// 日志[生产必须关闭！]
	if os.Getenv("RUNMODE") == "debug" {
		readDb.LogMode(true)
	}
}

// CloseReadDB closes database connection (unnecessary)
func CloseReadDB() {
	defer readDb.Close()
}