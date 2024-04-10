package repository

import (
	"basic-go/week2/webook/internal/domain"
	"basic-go/week2/webook/internal/repository/dao"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

// ErrDuplicateEmail 小技巧：如果dao层返回了这个err，则service层可直接从repo层调用来进行判定
var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	// 得重新命名为 User 相关的，因为Service在通过repo层调用时是在具体业务中的（如User业务，而不能用Record）
	ErrUserNotFound = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

//func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
//	return repo.dao.Insert(ctx, dao.User{
//		Email:    u.Email,
//		Password: u.Password,
//	})
//}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, toPersistent(u))
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return toDomain(u), nil
}

func (repo *UserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, toPersistent(user))
}

func (repo *UserRepository) FindById(ctx *gin.Context, id int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, id)
	return toDomain(u), err
}

// 私有方法（首字母小写）
func toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		// UTC 0的毫秒 -> time
		Birthday: time.UnixMilli(u.Birthday),
		Resume:   u.Resume,
		// UTC 0的毫秒 -> time
		Ctime: time.UnixMilli(u.Ctime),
	}
}

func toPersistent(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday.UnixMilli(),
		Resume:   u.Resume,
	}
}
