package ioc

import (
	"basic-go/week2/webook/config"
	"basic-go/week2/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	// note viper：用一个内部结构体接收配置
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var c Config
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic("viper初始化db失败")
	}
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
