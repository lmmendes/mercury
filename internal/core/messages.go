package core

import (
	"context"
	"fmt"
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type MessageService interface {
	Store(ctx context.Context, message *models.Message) error
	ListByInbox(ctx context.Context, inboxID, limit, offset int) (*models.PaginatedResponse, error)
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

func (s *messageService) Store(ctx context.Context, message *models.Message) error {
	// First try to match against rules
	rules, _, err := s.repo.ListRules(ctx, 1000, 0)
	if err != nil {
		s.logger.Error("Error querying rules: %v", err)
		return err
	}

	for _, rule := range rules {
		if rule.Sender == message.Sender && rule.Receiver == message.Receiver && rule.Subject == message.Subject {
			message.InboxID = rule.InboxID
			s.logger.Debug("Message matched rule for inbox %d", rule.InboxID)
			return s.repo.CreateMessage(ctx, message)
		}
	}

	// If no rule matches, store in the inbox with the matching email address
	inbox, err := s.repo.GetInboxByEmail(ctx, message.Receiver)
	if err != nil {
		s.logger.Error("Error finding inbox for email %s: %v", message.Receiver, err)
		return err
	}

	if inbox == nil {
		s.logger.Error("No inbox found for email %s", message.Receiver)
		return fmt.Errorf("no inbox found for email address: %s", message.Receiver)
	}

	message.InboxID = inbox.ID
	s.logger.Info("Storing message in inbox %d", inbox.ID)
	return s.repo.CreateMessage(ctx, message)
}

func (s *messageService) ListByInbox(ctx context.Context, inboxID, limit, offset int) (*models.PaginatedResponse, error) {
	messages, total, err := s.repo.ListMessagesByInbox(ctx, inboxID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list messages for inbox %d: %v", inboxID, err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: messages,
	}
	response.Pagination.Total = total
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	s.logger.Debug("Retrieved %d messages for inbox %d (total: %d)", len(messages), inboxID, total)
	return response, nil
}
