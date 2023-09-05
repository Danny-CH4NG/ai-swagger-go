package respositories

import (
	"api-service/models"
	"context"
)

type Account interface {
	GetAccounts(ctx context.Context) ([]*models.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (*models.Account, error)
	InsertAccount(ctx context.Context, account *models.Account) error
	UpdateAccount(ctx context.Context, account *models.Account) error
	DeleteAccount(ctx context.Context, account *models.Account) error
}

type AccountRepository struct {
	dbRepo

	Account
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (r *AccountRepository) GetAccounts(ctx context.Context) ([]*models.Account, error) {
	var accounts []*models.Account
	db := r.getDB(ctx)
	if err := db.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) GetAccountByUsername(ctx context.Context, username string) (*models.Account, error) {
	var account models.Account
	db := r.getDB(ctx)
	if err := db.Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) InsertAccount(ctx context.Context, account *models.Account) error {
	db := r.getDB(ctx)
	return db.Create(account).Error
}

func (r *AccountRepository) UpdateAccount(ctx context.Context, account *models.Account) error {
	db := r.getDB(ctx)
	return db.Save(account).Error
}

func (r *AccountRepository) DeleteAccount(ctx context.Context, account *models.Account) error {
	db := r.getDB(ctx)
	return db.Delete(account).Error
}
