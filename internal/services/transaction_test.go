package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/usecases"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock "mfv-challenge/mocks/services"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateTransaction(t *testing.T) {
	cases := []struct {
		name          string
		userID        string
		body          usecases.TransactionRequest
		buildStubs    func(tranUC *mock.MockTransactionUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "success",
			userID: "1",
			body: usecases.TransactionRequest{
				AccountID:       2,
				Amount:          10000,
				TransactionType: "deposit",
			},
			buildStubs: func(tranUC *mock.MockTransactionUsecase) {
				tranUC.EXPECT().Create(gomock.Any(), int64(1), gomock.Any()).Return(&usecases.TransactionResponse{
					ID:              1,
					AccountID:       2,
					Amount:          10000,
					Bank:            "ACB",
					TransactionType: "deposit",
					CreatedAt:       time.Now(),
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "bad request with unbalanced withdraw",
			userID: "1",
			body: usecases.TransactionRequest{
				AccountID:       2,
				Amount:          10000,
				TransactionType: "withdraw",
			},
			buildStubs: func(tranUC *mock.MockTransactionUsecase) {
				tranUC.EXPECT().Create(gomock.Any(), int64(1), gomock.Any()).Return(nil, constants.ErrorWithdraw)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "bad request with unsupported type",
			userID: "1",
			body: usecases.TransactionRequest{
				AccountID:       2,
				Amount:          10000,
				TransactionType: "withdraw",
			},
			buildStubs: func(tranUC *mock.MockTransactionUsecase) {
				tranUC.EXPECT().Create(gomock.Any(), int64(1), gomock.Any()).Return(nil, constants.ErrUnsupportedTransactionType)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockTransactionUsecase(ctrl)
			tc.buildStubs(uc)

			service := NewTransaction(uc)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/users/%s/transactions", tc.userID)

			data, _ := json.Marshal(tc.body)
			request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			request = mux.SetURLVars(request, map[string]string{
				"user_id": tc.userID,
			})

			service.Create(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListTransaction(t *testing.T) {
	cases := []struct {
		name          string
		userID        string
		buildStubs    func(tranUC *mock.MockTransactionUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "success",
			userID: "1",
			buildStubs: func(tranUC *mock.MockTransactionUsecase) {
				tranUC.EXPECT().List(gomock.Any(), int64(1), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*usecases.TransactionResponse{
					{
						ID:              1,
						AccountID:       2,
						Amount:          12000,
						Bank:            "",
						TransactionType: "deposit",
					},
					{
						ID:              1,
						AccountID:       2,
						Amount:          10000,
						Bank:            "",
						TransactionType: "withdraw",
					},
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "bad request",
			userID:     "a",
			buildStubs: func(tranUC *mock.MockTransactionUsecase) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "not found",
			userID: "1",
			buildStubs: func(tranUC *mock.MockTransactionUsecase) {
				tranUC.EXPECT().List(gomock.Any(), int64(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, constants.ErrorRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockTransactionUsecase(ctrl)
			tc.buildStubs(uc)

			service := NewTransaction(uc)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/users/%s/transactions", tc.userID)

			request, _ := http.NewRequest(http.MethodGet, url, nil)
			request = mux.SetURLVars(request, map[string]string{
				"user_id": tc.userID,
			})

			service.List(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
