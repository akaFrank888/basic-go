package middleware

import (
	ijwt "basic-go/week2/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(hdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: hdl,
	}
}

func (m *LoginJWTMiddlewareBuilder) CheckLoginJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 登录校验
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" {
			// 不需要校验是否登录
			return
		}
		tokenStr := m.ExtractToken(ctx)
		// 解析jwt
		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})
		if err != nil {
			// token不对
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// token解析了，但不合法或者过期了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//if uc.UserAgent != ctx.GetHeader("User-Agent") {
		//	// 后期讲到监控告警的时候，这个地方要埋点
		//	// 能够进来这个分支的，大概率是攻击者
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		//// NumericDate类型（组合了time.Time）
		//expireTime := uc.ExpiresAt
		//// 刷新规则：每过10s就刷新一次【登录成功时在login设置的是1min的ExpireTime】
		//// 取上一次的更新时间，相减，若<50s就更新
		//if expireTime.Sub(time.Now()) < time.Second*50 {
		//	uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err := token.SignedString(web.JWTKey)
		//	ctx.Header("x-jwt-token", tokenStr)
		//	if err != nil {
		//		// 这边不要中断，因为仅仅是过期时间没有刷新，但是用户是登录了的
		//		log.Println(err)
		//	}
		//}

		// 在redis中查看ssid
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// redis有问题或者ssid存在（表明用户已退出） ==> redis有问题或者 token无效
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user", uc)
	}
}
