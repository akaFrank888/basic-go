package service

import (
	"basic-go/week2/webook/internal/domain"
	"basic-go/week2/webook/internal/repository"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ErrDuplicateEmail 小技巧：如果repo层返回了这个err，则web层可直接从service层调用来进行判定
var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号或密码错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 对用户密码的加密设计在service层
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

// Login 因为session，所以要返回一个domain.User
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
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

func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *UserService) FindById(ctx *gin.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}
