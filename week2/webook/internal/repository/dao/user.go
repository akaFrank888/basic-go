package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ErrDuplicateEmail 预自定义一个错误
var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrDuplicatePhone = errors.New("手机号冲突")
	// ErrRecordNotFound gorm框架有 未找到某条数据 得错误
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	UpdateById(ctx context.Context, persistent User) error
	FindById(ctx context.Context, id int64) (User, error)
}
type GORMUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{
		db: db,
	}
}

// Insert dao层要返回自己定义的User，而不是domain.User
func (dao *GORMUserDao) Insert(ctx context.Context, u User) error {
	// 取当前毫秒数
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	// 获取err，检查是否是邮箱冲突
	err := dao.db.WithContext(ctx).Create(&u).Error
	// note 类型断言
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062 // 从控制台看到的具体Error Number
		if me.Number == duplicateErr {
			// 邮箱冲突
			// return一个特定的错误
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) UpdateById(ctx context.Context, persistent User) error {
	return dao.db.WithContext(ctx).Model(&persistent).Where("id=?", persistent.Id).Updates(map[string]any{
		"utime":    time.Now().UnixMilli(), // 更新时间
		"nickname": persistent.Nickname,
		"birthday": persistent.Birthday,
		"resume":   persistent.Resume,
	}).Error
}

func (dao *GORMUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
	return u, err
}

// User 相当于PO，即属性与表字段一一对应
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 代表可为Null（因为用户用phone注册的话，会没有email）
	Email    sql.NullString `gorm:"unique"`
	Password string
	Nickname string `gorm:"type=varchar(20)"`
	Birthday int64
	Resume   string `gorm:"type=varchar(200)"`

	Phone sql.NullString `gorm:"unique"`

	// 创建时间  避免时区问题，一律用 UTC 0 的毫秒数【若要转成符合中国的时区，要么让前端处理，要么在web层给前端的时候转成UTC 8 的时区】
	Ctime int64
	// 更新时间
	Utime int64
}
