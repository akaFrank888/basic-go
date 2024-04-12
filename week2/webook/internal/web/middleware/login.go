package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 为了将go中time类型当作session的val存入redis，需要将time类型注册一下
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 登录校验
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 这两个页面不需要校验是否登录
			return
		}

		// 获取session，登录校验
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 未登录
			ctx.AbortWithStatus(http.StatusUnauthorized) // 401
			return
		}

		// 每一分钟刷新一次登录状态，
		now := time.Now()
		const updateTimeKey = "update_time"
		// 取出上次的更新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		// 两种情况都要更新上次的更新时间为现在的时间：“val==nil：第一次登录进来”和“距上次更新已过10s”
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			sess.Set(updateTimeKey, now)
			// TODO go的session机制问题：set后也要更新其他的key，不然会丢失
			sess.Set("userId", userId)
			// session进行set后需要save
			err := sess.Save()
			if err != nil {
				// 打日志
				fmt.Println(err)
			}

		}
	}
}
