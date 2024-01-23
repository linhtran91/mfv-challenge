package services

import (
	"net/http"
)

type account struct {
	accountRepo AccountRepository
}

type AccountRepository interface{}

func NewAccount(accountRepository AccountRepository) *account {
	return &account{accountRepo: accountRepository}
}

func (s *account) Get(w http.ResponseWriter, r *http.Request) {}
