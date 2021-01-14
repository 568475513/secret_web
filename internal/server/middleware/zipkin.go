package middleware

import (
	// "fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	zkOt "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"

	"abs/pkg/logging"
	pzipkin "abs/pkg/provider/zipkin"
)

// 全局设置
// serviceName string AND ip string
const (
	// 显然我是不知道这玩意到底有啥用处的
	appName = "abs-go"
	hostPort = "127.0.0.1:8888"
)

// 调用链路初始化处理Tracer
func ZipkinTracer() gin.HandlerFunc {
	// create a reporter to be used by the tracer
	reporter := pzipkin.NewReporter(logging.ZLogger)
	// defer reporter.Close()
	// set-up the local endpoint for our service
	endpoint, err := zipkin.NewEndpoint(appName, hostPort)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	// set-up our sampling strategy
	sampler := zipkin.NewModuloSampler(1)
	// initialize the tracer
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	// 设置全局
	zkTracer := zkOt.Wrap(tracer)
	opentracing.SetGlobalTracer(zkTracer)

	// 返回中间件处理
	return func(c *gin.Context) {
		// 实际记录
		span := zkTracer.StartSpan(c.FullPath())
		defer span.Finish()
		// 设置父级Span
		c.Set("ParentSpanContext", span.Context())
		// 往下处理请求
		c.Next()
	}
}
