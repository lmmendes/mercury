package core

import (
	"context"
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type InboxService interface {
	Create(ctx context.Context, inbox *models.Inbox) error
	Get(ctx context.Context, id int) (*models.Inbox, error)
	Update(ctx context.Context, inbox *models.Inbox) error
	Delete(ctx context.Context, id int) error
	ListByAccount(ctx context.Context, accountID, limit, offset int) (*models.PaginatedResponse, error)
}

type inboxService struct {
	repo   storage.Repository
	logger *logger.Logger
}

func NewInboxService(core *Core) InboxService {
	return &inboxService{
		repo:   core.Repository,
		logger: core.Logger,
	}
}

func (s *inboxService) Create(ctx context.Context, inbox *models.Inbox) error {
	if err := s.repo.CreateInbox(ctx, inbox); err != nil {
		s.logger.Error("Failed to create inbox: %v", err)
		return err
	}
	s.logger.Info("Created inbox %d for account %d", inbox.ID, inbox.AccountID)
	return nil
}

func (s *inboxService) Get(ctx context.Context, id int) (*models.Inbox, error) {
	inbox, err := s.repo.GetInbox(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get inbox %d: %v", id, err)
		return nil, err
	}
	s.logger.Debug("Retrieved inbox: %d", id)
	return inbox, nil
}

func (s *inboxService) Update(ctx context.Context, inbox *models.Inbox) error {
	if err := s.repo.UpdateInbox(ctx, inbox); err != nil {
		s.logger.Error("Failed to update inbox %d: %v", inbox.ID, err)
		return err
	}
	s.logger.Info("Updated inbox: %d", inbox.ID)
	return nil
}

func (s *inboxService) Delete(ctx context.Context, id int) error {
	if err := s.repo.DeleteInbox(ctx, id); err != nil {
		s.logger.Error("Failed to delete inbox %d: %v", id, err)
		return err
	}
	s.logger.Info("Deleted inbox: %d", id)
	return nil
}

func (s *inboxService) ListByAccount(ctx context.Context, accountID, limit, offset int) (*models.PaginatedResponse, error) {
	inboxes, total, err := s.repo.ListInboxesByAccount(ctx, accountID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list inboxes for account %d: %v", accountID, err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: inboxes,
	}
	response.Pagination.Total = total
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	s.logger.Debug("Retrieved %d inboxes for account %d (total: %d)", len(inboxes), accountID, total)
	return response, nil
}
