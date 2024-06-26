package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	// YYYY-MM-DD
	Birthday time.Time
	Resume   string

	Phone string

	Ctime time.Time

	// 组合一下wechatInfo
	WechatInfo WechatInfo
}
