package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {

	initViperV1()

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

func initViper() {
	// 配置文件名称和类型
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	// 当前工作目录（Working Directory）的子目录是config
	viper.AddConfigPath("config")
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// 读取配置文件内容
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperV1() {
	// 配置文件类型
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/dev.yaml")
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// 读取配置文件内容
	val := viper.Get("test.key")
	log.Println(val)
}
