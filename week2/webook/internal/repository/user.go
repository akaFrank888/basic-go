package repository

import (
	"basic-go/week2/webook/internal/domain"
	"basic-go/week2/webook/internal/repository/cache"
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
	dao   *dao.UserDao
	cache *cache.UserCache
}

//func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
//	return repo.dao.Insert(ctx, dao.User{
//		Email:    u.Email,
//		Password: u.Password,
//	})
//}

func NewUserRepository(d *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   d,
		cache: c,
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

func (repo *UserRepository) FindById(ctx *gin.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	if err == nil {
		// 从缓存中查到了
		return du, nil
	}

	// TODO err不为nil有多种可能：
	// 1） 缓存中没有key，但redis正常
	// 2） 访问redis有问题。可能是连不上网，也可能redis本身崩了
	u, err := repo.dao.FindById(ctx, uid)
	// 将du存入缓存
	du = toDomain(u)
	_ = repo.cache.Set(ctx, du)
	// TODO 可以不接收err，因为这次没存进缓存，下次直接查数据库就行了。而且接受了err，也只说明连接redis的网络和本身有问题，无法解决。
	return du, nil
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
