package ioc

import (
	"basic-go/week2/webook/internal/service/sms"
	"basic-go/week2/webook/internal/service/sms/localsms"
)

func InitSMSService() sms.Service {
	return localsms.NewService()

	// 或者是腾讯云的短信
	// return tencent.NewTencentSMSService()
}
