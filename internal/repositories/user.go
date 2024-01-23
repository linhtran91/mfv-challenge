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

func (r *user) GetDetail(ctx context.Context, id int64) (string, []int64, error) {

	return "", nil, nil
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

// func (r *loan) Create(ctx context.Context, loan *models.Loan, repayments []*models.Repayment) (int64, error) {
// 	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Create(&loan).Error; err != nil {
// 			return err
// 		}

// 		for i := range repayments {
// 			repayments[i].LoanID = loan.ID
// 		}

// 		if err := tx.CreateInBatches(repayments, 50).Error; err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return -1, err
// 	}
// 	return loan.ID, nil
// }

// func (r *loan) Approve(ctx context.Context, loanID int64, at time.Time) error {
// 	if err := r.db.WithContext(ctx).
// 		Model(&models.Loan{}).
// 		Where("id = ?", loanID).
// 		Updates(map[string]interface{}{
// 			"status":     constants.APPROVED,
// 			"updated_at": at,
// 		}).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *loan) View(ctx context.Context, customerID int64, limit, offset int) ([]*models.Loan, error) {
// 	var result []*models.Loan
// 	err := r.db.WithContext(ctx).
// 		Model(&models.Loan{}).
// 		Where("customer_id = ?", customerID).
// 		Order("scheduled_date desc").
// 		Limit(limit).
// 		Offset(offset).
// 		Find(&result).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, constants.ErrorRecordNotFound
// 		}
// 		return nil, err
// 	}
// 	return result, nil
// }

// func (r *loan) UpdateStatus(ctx context.Context, loanID int64, at time.Time) error {
// 	if err := r.db.WithContext(ctx).Model(&models.Loan{}).
// 		Where("id = ?", loanID).
// 		Updates(map[string]interface{}{
// 			"status":     constants.PAID,
// 			"updated_at": at,
// 		}).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
