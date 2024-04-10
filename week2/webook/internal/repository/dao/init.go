package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {

	// 自动建表会有字段不符合公司标准的风险
	return db.AutoMigrate(&User{})
}
