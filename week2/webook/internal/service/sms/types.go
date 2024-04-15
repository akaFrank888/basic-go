package sms

import "context"

// Service 发送短信的抽象，为了屏蔽不同供应商的区别
type Service interface {
	// Send 给多个numbers发短信
	Send(ctx context.Context, tplId string,
		args []string, numbers ...string) error
}
