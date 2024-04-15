package service

import "context"

type CodeService struct {
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	return nil
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return false, nil
}
