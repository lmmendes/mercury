package core

import (
	"mercury/internal/logger"
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
	logger *logger.Logger
}

func NewAccountService(core *Core) AccountService {
	return &accountService{
		repo:   core.Repository,
		logger: core.Logger,
	}
}

func (s *accountService) Create(account *models.Account) error {
	if err := s.repo.CreateAccount(account); err != nil {
		s.logger.Error("Failed to create account: %v", err)
		return err
	}
	s.logger.Info("Created account: %d", account.ID)
	return nil
}

func (s *accountService) Get(id int) (*models.Account, error) {
	account, err := s.repo.GetAccount(id)
	if err != nil {
		s.logger.Error("Failed to get account %d: %v", id, err)
		return nil, err
	}
	s.logger.Debug("Retrieved account: %d", id)
	return account, nil
}

func (s *accountService) Update(account *models.Account) error {
	if err := s.repo.UpdateAccount(account); err != nil {
		s.logger.Error("Failed to update account %d: %v", account.ID, err)
		return err
	}
	s.logger.Info("Updated account: %d", account.ID)
	return nil
}

func (s *accountService) Delete(id int) error {
	if err := s.repo.DeleteAccount(id); err != nil {
		s.logger.Error("Failed to delete account %d: %v", id, err)
		return err
	}
	s.logger.Info("Deleted account: %d", id)
	return nil
}

func (s *accountService) List() ([]*models.Account, error) {
	accounts, err := s.repo.ListAccounts()
	if err != nil {
		s.logger.Error("Failed to list accounts: %v", err)
		return nil, err
	}
	s.logger.Debug("Retrieved %d accounts", len(accounts))
	return accounts, nil
}
