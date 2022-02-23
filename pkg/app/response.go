package app

import (
	"abs/pkg/conf"
	"abs/pkg/logging"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	e "abs/pkg/enums"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	RequestId string `json:"requestId"`
}

// Response setting gin.JSON
func Result(httpCode, errCode int, errMsg string, data interface{}, c *gin.Context) {

	requestId := c.GetString(conf.AbsRequestId)

	logging.GetLogger().Info("responseData",
		zap.String("appId", GetAppId(c)),
		zap.String("userId", GetUserId(c)),
		zap.String("requestId", requestId),
		zap.Int("status", c.Writer.Status()),
		zap.Int("httpCode", httpCode),
		zap.Int("errCode", errCode),
		zap.String("errMsg", errMsg),
		zap.Any("data", data),
	)

	c.JSON(httpCode, Response{
		Code: errCode,
		Msg:  errMsg,
		Data: data,
		RequestId: requestId,
	})
}

// Response Success
func OK(c *gin.Context) {
	Result(http.StatusOK, e.SUCCESS, "操作成功", map[string]interface{}{}, c)
}

// Response OkWithData
func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, e.SUCCESS, message, map[string]interface{}{}, c)
}

// Response OkWithData
func OkWithData(data interface{}, c *gin.Context) {
	Result(http.StatusOK, e.SUCCESS, "OK", data, c)
}

// Response OkWithCodeData
func OkWithCodeData(message string, data interface{}, code int, c *gin.Context) {
	Result(http.StatusOK, code, message, data, c)
}

// Response Fail
func Fail(c *gin.Context) {
	Result(http.StatusOK, e.ERROR, "操作失败", map[string]interface{}{}, c)
}

// Response FailWithMessage
func FailWithMessage(message string, code int, c *gin.Context) {
	Result(http.StatusOK, code, message, map[string]interface{}{}, c)
}

// Response FailWithMessage
func FailWithParameter(message string, c *gin.Context) {
	Result(http.StatusOK, e.INVALID_PARAMS, message, map[string]interface{}{}, c)
}
