package routers

import (
	"fmt"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"

	"abs/cmd/server/routers/groups"
	"abs/internal/server/middleware"
	e "abs/pkg/enums"
	"abs/pkg/logging"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	// 收到panics抛错后返回500中间件
	r.Use(middleware.GinRecovery(logging.BLogger))

	// 自定义中间件在此处添加...[注意顺序]
	// 暂存当前上下文
	r.Use(middleware.SetCurrentContext())
	// 跨域中间件
	r.Use(middleware.Cors())
	// 处理公共请求参数以及网关转发参数
	r.Use(middleware.ReqParamHandle())
	// 调用链路zipkin
	r.Use(middleware.ZipkinTracer(true))

	// 请求的入参和出参日志
	r.Use(middleware.GinLogger(logging.BLogger))

	// 此处可写公共路由...
	// 健康检测接口
	r.GET("/health", func(c *gin.Context) {
		c.String(e.SUCCESS, "health - "+fmt.Sprint(time.Now().Unix()))
	})

	// 加载其它路由组
	group := r.Group("")
	groups.AliveBaseRouter(group) // 注册直播基础路由组
	groups.AliveAppRouter(group)  // 注册app接口路由组

	// 性能分析 ...
	// goPprof handel
	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	ginpprof.Wrap(r)

	return r
}
