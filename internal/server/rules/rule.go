package rules

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var v *validator.Validate
// 自定义校验
// use a single instance of Validate, it caches struct info
// var Validate *validator.Validate

// 初始化系统验证器
func InitVali() {
	// Validate = validator.New()
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 自定义手机号码验证
		v.RegisterValidation("checkMobile", checkMobile)
		// 自定义pay_info验证【base_info】1
		v.RegisterValidation("checkPayInfo1", checkPayInfoBaseInfo)
	}
	fmt.Println(">>>初始化系统验证器完成")
}

// 以下是公共常见的校验写在这个文件，其它特殊校验分拆文件
// 手机号码验证方法
func checkMobile(fl validator.FieldLevel) bool {
	mobile := strconv.Itoa(int(fl.Field().Uint()))
	re := `^1[3456789]\d{9}$`
	r := regexp.MustCompile(re)
	return r.MatchString(mobile)
}

// 参数验证兼容baseinfo的pay_info
func checkPayInfoBaseInfo(fl validator.FieldLevel) bool {
	// 能用数组就不要用切片
	fields := [6]string{"type", "app_id", "resource_id", "resource_type", "payment_type", "product_id"}
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(fl.Field().String()), &dat); err != nil {
		return false
	}
	for _, v := range fields {
		if _, ok := dat[v]; !ok {
			return false
		}
	}
	// 不能都为空
	if dat["resource_id"].(string) == "" && dat["product_id"].(string) == "" {
		return false
	}

	// content_app_id限制
	if v, ok := dat["content_app_id"]; ok && len(v.(string)) < 3 {
		return false
	}

	return true
}
