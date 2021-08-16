package app

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
)

// 获取调用链路父级初始化
func GetTracingSpan(c *gin.Context) opentracing.SpanContext {
	parentSpanContext, isExists := c.Get("ParentSpanContext")
	if isExists {
		return parentSpanContext.(opentracing.SpanContext)
	} else {
		// 容错处理！
		return opentracing.StartSpan("ParentSpanContextError").Context()
	}
}

// Query参数验证
func ParseQueryRequest(c *gin.Context, request interface{}) (err error) {
	err = c.ShouldBindQuery(request)
	handleParamsError(c, err)
	return
}

// 普通参数验证
func ParseRequest(c *gin.Context, request interface{}) (err error) {
	err = c.ShouldBind(request)
	handleParamsError(c, err)
	return
}

// 处理422错误信息
func handleParamsError(c *gin.Context, err error) {
	if err != nil {
		var errStr string
		switch err.(type) {
		case validator.ValidationErrors:
			errStr = err.Error()
		case *json.UnmarshalTypeError:
			unmarshalTypeError := err.(*json.UnmarshalTypeError)
			errStr = fmt.Errorf("%s [类型错误，期望类型] %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
		default:
			// logging.Error("unknown error：" + err.Error())
			errStr = errors.New("参数校验[unknown error]").Error()
		}
		// 此处错误信息该修正下
		FailWithParameter(errStr, c)
	}
}

// 没有用户标识参数上传的错误
func NoUserParseRequest(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	universalUnionId := c.Query("universal_union_id")
	if userId == "" && universalUnionId == "" {
		FailWithParameter("用户参数未上传（user_id或者universal_union_id）", c)
		return errors.New("用户参数未上传（user_id或者universal_union_id）")
	}
	return
}

// 获取AppId
func GetAppId(c *gin.Context) string {
	appId := c.GetString("app_id")
	if appId == "" {
		// 兼容非网关
		appId = c.DefaultQuery("app_id", c.DefaultPostForm("app_id", ""))
	}
	return appId
}

// 获取UserId
func GetUserId(c *gin.Context) string {
	userId := c.GetString("user_id")
	anonUserId := c.GetString("anon_user_id")
	if len(userId) > 0 {
		return userId
	} else {
		if len(anonUserId) > 0 {
			return anonUserId
		} else {
			userInfo := c.GetStringMap("userInfo")
			if userId, ok := userInfo["user_id"]; ok {
				return userId.(string)
			} else {
				return ""
			}
		}
	}
}
