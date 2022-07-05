package secret

import (
	"abs/internal/server/repository/user"
	"abs/pkg/app"
	"github.com/gin-gonic/gin"
)

//缓存用户数据脚本
func WeekUserDataScript(c *gin.Context) {
	var (
		err error
		u   user.User
	)
	err = u.WeekGetUserData()
	app.OkWithData(err, c)
}

// 发送每日用户周报
func DayUserDataScript(c *gin.Context) {
	var (
		err error
		u   user.User
	)
	err = u.GetUserPriceDay()
	app.OkWithData(err, c)
}

// 发送每日用户日报
func DayUserDataScript2(c *gin.Context) {
	var (
		err error
		u   user.User
	)
	err = u.GetUserPriceDay2()
	app.OkWithData(err, c)
}

// 发送每周日用户周报推送
func DayUserDataPushScript(c *gin.Context) {
	var (
		err error
		u   user.User
	)
	err = u.GetUserDataWeekPush()
	app.OkWithData(err, c)
}
