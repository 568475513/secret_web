package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"abs/models/business"
	"abs/pkg/enums"
	"abs/pkg/util"
	"abs/service"
)

var bLogger *bussinessLogger

type BuryingPoint struct {
	AppId       string
	UserId      string
	ResourceId  string
	ProductId   string
}

type UserPurchaseData struct {
	AppId          string `json:"app_id"`
	UserId         string `json:"user_id"`
	RawUrl         string `json:"raw_url"`
	Url            string `json:"url"`
	Referer        string `json:"referer"`
	AppVersion     string `json:"app_version"`
	Agent          string `json:"agent"`
	Client         int8
	UserCollection bool   `json:"use_collection"`
	Ip             string `json:"ip"`
	ResourceType   string `json:"resource_type"`
	ResourceId     string `json:"resource_id"`
	ProductId      string `json:"product_id"`
	IsResourcePay  bool   `json:"is_resource_pay"`
	Params         string `json:"params"`
}

type LogData struct {
	*UserPurchaseData
	IsPay           int    `json:"is_pay"`
	CreatedAt       string `json:"created_at"`
	UniversalOpenId string `json:"universal_open_id"`
}

type bussinessLogger struct {
	MaxAge    int
	MaxSize   int
	MaxBackup int
	FileName  string
	Path      string
	*zap.Logger
	rule func() string
}

// 购买关系埋点数据上报
func (b *BuryingPoint) InsertDataUserPurchase(c *gin.Context, available bool) {
	// 接口参数
	params, _ := util.JsonEncode(c.Request.Form)
	userPurchase := &UserPurchaseData{
		AppId:          b.AppId,
		UserId:         b.UserId,
		RawUrl:         c.Request.URL.String(),
		Url:            c.Request.URL.Path,
		Referer:        c.Request.Referer(),
		AppVersion:     c.GetString("app_version"),
		Agent:          c.Request.UserAgent(),
		Client:         int8(c.GetInt("client")),
		UserCollection: true, // 小程序的
		Ip:             c.ClientIP(),
		ResourceType:   strconv.Itoa(enums.ResourceTypeLive),
		ResourceId:     b.ResourceId,
		ProductId:      b.ProductId,
		IsResourcePay:  available,
		Params:         string(params),
	}

	b.InsertUserPurchaseLog(userPurchase)
}

// 写入逻辑
func (b *BuryingPoint) InsertUserPurchaseLog(field *UserPurchaseData) {
	createdAt, isPlay := time.Now().Format(util.TIME_LAYOUT), 0
	userService := service.UserService{AppId: field.AppId, UserId: field.UserId}
	// 查用户了...
	userInfoResponse, _ := userService.RequestUserInfo()
	if isExist, err := business.ExistPurchaseRecord(field.AppId, field.UserId); isExist && err == nil {
		isPlay = 1
	}
	logData := LogData{
		IsPay:            isPlay,
		CreatedAt:        createdAt,
		UniversalOpenId:  userInfoResponse.Data.UniversalOpenId,
		UserPurchaseData: field,
	}
	jsonLog, _ := json.Marshal(logData)
	b.initBussinessLogger(1, 3, 1024*1024, func() string {
		return fmt.Sprintf("user_purchase_log_%s.log", time.Now().Format("2006_01_02"))
	}, util.GetRuntimeDir())
	bLogger.SuperInfo(string(jsonLog))
}

// 私有方法 ==============================
// 写入具体逻辑
func (b *BuryingPoint) initBussinessLogger(maxAge int, maxBackup int, maxSize int, rule func() string, path string) {
	if bLogger != nil {
		return
	}
	bLogger = &bussinessLogger{}
	bLogger.MaxAge = maxSize
	bLogger.MaxSize = maxSize
	bLogger.MaxBackup = maxBackup
	bLogger.Path = path
	bLogger.rule = rule
	bLogger.FileName = BornFileName(rule)
	bLogger.Logger = initZapLogeer(bLogger.MaxSize, bLogger.MaxBackup, bLogger.MaxAge, path + "/" + bLogger.FileName)
}

func (*bussinessLogger) SuperInfo(msg string) {
	currentFileName := BornFileName(bLogger.rule)
	if bLogger.FileName != currentFileName {
		fmt.Println("生成了新的实例.....", currentFileName)
		bLogger.Logger = initZapLogeer(bLogger.MaxSize, bLogger.MaxBackup, bLogger.MaxAge, bLogger.Path+"/"+currentFileName)
		bLogger.FileName = currentFileName
	}
	bLogger.Info(msg)
}

func initZapLogeer(maxSize int, maxBackup int, maxAge int, fileName string) *zap.Logger {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		LocalTime:  true,
	}
	encoder := func() zapcore.Encoder {
		return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
		})
	}
	core := zapcore.NewCore(encoder(), zapcore.AddSync(lumberJackLogger), zap.InfoLevel)
	return zap.New(core)
}

func BornFileName(rule func() string) string {
	return rule()
}