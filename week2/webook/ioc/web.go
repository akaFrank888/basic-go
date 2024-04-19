package ioc

import (
	"basic-go/week2/webook/internal/web"
	"basic-go/week2/webook/internal/web/middleware"
	"basic-go/week2/webook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// middleware解决跨域问题【有跨域问题时才会触发】
		cors.New(cors.Config{
			// 1. origin
			//AllowAllOrigins: true,
			AllowOrigins: []string{"https://foo.com"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost")
					return true
				}
				return strings.Contains(origin, "your_company.com") // 只有公司的域名可以跨域访问
			},
			// 2. method（不写就是默许全部方法）
			//AllowMethods:     []string{"PUT", "POST"},
			// 3. headers
			// 前端要把token放在authorization里面
			AllowHeaders: []string{"content-type", "authorization"},
			// 允许前端访问到你的后端响应中带的header【跨域问题类型】【加几个header就要在这允许几个】
			ExposeHeaders: []string{"x-jwt-token"},
			// 4. 是否允许cookie
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		func(ctx *gin.Context) { // TODO 因为是 HandlerFunc 类型的不定参数，所以可以传多个
			println("这是一个middleware")
		},
		ratelimit.NewBuilder(redisClient, time.Second, 1000).Build(),
		(&middleware.LoginJWTMiddlewareBuilder{}).CheckLoginJWT(),
	}
}
