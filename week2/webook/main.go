package main

import (
	"basic-go/week2/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func main() {

	server := gin.Default()

	// middleware解决跨域问题【有跨域问题时才会触发】
	server.Use(cors.New(cors.Config{
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
		AllowHeaders: []string{"content-type"},
		//ExposeHeaders: []string{"Content-Length"},
		// 4. 是否允许cookie
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}), func(ctx *gin.Context) { // TODO 因为是 HandlerFunc 类型的不定参数，所以可以传多个
		println("第一个middleware")
	}, func(ctx *gin.Context) {
		println("第二个middleware")
	})

	c := web.NewUserHandler()
	c.RegisterRoutes(server)

	server.Run(":8080")

}
