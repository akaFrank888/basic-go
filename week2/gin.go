package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Gin的基础知识
func main() {

	// 创建一个逻辑上的web服务器server 【golang中一个web服务器被抽象成为Engine】
	server := gin.Default()
	// Engine注册路由和接入middleware
	server.GET("/user", func(ctx *gin.Context) {
		// context 处理请求
		ctx.String(http.StatusOK, "hello world")
	})
	// 除上面的静态路由，Gin路由还有参数路由和通配符路由
	// 参数路由
	server.GET("/user/:name", func(ctx *gin.Context) {
		// 调用ctx的Param方法获取参数
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "这是传过来的name：%s", name)
	})
	// 通配符路由
	server.GET("/user/*.html", func(ctx *gin.Context) {
		path := ctx.Param(".html")
		ctx.String(http.StatusOK, "通配符路由上匹配的值：%s", path)
	})

	/*
		// 一个go进程可以开多个Engine，如下建立协程
		go func() {
			server_another := gin.Default()
			server_another.Run(":8081")
		}()
	*/

	// server监听端口
	server.Run(":8080") // 别忘记冒号

}
