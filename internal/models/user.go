package models

import "time"

type User struct {
	ID       int64
	Username string
	Password string
}

type Account struct {
	ID      int64
	Name    string
	Bank    string
	Balance float64
	UserID  int64
}

type Transaction struct {
	ID              int64
	Amount          float64
	TransactionType string
	CreatedAt       time.Time
	AccountID       int64
}
