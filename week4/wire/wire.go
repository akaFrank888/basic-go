//go:build wireinject

package wire

import (
	"basic-go/week4/wire/repository"
	"basic-go/week4/wire/repository/dao"
	"github.com/google/wire"
)

func initUserRepository() *repository.UserRepository {
	// 参数是方法本身，方法名后不要加()，否则就是调用该方法了
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, initDB)
	return &repository.UserRepository{}
}
