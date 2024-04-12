package main

import (
	"basic-go/week2/webook/internal/repository"
	"basic-go/week2/webook/internal/repository/dao"
	"basic-go/week2/webook/internal/service"
	"basic-go/week2/webook/internal/web"
	"basic-go/week2/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {

	// 连接mysql + 简表
	db := initDB()

	server := initServer()
	initUser(db, server)
	server.Run(":8080")

}

func initUser(db *gorm.DB, server *gin.Engine) {
	// 从dao层依次创建
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	svc := service.NewUserService(ur)
	c := web.NewUserHandler(svc)
	c.RegisterRoutes(server)
}

func initServer() *gin.Engine {
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
		println("这是一个middleware")
	})

	loginMiddleWare := &middleware.LoginMiddlewareBuilder{}
	// 方式一：创建cookie的存储方式
	// store := cookie.NewStore([]byte("secret"))
	// 方式二：基于内存的实现，第一个参数是 authentication key 32位或64位无特殊字符；第二个参数是 encryption key
	// store := memstore.NewStore([]byte("IKD20XkWAXJus2zS7R97SH51K7XgQrLb"),
	// 	[]byte("TYJ5tKRWpIfBYWBPLMK9bGxKLAgkpXXN"))
	// 方式三：
	// 初始化一个session，命名为ssid，并以cookie存储ssid
	store, err := redis.NewStore(16, "tcp", "localhost:6379",
		"", []byte("IKD20XkWAXJus2zS7R97SH51K7XgQrLb"), []byte("IKD20XkWAXJus2zS7R97SH51K7XgQrLa"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store))
	// 检查登录状态
	server.Use(loginMiddleWare.CheckLogin())
	return server
}

func initDB() *gorm.DB {
	// gorm连接mysql
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:13316)/webook"))
	if err != nil {
		// 服务器都出错就直接panic不用return啦
		panic(err)
	}
	// 建表
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
