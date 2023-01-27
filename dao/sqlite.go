package dao

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSQLiet() {
	db, err = gorm.Open(sqlite.Open("potatoes.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Couldn't connect to SQLite:", err)
		os.Exit(1)
	}
	// Auto Migrate
	db.AutoMigrate(&Potato{})
}
