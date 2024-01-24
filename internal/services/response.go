package services

import (
	"encoding/json"
	"mfv-challenge/internal/constants"
	"net/http"
)

// Writes the response as a standard JSON response with StatusOK
func writeOKResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		writeErrorResponse(w, nil, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, err error, errorCode int, errorMsg string) {
	switch err {
	case constants.ErrorRecordNotFound:
		errorCode = http.StatusNotFound
		errorMsg = "Not Found"
	case constants.ErrorWithdraw:
		errorCode = http.StatusBadRequest
		errorMsg = "Not balanced amount"
	case constants.ErrUnsupportedTransactionType:
		errorCode = http.StatusBadRequest
		errorMsg = err.Error()
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Message: errorMsg}})
}

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
