package repositories

import (
	"gorm.io/gorm"
)

type transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) *transaction {
	return &transaction{db: db}
}
