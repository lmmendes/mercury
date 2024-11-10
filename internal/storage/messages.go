package storage

import (
	"context"
	"database/sql"
	"errors"
	"inbox451/internal/models"
)

func (r *repository) CreateMessage(ctx context.Context, message *models.Message) error {
	return r.queries.CreateMessage.QueryRowContext(ctx,
		message.InboxID, message.Sender, message.Receiver, message.Subject, message.Body).
		Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)
}

func (r *repository) GetMessage(ctx context.Context, id int) (*models.Message, error) {
	var message models.Message
	err := r.queries.GetMessage.GetContext(ctx, &message, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *repository) ListMessagesByInbox(ctx context.Context, inboxID, limit, offset int) ([]*models.Message, int, error) {
	var total int
	err := r.queries.CountMessagesByInbox.GetContext(ctx, &total, inboxID)
	if err != nil {
		return nil, 0, err
	}

	var messages []*models.Message
	err = r.queries.ListMessagesByInbox.SelectContext(ctx, &messages, inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}
