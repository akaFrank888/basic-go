//go:build wireinject

package main

import (
	"basic-go/week2/webook/internal/repository"
	"basic-go/week2/webook/internal/repository/cache"
	"basic-go/week2/webook/internal/repository/dao"
	"basic-go/week2/webook/internal/service"
	"basic-go/week2/webook/internal/web"
	"basic-go/week2/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		// dao和cache
		dao.NewUserDao, cache.NewUserCache, cache.NewCodeCache,
		// repository
		repository.NewCachedUserRepository, repository.NewCodeRepository,
		// service
		ioc.InitSMSService, service.NewUserService, service.NewCodeService,
		// handler
		web.NewUserHandler,
		// gin.Engine部分
		ioc.InitGinMiddlewares, ioc.InitWebServer,
	)
	return gin.Default()
}
