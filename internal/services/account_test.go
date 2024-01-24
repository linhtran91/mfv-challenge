package services

import (
	"fmt"
	"mfv-challenge/internal/usecases"
	mock "mfv-challenge/mocks/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	cases := []struct {
		name          string
		accountID     string
		buildStubs    func(uc *mock.MockAccountUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "success",
			accountID: "1",
			buildStubs: func(uc *mock.MockAccountUsecase) {
				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(&usecases.Account{
					ID:      1,
					UserID:  1,
					Name:    "alice",
					Balance: 10000,
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "bad request",
			accountID:  "a",
			buildStubs: func(tranUC *mock.MockAccountUsecase) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockAccountUsecase(ctrl)
			tc.buildStubs(uc)

			service := NewAccount(uc)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/accounts/%s", tc.accountID)

			request, _ := http.NewRequest(http.MethodGet, url, nil)
			request = mux.SetURLVars(request, map[string]string{
				"account_id": tc.accountID,
			})

			service.Get(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
