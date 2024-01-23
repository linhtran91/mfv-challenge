package repositories

import (
	"gorm.io/gorm"
)

type account struct {
	db *gorm.DB
}

func NewAccount(db *gorm.DB) *account {
	return &account{db: db}
}
