package service

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	e "abs/pkg/enums"
	"abs/pkg/util"
)

// http请求池配置
const (
	maxIdleConns        = 32
	maxIdleConnsPerHost = 32
	idleConnTimeout     = 10
	dialerTimeout       = 10
	dialerKeepAlive     = 60
)

// http pool
var client http.Client

type XiaoeHttpSettings struct {
	Timeout    time.Duration
	Retries    int
	RetryDelay time.Duration
}

type XiaoeHttpRequest struct {
	url      string
	params   interface{}
	settings XiaoeHttpSettings
	req      *http.Request
	resp     *http.Response
}

// 用户请求连接池
func InitService() {
	// 全局超时
	timeout, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
	// 连接池配置
	client = http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			IdleConnTimeout:     idleConnTimeout * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   dialerTimeout * time.Second,
				KeepAlive: dialerKeepAlive * time.Second,
			}).DialContext,
		},
	}
}

// 创建请求体
func newXiaoeRequest(rawUrl, method string) *XiaoeHttpRequest {
	var resp http.Response
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Println("XiaoeRequest:", err)
	}
	req := http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	return &XiaoeHttpRequest{
		url: rawUrl,
		req: &req,
		settings: XiaoeHttpSettings{
			Timeout: 3 * time.Second,
		},
		resp: &resp,
	}
}

// 发起Get请求
func Get(url string) *XiaoeHttpRequest {
	return newXiaoeRequest(url, "GET")
}

// 发起Post请求
func Post(url string) *XiaoeHttpRequest {
	return newXiaoeRequest(url, "POST")
}

// 设置Host
func (x *XiaoeHttpRequest) SetHost(host string) *XiaoeHttpRequest {
	x.req.Host = host
	return x
}

// 设置超时
func (x *XiaoeHttpRequest) SetTimeout(timeout time.Duration) {
	x.settings.Timeout = timeout
}

// 设置头部信息
func (x *XiaoeHttpRequest) SetHeader(key, value string) *XiaoeHttpRequest {
	x.req.Header.Set(key, value)
	return x
}

// 设置参数
func (x *XiaoeHttpRequest) SetParams(params interface{}) *XiaoeHttpRequest {
	x.params = params
	return x
}

func (x *XiaoeHttpRequest) Setting(settings XiaoeHttpSettings) *XiaoeHttpRequest {
	x.settings = settings
	return x
}

func (x *XiaoeHttpRequest) Request() *http.Request {
	return x.req
}

func (x *XiaoeHttpRequest) Response() (*http.Response, error) {
	return x.getResponse()
}

// 获取响应
func (x *XiaoeHttpRequest) getResponse() (resp *http.Response, err error) {
	if x.resp.StatusCode != 0 {
		return x.resp, nil
	}

	done := make(chan bool)
	go func() {
		defer close(done)
		resp, err = x.doRequest()
		x.resp = resp
	}()

	// 超时控制
	select {
	case <-time.After(x.settings.Timeout):
		data := make(map[string]interface{})
		data["code"] = e.TIMEOUT
		data["msg"] = fmt.Sprintf("url: %s[request timeout: %s]", x.req.URL, x.settings.Timeout)
		byteData, _ := util.JsonEncode(data)
		// 这个Error得改，还不能打开
		// logging.Error(string(byteData))
		resp = new(http.Response)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(byteData))
		resp.Body.Close()
		return
	case <-done:
		return
	}

	// 1.配置config
	// configA := hystrix.CommandConfig{
	// 	Timeout: x.settings.Timeout,
	// }
	// 2.配置command
	// hystrix.ConfigureCommand("getResponse", configA)
	// 3.执行Do方法
	// err = hystrix.Do("getResponse", func() error {
	// 	resp, err = x.DoRequest()
	// 	return err
	// }, func(e error) error {
	// 	data := make(map[string]string)
	// 	data["code"] = strconv.Itoa(e.TIMEOUT)
	// 	data["msg"] = x.url + "[request timeout]"
	// 	data["data"] = ""
	// 	byteData, _ := json.Marshal(data)
	// 	resp = new(http.Response)
	// 	resp.Body = ioutil.NopCloser(bytes.NewBuffer(byteData))
	// 	resp.Body.Close()
	// 	return
	// })
}

