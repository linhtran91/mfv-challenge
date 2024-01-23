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
	userUsecase UserUsecase
}

type UserUsecase interface {
	GetDetail(ctx context.Context, id int64) (*usecases.UserAccounts, error)
	ListAccount(ctx context.Context, id int64) ([]*models.Account, error)
}

func NewUser(userUsecase UserUsecase) *user {
	return &user{userUsecase: userUsecase}
}

func (s *user) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()

	inputs := mux.Vars(r)
	userID, err := strconv.Atoi(inputs["user_id"])
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response, err := s.userUsecase.GetDetail(ctx, int64(userID))
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, response)
}

func (s *user) ListAccounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()

	inputs := mux.Vars(r)
	userID, err := strconv.Atoi(inputs["user_id"])
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	accounts, err := s.userUsecase.ListAccount(ctx, int64(userID))
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError, "Internal Server Error")
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
