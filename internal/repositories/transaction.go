package repositories

import (
	"context"
	"mfv-challenge/internal/models"

	"gorm.io/gorm"
)

type transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) *transaction {
	return &transaction{db: db}
}

func (r *transaction) Create(ctx context.Context, account *models.Account, tran *models.Transaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tran).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Account{}).
			Where(`id = ?`, account.ID).
			Updates(map[string]interface{}{
				"balance": account.Balance,
			}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *transaction) List(ctx context.Context, userID, accountID int64, limit, offset int) ([]*models.Transaction, error) {
	var result []*models.Transaction
	query := r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Joins(`JOIN accounts ON transactions.account_id = accounts.id`).
		Joins(`JOIN users ON accounts.user_id = users.id`).
		Where(`users.id = ?`, userID)
	if accountID > 0 {
		query = query.Where(`accounts.id = ?`, accountID)
	}
	if err := query.
		Limit(limit).
		Offset(offset).
		Order(`transactions.id desc`).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
