package service

import (
	//"io/ioutil"
	"os/exec"
)

type TimeRestApi struct {
	SdkAppId   string
	Identifier string
}

const (
	protectedKeyPath = "./storage/key/private_key"
	toolPath         = "./storage/alive_video/signature/linux-signature64"
)

// 独立模式根据Identifier生成UserSig的方法
func (t *TimeRestApi) GenerateUserSig() (string, error) {
	result, err := exec.Command(toolPath, protectedKeyPath, t.SdkAppId, t.Identifier).Output()
	if err != nil {
		return "", err
	}
	return string(result), nil
}
