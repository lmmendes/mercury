package core

import (
	"log"
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
	logger *log.Logger
}

func NewInboxService(core *Core) InboxService {
	return &inboxService{
		repo:   core.Repository,
		logger: core.Logger,
	}
}

func (s *inboxService) Create(inbox *models.Inbox) error {
	return s.repo.CreateInbox(inbox)
}

func (s *inboxService) Get(id int) (*models.Inbox, error) {
	return s.repo.GetInbox(id)
}

func (s *inboxService) Update(inbox *models.Inbox) error {
	return s.repo.UpdateInbox(inbox)
}

func (s *inboxService) Delete(id int) error {
	return s.repo.DeleteInbox(id)
}

func (s *inboxService) GetByAccountID(accountID int) ([]*models.Inbox, error) {
	return s.repo.ListInboxesByAccount(accountID)
}
