package wire

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("***"))
	if err != nil {
		panic(err)
	}
	return db
}
