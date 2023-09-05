package services

import (
	"api-service/models"
	"api-service/respositories"
	"context"
)

type AccountService interface {
	GetAccounts() ([]*models.Account, error)
	GetAccountByUsername(username string) (*models.Account, error)
	CreateAccount(account *models.Account) error
	UpdateAccount(account *models.Account) error
	DeleteAccount(account *models.Account) error
}

type accountService struct {
	ctx         context.Context
	accountRepo respositories.Account
}

func NewAccountService(ctx context.Context) *accountService {
	return &accountService{
		ctx:         ctx,
		accountRepo: respositories.NewAccountRepository(),
	}
}

func (s *accountService) GetAccounts() ([]*models.Account, error) {
	return s.accountRepo.GetAccounts(s.ctx)
}

func (s *accountService) GetAccountByUsername(username string) (*models.Account, error) {
	result, err := s.accountRepo.GetAccountByUsername(s.ctx, username)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrorAccountNotFound
	}
	return result, nil
}

func (s *accountService) CreateAccount(account *models.Account) error {
	return s.accountRepo.InsertAccount(s.ctx, account)
}

func (s *accountService) UpdateAccount(account *models.Account) error {
	return s.accountRepo.UpdateAccount(s.ctx, account)
}

func (s *accountService) DeleteAccount(account *models.Account) error {
	return s.accountRepo.DeleteAccount(s.ctx, account)
}
