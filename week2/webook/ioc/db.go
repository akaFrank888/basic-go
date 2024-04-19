package ioc

import (
	"basic-go/week2/webook/config"
	"basic-go/week2/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// gorm连接mysql
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 服务器都出错就直接panic不用return啦
		panic(err)
	}
	// 建表（有点耦合，但没优化办法）
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
