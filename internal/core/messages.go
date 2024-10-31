package core

import (
	"log"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type MessageService interface {
	Store(message *models.Message) error
	GetByInboxID(inboxID int) ([]*models.Message, error)
}

type messageService struct {
	repo   storage.Repository
	logger *log.Logger
}

func NewMessageService(core *Core) MessageService {
	return &messageService{
		repo:   core.Repository,
		logger: core.Logger,
	}
}

func (s *messageService) Store(message *models.Message) error {
	// First try to match against rules
	rules, err := s.repo.ListRules()
	if err != nil {
		s.logger.Printf("Error querying rules: %v", err)
		return err
	}

	for _, rule := range rules {
		if rule.Sender == message.Sender && rule.Receiver == message.Receiver && rule.Subject == message.Subject {
			message.InboxID = rule.InboxID
			return s.repo.CreateMessage(message)
		}
	}

	// If no rule matches, store in the inbox with the matching email address
	inbox, err := s.repo.GetInboxByEmail(message.Receiver)
	if err != nil {
		s.logger.Printf("Error finding inbox: %v", err)
		return err
	}

	message.InboxID = inbox.ID
	return s.repo.CreateMessage(message)
}

func (s *messageService) GetByInboxID(inboxID int) ([]*models.Message, error) {
	return s.repo.ListMessagesByInbox(inboxID)
}
