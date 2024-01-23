package services

import (
	"net/http"
)

type transaction struct {
	transactionRepo TransactionRepository
}

type TransactionRepository interface {
}

func NewTransaction(transactionRepo TransactionRepository) *transaction {
	return &transaction{transactionRepo: transactionRepo}
}

func (s *transaction) List(w http.ResponseWriter, r *http.Request) {}

func (s *transaction) Create(w http.ResponseWriter, r *http.Request) {}
