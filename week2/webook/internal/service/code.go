package service

import (
	"basic-go/week2/webook/internal/repository"
	"basic-go/week2/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

//	type CodeService interface {
//		Send(ctx context.Context, biz, phone string) error
//		Verify(ctx context.Context,
//			biz, phone, inputCode string) (bool, error)
//	}
type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		sms:  smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 你在这儿，是不是要开始发送验证码了？
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *CodeService) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
