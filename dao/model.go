package dao

import (
	"time"

	"gorm.io/gorm"
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
