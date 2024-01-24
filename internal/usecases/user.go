package usecases

import (
	"context"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/models"
)

type UserRepository interface {
	GetDetail(ctx context.Context, id int64) ([]*models.UserAccount, error)
	ListAccount(ctx context.Context, id int64) ([]*models.Account, error)
	GetCredential(ctx context.Context, username string) (*models.User, error)
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

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *user) ComparePassword(ctx context.Context, info LoginInfo) (bool, int64, error) {
	us, err := u.userRepo.GetCredential(ctx, info.Username)
	if err != nil {
		if err == constants.ErrorRecordNotFound {
			return false, 0, nil
		}
		return false, 0, err
	}

	return us.Password == info.Password, us.ID, nil
}
