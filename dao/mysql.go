package dao

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func InitMySQL() {
	dsn := "username:password@tcp(127.0.0.1:3306)/potatoes?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Couldn't connect to MySQL:", err)
		os.Exit(1)
	}

	// Auto Migrate
	db.AutoMigrate(&Potato{})

	// NOTE: For Test Purpose
	db.Where("id >=0").Unscoped().Delete(&Potato{})
	potatoes := []Potato{{Entry: "Buy carrots", Checked: true}, {Entry: "Buy chocolates"}, {Entry: "Buy milk", Checked: true}}
	for _, potato := range potatoes {
		db.Create(&potato)
	}
}
