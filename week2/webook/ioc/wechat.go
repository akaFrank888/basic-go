package ioc

import (
	"basic-go/week2/webook/internal/service/oauth2/wechat"
	"os"
)

func InitWechatService() wechat.Service {
	// 从环境变量中取出 appid （因为appid是敏感信息，所以放在环境变量或配置文件中）
	// os.getenv("WECHAT_APP_ID")是假设一定有这个环境变量

	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		// panic("找不到环境变量WECHAT_APP_ID")
	}

	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		// panic("找不到环境变量WECHAT_APP_SECRET")
	}
	return wechat.NewService(appID, appSecret)
}
