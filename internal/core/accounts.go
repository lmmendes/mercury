package core

import (
	"context"
	"mercury/internal/models"
)

type AccountService struct {
	core *Core
}

func NewAccountService(core *Core) AccountService {
	return AccountService{core: core}
}

func (s *AccountService) Create(ctx context.Context, account *models.Account) error {
	s.core.Logger.Info("Creating new account: %s", account.Name)

	if err := s.core.Repository.CreateAccount(ctx, account); err != nil {
		s.core.Logger.Error("Failed to create account: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully created account with ID: %d", account.ID)
	return nil
}

func (s *AccountService) Get(ctx context.Context, id int) (*models.Account, error) {
	s.core.Logger.Debug("Fetching account with ID: %d", id)

	account, err := s.core.Repository.GetAccount(ctx, id)
	if err != nil {
		s.core.Logger.Error("Failed to fetch account: %v", err)
		return nil, err
	}

	if account == nil {
		s.core.Logger.Info("Account not found with ID: %d", id)
		return nil, ErrNotFound
	}

	return account, nil
}

func (s *AccountService) Update(ctx context.Context, account *models.Account) error {
	s.core.Logger.Info("Updating account with ID: %d", account.ID)

	if err := s.core.Repository.UpdateAccount(ctx, account); err != nil {
		s.core.Logger.Error("Failed to update account: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully updated account with ID: %d", account.ID)
	return nil
}

func (s *AccountService) Delete(ctx context.Context, id int) error {
	s.core.Logger.Info("Deleting account with ID: %d", id)

	if err := s.core.Repository.DeleteAccount(ctx, id); err != nil {
		s.core.Logger.Error("Failed to delete account: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully deleted account with ID: %d", id)
	return nil
}

func (s *AccountService) List(ctx context.Context, limit, offset int) (*models.PaginatedResponse, error) {
	s.core.Logger.Info("Listing accounts with limit: %d and offset: %d", limit, offset)

	accounts, total, err := s.core.Repository.ListAccounts(ctx, limit, offset)
	if err != nil {
		s.core.Logger.Error("Failed to list accounts: %v", err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: accounts,
	}
	response.Pagination.Total = total
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	s.core.Logger.Info("Successfully retrieved %d accounts (total: %d)", len(accounts), total)
	return response, nil
}
