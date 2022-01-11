module abs

go 1.14

require (
	github.com/DeanThompson/ginpprof v0.0.0-20201112072838-007b1e56b2e1
	github.com/RichardKnop/machinery v1.10.0
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.62.0
	github.com/go-playground/validator/v10 v10.2.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.9
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/spf13/cobra v1.1.1
	github.com/tencentyun/tls-sig-api-v2-golang v1.1.0
	go.uber.org/zap v1.10.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/gomodule/redigo => talkcheap.xiaoeknow.com/AliveDev/redigo v1.8.4
