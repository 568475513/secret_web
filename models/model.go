package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	"abs/models/alive"
	"abs/models/business"
	"abs/models/data"
	"abs/models/sub_business"
	// "abs/models/user"
)

// 表前缀
const TablePrefix = "t_"

// 初始化各数据库连接
func Init() {
	fmt.Println(">开始初始化各数据库连接池...")

	// 初始化连接池
	// 用户的暂时不要查了，走用户服务
	// user.Init()
	alive.Init()
	business.Init()
	business.InitRw()
	sub_business.Init()
	// data.Init()

	// 设置表名【注意所有数据库链接都会通用这个方法】
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if strings.Contains(defaultTableName, ".") {
			return defaultTableName
		} else {
			return TablePrefix + defaultTableName
		}
	}
	fmt.Println(">>>初始化数据库连接池完成")
}

// 初始化各数据库连接
func InitJob() {
	fmt.Println(">开始初始化Job各数据库连接池...")

	// 初始化连接池
	business.Init()
	business.InitRw()
	data.Init()

	// 设置表名【注意所有数据库链接都会通用这个方法】
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if strings.Contains(defaultTableName, ".") {
			return defaultTableName
		} else {
			return TablePrefix + defaultTableName
		}
	}
	fmt.Println(">>>初始化Job数据库连接池完成")
}

// 关闭各数据库链接
func CloseDB() {
	// user.CloseDB()
	alive.CloseDB()
	business.CloseDB()
	business.CloseDB()
	sub_business.CloseDB()
	data.CloseDB()
}
