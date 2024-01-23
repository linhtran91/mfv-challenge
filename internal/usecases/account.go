package usecases

import (
	"context"
	"mfv-challenge/internal/models"
)

type AccountRepository interface {
	Get(ctx context.Context, id int64) (*models.Account, error)
	GetByUserIDAndAccountID(ctx context.Context, userID, accountID int64) (*models.Account, error)
}

type account struct {
	accountRepo AccountRepository
}

func NewAccount(accountRepository AccountRepository) *account {
	return &account{accountRepo: accountRepository}
}

type Account struct {
	ID      int64   `json:"id"`
	UserID  int64   `json:"user_id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

func (u *account) Get(ctx context.Context, id int64) (*Account, error) {
	r, err := u.accountRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:      id,
		UserID:  r.UserID,
		Name:    r.Name,
		Balance: r.Balance,
	}, nil
}

func (u *account) GetByUserIDAndAccountID(ctx context.Context, userID, accountID int64) (*Account, error) {
	r, err := u.accountRepo.GetByUserIDAndAccountID(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:      accountID,
		UserID:  r.UserID,
		Name:    r.Name,
		Balance: r.Balance,
	}, nil
}
