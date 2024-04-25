package service

import (
	"basic-go/week2/webook/internal/domain"
	"basic-go/week2/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// ErrDuplicateEmail 小技巧：如果repo层返回了这个err，则web层可直接从service层调用来进行判定
var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号或密码错误")
)

type UserService interface {
	Login(ctx context.Context, email string, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// note 对用户密码的加密设计在 service 层
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	u.Password = string(hash) // []byte ==> string
	return svc.repo.Create(ctx, u)
}

// Login 因为session，所以要返回一个domain.User
func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 对比密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil

}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *userService) FindById(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// 两种情况
		// 1. err != nil ==> 系统错误
		// 2. err == nil ==> u可用
		return u, err
	}
	// Find失败就Create
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != repository.ErrDuplicatePhone {
		// 系统错误
		return domain.User{}, err
	}
	// 两种情况
	// 1. err == ErrDuplicatePhone ==> 手机号冲突
	// 2. err == nil ==> 创建成功
	// TODO 主从延迟 ==>插入进的是主库，查询查的是从库，所以可能刚插进去就查的话查不到，因为主从库还没同步完成【解决方式是强制查主库，但还没做】
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenId)
	if err != repository.ErrUserNotFound {
		// 两种情况
		// 1. err != nil ==> 系统错误
		// 2. err == nil ==> u可用
		return u, err
	}
	// Find失败就Create
	err = svc.repo.Create(ctx, domain.User{
		WechatInfo: info,
	})
	if err != nil && err != repository.ErrDuplicatePhone {
		// 系统错误
		return domain.User{}, err
	}
	// 两种情况
	// 1. err == ErrDuplicatePhone ==> 手机号冲突
	// 2. err == nil ==> 创建成功
	// TODO 主从延迟 ==>插入进的是主库，查询查的是从库，所以可能刚插进去就查的话查不到，因为主从库还没同步完成【解决方式是强制查主库，但还没做】
	return svc.repo.FindByWechat(ctx, info.OpenId)
}
