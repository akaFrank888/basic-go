package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 登录校验
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 这两个页面不需要校验是否登录
			return
		}

		// 获取session，登录校验
		sess := sessions.Default(ctx)
		if sess.Get("userId") == nil {
			// 未登录
			ctx.AbortWithStatus(http.StatusUnauthorized) // 401
			return
		}
	}
}
