package services

import (
	"context"
	"encoding/json"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/models"
	"mfv-challenge/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	userUsecase UserUsecase
	encoder     TokenEncoder
}

type UserUsecase interface {
	GetDetail(ctx context.Context, id int64) (*usecases.UserAccounts, error)
	ListAccount(ctx context.Context, id int64) ([]*models.Account, error)
	ComparePassword(ctx context.Context, info usecases.LoginInfo) (bool, int64, error)
}

type TokenEncoder interface {
	Encode(userID int64) (string, error)
}

func NewUser(userUsecase UserUsecase, encoder TokenEncoder) *user {
	return &user{
		userUsecase: userUsecase,
		encoder:     encoder,
	}
}

func (s *user) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()

	inputs := mux.Vars(r)
	userID, err := strconv.Atoi(inputs["user_id"])
	if err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest, "Bad Request")
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
		writeErrorResponse(w, err, http.StatusBadRequest, "Bad Request")
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

func (s *user) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()
	var user usecases.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest, "Bad request")
		return
	}

	isValid, id, err := s.userUsecase.ComparePassword(ctx, user)
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if !isValid {
		writeErrorResponse(w, err, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := s.encoder.Encode(id)
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, map[string]interface{}{
		"token": token,
	})
}
