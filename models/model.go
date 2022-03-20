package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	s "abs/models/secret"
)

// 表前缀
const TablePrefix = "t_"

// 初始化各数据库连接
func Init() {
	fmt.Println(">开始初始化各数据库连接池...")
	// 初始化连接池
	s.Init()
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

// 关闭各数据库链接
func CloseDB() {
	s.CloseDB()
}
