package repositories

import (
	"context"
	"mfv-challenge/internal/models"

	"gorm.io/gorm"
)

type account struct {
	db *gorm.DB
}

func NewAccount(db *gorm.DB) *account {
	return &account{db: db}
}

func (r *account) Get(ctx context.Context, id int64) (*models.Account, error) {
	var result *models.Account
	if err := r.db.WithContext(ctx).
		Model(&models.Account{}).
		Where(`id = ?`, id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *account) GetByUserIDAndAccountID(ctx context.Context, userID, accountID int64) (*models.Account, error) {
	var result *models.Account
	if err := r.db.WithContext(ctx).
		Model(&models.Account{}).
		Joins(`JOIN users ON users.id = accounts.user_id`).
		Where(`accounts.id = ?`, accountID).
		Where(`users.id = ?`, userID).
		First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
