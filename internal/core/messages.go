package core

import (
	"context"
	"mercury/internal/models"
)

type MessageService struct {
	core *Core
}

func NewMessageService(core *Core) MessageService {
	return MessageService{core: core}
}

func (s *MessageService) Store(ctx context.Context, message *models.Message) error {
	s.core.Logger.Info("Storing new message for inbox %d from %s", message.InboxID, message.Sender)

	if err := s.core.Repository.CreateMessage(ctx, message); err != nil {
		s.core.Logger.Error("Failed to store message: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully stored message with ID: %d", message.ID)
	return nil
}

func (s *MessageService) Get(ctx context.Context, id int) (*models.Message, error) {
	s.core.Logger.Debug("Fetching message with ID: %d", id)

	message, err := s.core.Repository.GetMessage(ctx, id)
	if err != nil {
		s.core.Logger.Error("Failed to fetch message: %v", err)
		return nil, err
	}

	if message == nil {
		s.core.Logger.Info("Message not found with ID: %d", id)
		return nil, ErrNotFound
	}

	return message, nil
}

func (s *MessageService) ListByInbox(ctx context.Context, inboxID, limit, offset int) (*models.PaginatedResponse, error) {
	s.core.Logger.Info("Listing messages for inbox %d with limit: %d and offset: %d", inboxID, limit, offset)

	messages, total, err := s.core.Repository.ListMessagesByInbox(ctx, inboxID, limit, offset)
	if err != nil {
		s.core.Logger.Error("Failed to list messages: %v", err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: messages,
	}
	response.Pagination.Total = total
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	s.core.Logger.Info("Successfully retrieved %d messages (total: %d)", len(messages), total)
	return response, nil
}
