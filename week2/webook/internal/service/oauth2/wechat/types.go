package wechat

import (
	"basic-go/week2/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Service interface {
	// AuthURL 拼接请求wx接口的url
	AuthURL(ctx context.Context, state string) (string, error)
	// VerifyCode 拿到wx确认后跳转回来带有code临时授权码的url
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

// 替换掉 APPID  REDIRECT_URI SCOPE STATE
const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`

// 1. 从wx返回的地址 2. 将AuthURL中的redirectURI进行URL encode（使得url中的"://"编码成"%3A%2F%2F"）
var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type service struct {
	// 一个应用中appID是不会变的
	appID     string
	appSecret string
	client    *http.Client
}

func NewService(appID string, appSecret string) Service {
	return &service{
		appID:     appID,
		client:    http.DefaultClient, // 先用一个默认的client，但并不符合依赖注入，因为目前不需要定制client
		appSecret: appSecret,
	}
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(authURLPattern, s.appID, redirectURI, state), nil
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	// 在代码里面向”https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code“发送请求
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		s.appID, s.appSecret, code)
	// 创建一个req
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	// 发送请求
	rep, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	// 将返回的rep中body中的json反序列化成结构体
	// note 注意json.NewDecoder().Decode和json.Unmarshal的区别
	var res Result
	err = json.NewDecoder(rep.Body).Decode(&res)
	if err != nil {
		// 转json为结构体出错
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		// note error.New("")的缺点是不能用占位符
		return domain.WechatInfo{}, fmt.Errorf("微信接口调用失败，错误码：%d，错误信息：%s", res.ErrCode, res.ErrMsg)
	}
	// 成功验证临时授权码code，即登录成功
	return domain.WechatInfo{
		OpenId:  res.Openid,
		UnionId: res.UnionId,
	}, nil
}

type Result struct {
	// （1）正确情况下的返回字段
	// 接口调用凭证
	AccessToken string `json:"access_token"`
	// access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`
	// 用户刷新access_token
	RefreshToken string `json:"refresh_token"`
	// 授权用户唯一标识
	Openid string `json:"openid"`
	// 用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	// 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionId string `json:"unionid"`
	// （2）错误情况下的返回字段
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
