package dao

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

type PotatoType uint

const (
	NORMAL PotatoType = 0
	DAILY  PotatoType = 1
)

type Potato struct {
	ID      uint
	Entry   string
	Type    PotatoType
	Checked bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

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

func LoadPotatoes() ([]Potato, error) {
	var potatoes []Potato

	if err := db.Order("type desc").Order("id").Find(&potatoes).Error; err != nil {
		return nil, err
	}

	return potatoes, nil
}

func AddEntry(entry string, potatoType PotatoType) error {
	potato := &Potato{Entry: entry, Type: potatoType}

	if err := db.Create(potato).Error; err != nil {
		return err
	}

	return nil
}

func ToggleCheck(id uint) {
	var potato Potato

	if err := db.First(&potato, id).Error; err != nil {
		log.Println(err)
		return
	}

	if err := db.Model(&potato).Update("checked", !potato.Checked).Error; err != nil {
		log.Println(err)
	}
}

func DeleteEntry(id uint) {
	if err := db.Delete(&Potato{}, id).Error; err != nil {
		log.Println(err)
	}
}
