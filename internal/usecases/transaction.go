package usecases

import (
	"context"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/models"
	"time"
)

type transaction struct {
	transactionRepo TransactionRepository
	accountRepo     AccountRepository
}

type TransactionRepository interface {
	List(ctx context.Context, userID, accountID int64, limit, offset int) ([]*models.Transaction, error)
	Create(ctx context.Context, account *models.Account, tran *models.Transaction) error
}

func NewTransaction(transactionRepo TransactionRepository, accountRepo AccountRepository) *transaction {
	return &transaction{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (u *transaction) List(ctx context.Context, userID, accountID int64, limit, offset int) ([]*TransactionResponse, error) {
	transactions, err := u.transactionRepo.List(ctx, userID, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	result := make([]*TransactionResponse, 0, len(transactions))
	for _, tran := range transactions {
		result = append(result, &TransactionResponse{
			ID:              tran.ID,
			AccountID:       tran.AccountID,
			Amount:          tran.Amount,
			TransactionType: tran.TransactionType,
			CreatedAt:       tran.CreatedAt,
		})
	}
	return result, nil
}

type TransactionRequest struct {
	AccountID       int64   `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

type TransactionResponse struct {
	ID              int64     `json:"id"`
	AccountID       int64     `json:"account_id"`
	Amount          float64   `json:"amount,omitempty"`
	Bank            string    `json:"bank,omitempty"`
	TransactionType string    `json:"transaction_type,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

func (u *transaction) Create(ctx context.Context, userID int64, tran *models.Transaction) (*TransactionResponse, error) {
	// Deposit (money coming into your bank account)
	// Withdrawal (money going out of your bank account).
	account, err := u.accountRepo.GetByUserIDAndAccountID(ctx, userID, tran.AccountID)
	if err != nil {
		return nil, err
	}

	switch tran.TransactionType {
	case constants.Deposit:
		account.Balance += tran.Amount
	case constants.Withdraw:
		account.Balance -= tran.Amount
	default:
		return nil, constants.ErrUnsupportedTransactionType
	}

	if account.Balance <= 0 {
		return nil, constants.ErrorWithdraw
	}
	if err := u.transactionRepo.Create(ctx, account, tran); err != nil {
		return nil, err
	}
	return &TransactionResponse{
		ID:              tran.ID,
		AccountID:       tran.AccountID,
		Amount:          tran.Amount,
		Bank:            account.Bank,
		TransactionType: tran.TransactionType,
		CreatedAt:       tran.CreatedAt,
	}, nil
}
