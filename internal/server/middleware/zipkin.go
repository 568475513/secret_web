package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	zkOt "github.com/openzipkin-contrib/zipkin-go-opentracing"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go"

	"abs/pkg/logging"
	pzipkin "abs/pkg/provider/zipkin"
)

// 全局设置【serviceName string AND ip string】
const (
	// 显然我是不知道这玩意到底有啥用处的（说的是配置2333）
	appName = "abs-go"
	hostPort = "127.0.0.1:8888"
	// 运行前需要修改endpointUrl的值，从https://tracing-analysis.console.aliyun.com/ 获取zipkin网关
	enpoitUrl = "http://tracing-analysis-dc-hz.aliyuncs.com/adapt_gtlc5mrrui@e8b746cb9ebbeb8_gtlc5mrrui@53df7ad2afe8301/api/v2/spans"
)

// 调用链路初始化处理Tracer【文件写入模式】
func ZipkinTracer(isAlibaba bool) gin.HandlerFunc {
	// create a reporter to be used by the tracer
	var reporter reporter.Reporter
	if isAlibaba {
		// 【阿里爸爸解决方案~】
		reporter = httpreporter.NewReporter(enpoitUrl)
	} else {
		// 【文件写入模式~】
		reporter = pzipkin.NewReporter(logging.ZLogger)
	}
	// defer reporter.Close()
	// set-up the local endpoint for our service
	endpoint, err := zipkin.NewEndpoint(appName, hostPort)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	// set-up our sampling strategy
	// sampler := zipkin.NewModuloSampler(1)
	sampler, err := zipkin.NewBoundarySampler(float64(0.2), 2)
	if err != nil {
		log.Fatalf("[采样率]set-up our sampling strategy err: %+v\n", err)
	}
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