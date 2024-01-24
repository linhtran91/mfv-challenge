package constants

import (
	"errors"
	"time"
)

var (
	ErrorRecordNotFound           = errors.New("record not found")
	ErrorWithdraw                 = errors.New("balance is lesser than withdraw amount")
	ErrUnauthorized               = errors.New("unauthorized")
	ErrUnsupportedTransactionType = errors.New("unsupported transaction type")
)

const LengthOfID = 16

const DefaultPage = 1
const DefaultSize = 10
const DefaultAccountID = -1
const MaximumSize = 1e5

const DefaultTimeout = 15 * time.Second

const AuthorizationHeader = "Authorization"
const AuthorizationKey = "Bearer "

const Deposit = "deposit"
const Withdraw = "withdraw"