// 发起请求
func (x *XiaoeHttpRequest) doRequest() (resp *http.Response, err error) {
	if x.req.Method == "GET" {
		x.buildUrl(x.params)
		urlParsed, err := url.Parse(x.url)
		if err != nil {
			return nil, err
		}
		x.req.URL = urlParsed
	} else if x.req.Method == "POST" {
		if x.req.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			postParam := x.params.(map[string]string)
			if len(postParam) > 0 {
				var paramsStr string
				var buf bytes.Buffer
				for k, v := range postParam {
					buf.WriteString(url.QueryEscape(k))
					buf.WriteByte('=')
					buf.WriteString(url.QueryEscape(v))
					buf.WriteByte('&')
				}
				paramsStr = buf.String()
				paramsStr = paramsStr[0 : len(paramsStr)-1]
				x.body(paramsStr)
			}
		} else {
			x.body(x.params)
		}

	}

	// timeout, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
	// client := http.Client{
	// 	Timeout: time.Duration(timeout) * time.Second,
	// 	// 连接池配置
	// 	Transport: &http.Transport{
	// 		MaxIdleConns: maxIdleConns,
	// 		MaxIdleConnsPerHost: maxIdleConnsPerHost,
	// 		IdleConnTimeout: idleConnTimeout * time.Second,
	// 		DialContext: (&net.Dialer{
	// 			Timeout: dialerTimeout * time.Second,
	// 			KeepAlive: dialerKeepAlive * time.Second,
	// 		}).DialContext,
	// 	},
	// }
	for i := 0; x.settings.Retries == -1 || i <= x.settings.Retries; i++ {
		resp, err = client.Do(x.req)
		if err == nil {
			break
		}
		time.Sleep(x.settings.RetryDelay)
	}
	return resp, err
}

// 组装buildUrl
func (x *XiaoeHttpRequest) buildUrl(data interface{}) {
	var paramBody string
	switch data.(type) {
	case string:
		paramBody = data.(string)
	case map[string]string:
		mapParams := data.(map[string]string)
		var buff bytes.Buffer
		for k, v := range mapParams {
			buff.WriteString(url.QueryEscape(k))
			buff.WriteByte('=')
			buff.WriteString(url.QueryEscape(v))
			buff.WriteByte('&')
		}
		paramBody = buff.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}
	if x.req.Method == "GET" && len(paramBody) > 0 {
		if strings.Contains(x.url, "?") {
			x.url += "&" + paramBody
		} else {
			x.url = x.url + "?" + paramBody
		}
		return
	}
}

// post请求时设置请求体
func (x *XiaoeHttpRequest) body(data interface{}) *XiaoeHttpRequest {
	switch t := data.(type) {
	case string:
		buff := bytes.NewBufferString(t)
		x.req.Body = ioutil.NopCloser(buff)
		x.req.ContentLength = int64(len(t))
	case []byte:
		buff := bytes.NewBuffer(t)
		x.req.Body = ioutil.NopCloser(buff)
		x.req.ContentLength = int64(len(t))
	default:
		jsonStr, _ := util.JsonEncode(data)
		x.req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonStr))
	}

	return x
}

// 获取byte类型返回
func (x *XiaoeHttpRequest) Bytes() ([]byte, error) {
	bT := time.Now()
	resp, err := x.getResponse()
	if resp != nil {
		defer resp.Body.Close()
	}
	// 请求错误
	if err != nil {
		return nil, err
	}
	// 没有内容
	if resp.Body == nil {
		return nil, nil
	}
	// 压缩类型
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(reader)
		return body, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	// Http状态码错误
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("请求[%s]Url: %s，返回错误的状态码[%d]", x.req.Method, x.req.URL, resp.StatusCode)
		if err != nil {
			msg = fmt.Sprintf("%s, Error: %s", msg, err.Error())
		} else {
			msg = fmt.Sprintf("%s\n Params: %+v\n Resp: %s", msg, x.params, string(body))
		}
		return nil, errors.New(msg)
	}
	// debug日志
	if x.debug() {
		eT := time.Since(bT)
		// string(body)
		msg := fmt.Sprintf("- 发起第三方请求【%s】：Url: %s\n Method: %s\n Params: %+v\n Resp: %s\n", eT, x.req.URL, x.req.Method, x.params, "~不打印")
		// 输出打印
		fmt.Println(msg)
		// 输出相关日志
		// 这里暂时不能打开，除非job也有日志体系
		// logging.Info(msg)
	}
	return body, err
}

// 获取是否Debug模式
func (x *XiaoeHttpRequest) debug() bool {
	return os.Getenv("RUNMODE") == "debug"
}

// 获取json格式返回
func (x *XiaoeHttpRequest) ToJSON(v interface{}) error {
	data, err := x.Bytes()
	if err != nil {
		return err
	}

	return util.JsonDecode(data, v)
}

// 获取map格式的返回
func (x *XiaoeHttpRequest) ToMap() (map[string]interface{}, error) {
	var data map[string]interface{}
	err := x.ToJSON(&data)
	return data, err
}
