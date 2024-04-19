package tencent

// 引入sms
import (
	"basic-go/week2/webook/internal/service/sms"
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
	"os"
)

type Service struct {
	client *tencentSMS.Client
	// 腾讯云的短信SDK设计的就是string的指针
	appId    *string
	signName *string
}

func NewTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯 SMS 的 secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯 SMS 的 secret key")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return NewService(c, "1400842696", "妙影科技")
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := tencentSMS.NewSendSmsRequest()
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = &tplId
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)
	response, err := s.client.SendSms(request)
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		// 解引用
		status := *statusPtr
		if status.Code != nil && *status.Code != "Ok" {
			return fmt.Errorf("短信发送失败 code:%s mag: %s", *status.Code, status.Message)
		}
	}
	return nil
}

func NewService(client *tencentSMS.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}
