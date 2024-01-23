package services

import (
	"context"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/models"
	"mfv-challenge/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	userRepo UserRepository
}

type UserRepository interface {
	GetDetail(ctx context.Context, id int64) (string, []int64, error)
	ListAccount(ctx context.Context, id int64) ([]*models.Account, error)
}

func NewUser(userRepository UserRepository) *user {
	return &user{userRepo: userRepository}
}

func (s *user) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()

	inputs := mux.Vars(r)
	userID, err := strconv.Atoi(inputs["user_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	name, transactions, err := s.userRepo.GetDetail(ctx, int64(userID))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, usecases.UserAccounts{
		ID:       int64(userID),
		Username: name,
		Accounts: transactions,
	})
}

func (s *user) ListAccounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()

	inputs := mux.Vars(r)
	userID, err := strconv.Atoi(inputs["user_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	accounts, err := s.userRepo.ListAccount(ctx, int64(userID))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	result := make([]*usecases.Account, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, &usecases.Account{
			ID:      account.ID,
			UserID:  int64(userID),
			Name:    account.Name,
			Balance: account.Balance,
		})
	}
	writeOKResponse(w, result)
}
