package core

import (
	"log"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type AccountService interface {
	Create(account *models.Account) error
	Get(id int) (*models.Account, error)
	Update(account *models.Account) error
	Delete(id int) error
	List() ([]*models.Account, error)
}

type accountService struct {
	repo   storage.Repository
	logger *log.Logger
}

func NewAccountService(core *Core) AccountService {
	return &accountService{
		repo:   core.Repository,
		logger: core.Logger,
	}
}

func (s *accountService) Create(account *models.Account) error {
	if err := s.repo.CreateAccount(account); err != nil {
		s.logger.Printf("Failed to create account: %v", err)
		return err
	}
	return nil
}

func (s *accountService) Get(id int) (*models.Account, error) {
	account, err := s.repo.GetAccount(id)
	if err != nil {
		s.logger.Printf("Failed to get account: %v", err)
		return nil, err
	}
	return account, nil
}

func (s *accountService) Update(account *models.Account) error {
	if err := s.repo.UpdateAccount(account); err != nil {
		s.logger.Printf("Failed to update account: %v", err)
		return err
	}
	return nil
}

func (s *accountService) Delete(id int) error {
	if err := s.repo.DeleteAccount(id); err != nil {
		s.logger.Printf("Failed to delete account: %v", err)
		return err
	}
	return nil
}

func (s *accountService) List() ([]*models.Account, error) {
	accounts, err := s.repo.ListAccounts()
	if err != nil {
		s.logger.Printf("Failed to list accounts: %v", err)
		return nil, err
	}
	return accounts, nil
}
