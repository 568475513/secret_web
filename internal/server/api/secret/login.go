package secret

import (
	"abs/internal/server/repository/user"
	"abs/internal/server/rules/validator"
	"abs/pkg/app"
	"abs/pkg/enums"

	"github.com/gin-gonic/gin"
)

//用户登录注册接口
func UserLogin(c *gin.Context) {
	var (
		err error
		req validator.SecretUserLoginRule
		u   user.User
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	//如果用户id不存在则注册，反之则登录获取用户信息
	if req.UserId != "" { //	用户登录
		u.UserId = req.UserId
		_, err := u.GetUserInfo() //获取用户基本信息
		if err != nil {
			app.FailWithMessage("获取用户信息异常", enums.ERROR, c)
		}
	} else { //用户注册
		u.RegisterId = req.RegisterId
		u.GetUserOnlyId()

	}
	app.OkWithData(u, c)
}
