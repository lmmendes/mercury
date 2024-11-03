package core

import (
	"context"
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type AccountService interface {
	Create(ctx context.Context, account *models.Account) error
	Get(ctx context.Context, id int) (*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) (*models.PaginatedResponse, error)
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

func (s *accountService) Create(ctx context.Context, account *models.Account) error {
	if err := s.repo.CreateAccount(ctx, account); err != nil {
		s.logger.Error("Failed to create account: %v", err)
		return err
	}
	s.logger.Info("Created account: %d", account.ID)
	return nil
}

func (s *accountService) Get(ctx context.Context, id int) (*models.Account, error) {
	account, err := s.repo.GetAccount(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get account %d: %v", id, err)
		return nil, err
	}
	s.logger.Debug("Retrieved account: %d", id)
	return account, nil
}

func (s *accountService) Update(ctx context.Context, account *models.Account) error {
	if err := s.repo.UpdateAccount(ctx, account); err != nil {
		s.logger.Error("Failed to update account %d: %v", account.ID, err)
		return err
	}
	s.logger.Info("Updated account: %d", account.ID)
	return nil
}

func (s *accountService) Delete(ctx context.Context, id int) error {
	if err := s.repo.DeleteAccount(ctx, id); err != nil {
		s.logger.Error("Failed to delete account %d: %v", id, err)
		return err
	}
	s.logger.Info("Deleted account: %d", id)
	return nil
}

func (s *accountService) List(ctx context.Context, limit, offset int) (*models.PaginatedResponse, error) {
	accounts, total, err := s.repo.ListAccounts(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list accounts: %v", err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: accounts,
	}
	response.Pagination.Total = total
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	s.logger.Debug("Retrieved %d accounts (total: %d)", len(accounts), total)
	return response, nil
}
