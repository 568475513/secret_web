package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	e "abs/pkg/enums"
	"abs/pkg/logging"
)

// BaseRequest 发送http请求结构体
// 各种Http请求只能通过这个方法发送，也只能写在这个文件里面！！！【便于管理和防止代码污染】
type BaseRequest struct {
	// 请求基类
	// ...
}

// request发送请求到业务侧
func (br *BaseRequest) Request(method, url string, data interface{}) (m map[string]interface{}, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	// 这里是保底的超时设置，所有服务请求的时间都不可以超过3s
	client := &http.Client{
		Timeout: time.Duration(3 * time.Second),
		Transport: &http.Transport{
			MaxIdleConns: 32,
			MaxIdleConnsPerHost: 32,
			IdleConnTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,

		},
	}
	res, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logging.Error(err.Error())
		return
	}
	res.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(res)
	if err != nil {
		logging.Error(err.Error())
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	// 判断业务侧是否接收数据成功.
	case http.StatusOK:
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&m)
		if err == nil {
			// 可以添加多种返回的判断
			if v, ok := m["code"]; ok {
				if v.(float64) != e.SUCCESS {
					data, err := json.Marshal(m)
					if err == nil {
						logging.Error("Receive Fail, message: " + string(data) + " ,url:" + url + " ,data:" + string(jsonStr))
					}
				}
			} else {
				err = errors.New(fmt.Sprintf("Reponse not code!!! is: %+v", m))
			}
		}
		// 兼容奇怪错误
		if len(m) == 0 {
			err = errors.New(fmt.Sprintf("Reponse Nonstandard !!! is: %+v and method: %s and url: %s", m, method, url))
		}
	// 非200状态码错误情况.
	default:
		result, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(fmt.Sprintf("Send Data Fail: %s, Url: %s, Param: %+v, Res: %s", resp.Status, url, data, string(result)))
		logging.Warn(err.Error())
	}

	client.CloseIdleConnections()
	return
}

// RecordFailData 记录请求失败的数据
func (br *BaseRequest) RecordFailData(data interface{}) {
	logging.Error(fmt.Sprintf("RecordFailData!!!: %+v", data))
}

// ...
func NewBaseRequest() *BaseRequest {
	return &BaseRequest{}
}
