package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/models"
	"mfv-challenge/internal/token"
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

func TestLogin(t *testing.T) {
	secretKey := "9876hjds"
	duration := 5 * time.Minute
	cases := []struct {
		name          string
		body          usecases.LoginInfo
		buildStubs    func(uc *mock.MockUserUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			body: usecases.LoginInfo{
				Username: "tester",
				Password: "321ndsjkqe",
			},
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Return(true, int64(1), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "unauthorized",
			body: usecases.LoginInfo{
				Username: "tester",
				Password: "321ndsjkqe23",
			},
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Return(false, int64(0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "internal error",
			body: usecases.LoginInfo{
				Username: "tester",
				Password: "321ndsjkqe23",
			},
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Return(false, int64(0), errors.New("internal error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockUserUsecase(ctrl)
			tokenBuilder := token.NewjwtHMACBuilder(secretKey, duration)
			tc.buildStubs(uc)

			service := NewUser(uc, tokenBuilder)
			recorder := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			request, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(data))

			service.Login(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccounts(t *testing.T) {
	cases := []struct {
		name          string
		userID        string
		buildStubs    func(uc *mock.MockUserUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "success",
			userID: "1",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().ListAccount(gomock.Any(), int64(1)).Return([]*models.Account{
					{
						ID:      1,
						UserID:  1,
						Name:    "account 1",
						Balance: 10000,
					},
					{
						ID:      2,
						UserID:  1,
						Name:    "account 2",
						Balance: 10000,
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
			buildStubs: func(uc *mock.MockUserUsecase) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "not found",
			userID: "1",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().ListAccount(gomock.Any(), int64(1)).Return(nil, constants.ErrorRecordNotFound)
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

			uc := mock.NewMockUserUsecase(ctrl)
			tc.buildStubs(uc)

			service := NewUser(uc, mock.NewMockTokenEncoder(ctrl))
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/users/%s/transactions", tc.userID)

			request, _ := http.NewRequest(http.MethodGet, url, nil)
			request = mux.SetURLVars(request, map[string]string{
				"user_id": tc.userID,
			})

			service.ListAccounts(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetUser(t *testing.T) {
	cases := []struct {
		name          string
		userID        string
		buildStubs    func(uc *mock.MockUserUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "success",
			userID: "1",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().GetDetail(gomock.Any(), int64(1)).Return(&usecases.UserAccounts{
					ID:       1,
					Username: "alice",
					Accounts: []int64{1, 2, 3},
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "bad request",
			userID:     "a1",
			buildStubs: func(uc *mock.MockUserUsecase) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "not found",
			userID: "1",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().GetDetail(gomock.Any(), int64(1)).Return(nil, constants.ErrorRecordNotFound)
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

			uc := mock.NewMockUserUsecase(ctrl)
			tc.buildStubs(uc)

			service := NewUser(uc, mock.NewMockTokenEncoder(ctrl))
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/users/%s/transactions", tc.userID)

			request, _ := http.NewRequest(http.MethodGet, url, nil)
			request = mux.SetURLVars(request, map[string]string{
				"user_id": tc.userID,
			})

			service.Get(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
