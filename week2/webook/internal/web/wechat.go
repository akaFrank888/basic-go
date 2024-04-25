package web

import (
	"basic-go/week2/webook/internal/service"
	"basic-go/week2/webook/internal/service/oauth2/wechat"
	ijwt "basic-go/week2/webook/internal/web/jwt"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	// 结构体就不用通过注入来构建，指针需要
	ijwt.Handler
	jWTKey          []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, hdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		jWTKey:          []byte("IKD20XkWAXJus2zS7R97SH51K7XgQrLB"),
		stateCookieName: "jwt-state",
		Handler:         hdl,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	// 跳到wx的url
	g.GET("/authurl", h.OAuth2URL)
	// 处理wx跳转回来的请求
	g.Any("/callback", h.Callback)

}

func (h *OAuth2WechatHandler) OAuth2URL(ctx *gin.Context) {
	// 该state要放到jwt中
	state := uuid.New()

	val, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造跳转wx登录的url失败",
		})
		return
	}
	err = h.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "服务器异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})

}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 校验state，防止csrf攻击
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "非法请求",
		})
	}

	code := ctx.Query("code")
	// state := ctx.Query("state")
	wechatInfo, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "授权码有误",
		})
		return
	}

	// 临时授权码code校验成功，即登录成功
	u, err := h.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 登录成功后先设置refresh-token
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误（JWT的refresh-token）")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	return
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.jWTKey)
	if err != nil {
		return err
	}
	ctx.SetCookie(h.stateCookieName, tokenStr, 600, "/oauth2/wechat/callback", "", false, true)
	return nil
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	cookie, err := ctx.Cookie(h.stateCookieName)
	if err != nil {
		return fmt.Errorf("%w, 无法获得 cookie ", err)
	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(cookie, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.jWTKey, nil
	})
	if err != nil {
		return fmt.Errorf("%w, 无法获得 cookie ", err)
	}
	if state != sc.State {
		return errors.New("state 被篡改了")
	}
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
