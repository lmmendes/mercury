package core

import (
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type InboxService interface {
	Create(inbox *models.Inbox) error
	Get(id int) (*models.Inbox, error)
	Update(inbox *models.Inbox) error
	Delete(id int) error
	GetByAccountID(accountID int) ([]*models.Inbox, error)
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

func (s *inboxService) Create(inbox *models.Inbox) error {
	if err := s.repo.CreateInbox(inbox); err != nil {
		s.logger.Error("Failed to create inbox: %v", err)
		return err
	}
	s.logger.Info("Created inbox %d for account %d", inbox.ID, inbox.AccountID)
	return nil
}

func (s *inboxService) Get(id int) (*models.Inbox, error) {
	inbox, err := s.repo.GetInbox(id)
	if err != nil {
		s.logger.Error("Failed to get inbox %d: %v", id, err)
		return nil, err
	}
	s.logger.Debug("Retrieved inbox: %d", id)
	return inbox, nil
}

func (s *inboxService) Update(inbox *models.Inbox) error {
	if err := s.repo.UpdateInbox(inbox); err != nil {
		s.logger.Error("Failed to update inbox %d: %v", inbox.ID, err)
		return err
	}
	s.logger.Info("Updated inbox: %d", inbox.ID)
	return nil
}

func (s *inboxService) Delete(id int) error {
	if err := s.repo.DeleteInbox(id); err != nil {
		s.logger.Error("Failed to delete inbox %d: %v", id, err)
		return err
	}
	s.logger.Info("Deleted inbox: %d", id)
	return nil
}

func (s *inboxService) GetByAccountID(accountID int) ([]*models.Inbox, error) {
	inboxes, err := s.repo.ListInboxesByAccount(accountID)
	if err != nil {
		s.logger.Error("Failed to list inboxes for account %d: %v", accountID, err)
		return nil, err
	}
	s.logger.Debug("Retrieved %d inboxes for account %d", len(inboxes), accountID)
	return inboxes, nil
}
