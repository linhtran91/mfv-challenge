package services

import (
	"context"
	"encoding/json"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/models"
	"mfv-challenge/internal/usecases"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type transaction struct {
	transactionUC TransactionUsecase
}

type TransactionUsecase interface {
	List(ctx context.Context, userID, accountID int64, limit, offset int) ([]*models.Transaction, error)
	Create(ctx context.Context, userID int64, tran *models.Transaction) (*usecases.TransactionResponse, error)
}

func NewTransaction(transactionUC TransactionUsecase) *transaction {
	return &transaction{transactionUC: transactionUC}
}

func (s *transaction) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()
	inputs := mux.Vars(r)
	userID, err := strconv.Atoi(inputs["user_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	accountID := getValueFromUrl(r.URL.Query(), "account_id", constants.DefaultAccountID)
	limit, offset := getLimitOffset(r.URL.Query())
	response, err := s.transactionUC.List(ctx, int64(userID), int64(accountID), limit, offset)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOKResponse(w, response)
}

func (s *transaction) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()
	inputs := mux.Vars(r)
	userID, err := strconv.Atoi(inputs["user_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var current *usecases.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&current); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	tran := &models.Transaction{
		Amount:          current.Amount,
		TransactionType: current.TransactionType,
		CreatedAt:       time.Now(),
		AccountID:       current.AccountID,
	}
	response, err := s.transactionUC.Create(ctx, int64(userID), tran)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, response)
}
