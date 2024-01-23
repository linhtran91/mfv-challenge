package services

import (
	"context"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type account struct {
	accountUC AccountUsecase
}

type AccountUsecase interface {
	Get(ctx context.Context, id int64) (*usecases.Account, error)
}

func NewAccount(accountUC AccountUsecase) *account {
	return &account{accountUC: accountUC}
}

func (s *account) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()

	inputs := mux.Vars(r)
	userID, err := strconv.Atoi(inputs["account_id"])
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response, err := s.accountUC.Get(ctx, int64(userID))
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError, err.Error())
		return
	}
	writeOKResponse(w, response)
}
