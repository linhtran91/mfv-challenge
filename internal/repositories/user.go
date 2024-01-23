package repositories

import (
	"context"
	"mfv-challenge/internal/constants"
	"mfv-challenge/internal/models"

	"gorm.io/gorm"
)

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *user {
	return &user{db: db}
}

func (r *user) GetCredential(ctx context.Context, username string) (*models.User, error) {
	var user *models.User
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constants.ErrorRecordNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *user) GetDetail(ctx context.Context, id int64) ([]*models.UserAccount, error) {
	var result []*models.UserAccount
	if err := r.db.WithContext(ctx).Table(`users`).
		Select(`users.id, users.username, accounts.id as account_id`).
		Joins(`JOIN accounts ON users.id = accounts.user_id`).
		Where(`users.id = ?`, id).
		Order(`accounts.id asc`).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *user) ListAccount(ctx context.Context, id int64) ([]*models.Account, error) {
	var result []*models.Account
	if err := r.db.WithContext(ctx).
		Model(&models.Account{}).
		Joins(`JOIN users ON users.id = accounts.user_id`).
		Where(`users.id = ?`, id).
		Order(`accounts.id asc`).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
