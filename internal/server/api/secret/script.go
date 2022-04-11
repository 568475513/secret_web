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
