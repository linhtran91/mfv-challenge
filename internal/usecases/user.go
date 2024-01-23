package usecases

import (
	"context"
	"mfv-challenge/internal/models"
)

type UserRepository interface {
	GetDetail(ctx context.Context, id int64) ([]*models.UserAccount, error)
	ListAccount(ctx context.Context, id int64) ([]*models.Account, error)
}

type user struct {
	userRepo UserRepository
}

func NewUser(userRepository UserRepository) *user {
	return &user{userRepo: userRepository}
}

type UserAccounts struct {
	ID       int64   `json:"id"`
	Username string  `json:"name"`
	Accounts []int64 `json:"account_ids"`
}

func (u *user) GetDetail(ctx context.Context, id int64) (*UserAccounts, error) {
	accounts, err := u.userRepo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(accounts))
	for _, account := range accounts {
		ids = append(ids, account.AccountID)
	}
	return &UserAccounts{
		ID:       id,
		Username: accounts[0].Username,
		Accounts: ids,
	}, nil
}

func (u *user) ListAccount(ctx context.Context, id int64) ([]*models.Account, error) {
	return u.userRepo.ListAccount(ctx, id)
}
