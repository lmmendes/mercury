package core

import (
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type MessageService interface {
	Store(message *models.Message) error
	GetByInboxID(inboxID int) ([]*models.Message, error)
}

type messageService struct {
	repo   storage.Repository
	logger *logger.Logger
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
		s.logger.Error("Error querying rules: %v", err)
		return err
	}

	for _, rule := range rules {
		if rule.Sender == message.Sender && rule.Receiver == message.Receiver && rule.Subject == message.Subject {
			message.InboxID = rule.InboxID
			s.logger.Debug("Message matched rule for inbox %d", rule.InboxID)
			return s.repo.CreateMessage(message)
		}
	}

	// If no rule matches, store in the inbox with the matching email address
	inbox, err := s.repo.GetInboxByEmail(message.Receiver)
	if err != nil {
		s.logger.Error("Error finding inbox for email %s: %v", message.Receiver, err)
		return err
	}

	message.InboxID = inbox.ID
	s.logger.Info("Storing message in inbox %d", inbox.ID)
	return s.repo.CreateMessage(message)
}

func (s *messageService) GetByInboxID(inboxID int) ([]*models.Message, error) {
	messages, err := s.repo.ListMessagesByInbox(inboxID)
	if err != nil {
		s.logger.Error("Failed to list messages for inbox %d: %v", inboxID, err)
		return nil, err
	}
	s.logger.Debug("Retrieved %d messages for inbox %d", len(messages), inboxID)
	return messages, nil
}
