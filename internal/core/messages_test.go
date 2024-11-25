package core

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"inbox451/internal/logger"
	"inbox451/internal/mocks"
	"inbox451/internal/models"
	"inbox451/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	null "github.com/volatiletech/null/v9"
)

func setupMessageTestCore(t *testing.T) (*Core, *mocks.MockRepository) {
	mockRepo := mocks.NewMockRepository(t)
	logger := logger.New(io.Discard, logger.DEBUG)

	core := &Core{
		Logger:     logger,
		Repository: mockRepo,
	}
	core.MessageService = NewMessageService(core)

	return core, mockRepo
}

func TestMessageService_Store(t *testing.T) {
	tests := []struct {
		name    string
		message *models.Message
		mockFn  func(*mocks.MockRepository)
		wantErr bool
	}{
		{
			name: "successful store",
			message: &models.Message{
				InboxID:  1,
				Sender:   "sender@example.com",
				Receiver: "inbox@example.com",
				Subject:  "Test Subject",
				Body:     "Test Body",
			},
			mockFn: func(m *mocks.MockRepository) {
				m.On("CreateMessage", mock.Anything, mock.AnythingOfType("*models.Message")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			message: &models.Message{
				InboxID:  1,
				Sender:   "sender@example.com",
				Receiver: "inbox@example.com",
				Subject:  "Test Subject",
				Body:     "Test Body",
			},
			mockFn: func(m *mocks.MockRepository) {
				m.On("CreateMessage", mock.Anything, mock.AnythingOfType("*models.Message")).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, mockRepo := setupMessageTestCore(t)
			tt.mockFn(mockRepo)

			err := core.MessageService.Store(context.Background(), tt.message)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMessageService_Get(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		id      int
		mockFn  func(*mocks.MockRepository)
		want    *models.Message
		wantErr bool
		errType error
	}{
		{
			name: "existing message",
			id:   1,
			mockFn: func(m *mocks.MockRepository) {
				m.On("GetMessage", mock.Anything, 1).Return(&models.Message{
					Base: models.Base{
						ID:        1,
						CreatedAt: null.TimeFrom(now),
						UpdatedAt: null.TimeFrom(now),
					},
					InboxID:  1,
					Sender:   "sender@example.com",
					Receiver: "inbox@example.com",
					Subject:  "Test Subject",
					Body:     "Test Body",
				}, nil)
			},
			want: &models.Message{
				Base: models.Base{
					ID:        1,
					CreatedAt: null.TimeFrom(now),
					UpdatedAt: null.TimeFrom(now),
				},
				InboxID:  1,
				Sender:   "sender@example.com",
				Receiver: "inbox@example.com",
				Subject:  "Test Subject",
				Body:     "Test Body",
			},
			wantErr: false,
		},
		{
			name: "non-existent message",
			id:   999,
			mockFn: func(m *mocks.MockRepository) {
				m.On("GetMessage", mock.Anything, 999).Return(nil, storage.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
			errType: storage.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, mockRepo := setupMessageTestCore(t)
			tt.mockFn(mockRepo)

			got, err := core.MessageService.Get(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMessageService_ListByInbox(t *testing.T) {
	tests := []struct {
		name    string
		inboxID int
		limit   int
		offset  int
		mockFn  func(*mocks.MockRepository)
		want    *models.PaginatedResponse
		wantErr bool
	}{
		{
			name:    "successful list",
			inboxID: 1,
			limit:   10,
			offset:  0,
			mockFn: func(m *mocks.MockRepository) {
				messages := []*models.Message{
					{
						Base:     models.Base{ID: 1},
						InboxID:  1,
						Sender:   "sender1@example.com",
						Receiver: "inbox@example.com",
						Subject:  "Subject 1",
						Body:     "Body 1",
					},
					{
						Base:     models.Base{ID: 2},
						InboxID:  1,
						Sender:   "sender2@example.com",
						Receiver: "inbox@example.com",
						Subject:  "Subject 2",
						Body:     "Body 2",
					},
				}
				m.On("ListMessagesByInbox", mock.Anything, 1, 10, 0).Return(messages, 2, nil)
			},
			want: &models.PaginatedResponse{
				Data: []*models.Message{
					{
						Base:     models.Base{ID: 1},
						InboxID:  1,
						Sender:   "sender1@example.com",
						Receiver: "inbox@example.com",
						Subject:  "Subject 1",
						Body:     "Body 1",
					},
					{
						Base:     models.Base{ID: 2},
						InboxID:  1,
						Sender:   "sender2@example.com",
						Receiver: "inbox@example.com",
						Subject:  "Subject 2",
						Body:     "Body 2",
					},
				},
				Pagination: models.Pagination{
					Total:  2,
					Limit:  10,
					Offset: 0,
				},
			},
			wantErr: false,
		},
		{
			name:    "repository error",
			inboxID: 1,
			limit:   10,
			offset:  0,
			mockFn: func(m *mocks.MockRepository) {
				m.On("ListMessagesByInbox", mock.Anything, 1, 10, 0).
					Return([]*models.Message(nil), 0, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty inbox",
			inboxID: 2,
			limit:   10,
			offset:  0,
			mockFn: func(m *mocks.MockRepository) {
				m.On("ListMessagesByInbox", mock.Anything, 2, 10, 0).
					Return([]*models.Message{}, 0, nil)
			},
			want: &models.PaginatedResponse{
				Data: []*models.Message{},
				Pagination: models.Pagination{
					Total:  0,
					Limit:  10,
					Offset: 0,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, mockRepo := setupMessageTestCore(t)
			tt.mockFn(mockRepo)

			got, err := core.MessageService.ListByInbox(context.Background(), tt.inboxID, tt.limit, tt.offset)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
