package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//
	//// 连接mysql + 建表
	//db := initDB()
	//
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server := initServer()
	//codeSvc := initCodeSvc(redisClient)
	//initUser(db, server, redisClient, codeSvc)

	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello go")
	})
	server.Run(":8080")

}
