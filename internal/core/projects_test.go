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

func setupProjectTestCore(t *testing.T) (*Core, *mocks.Repository) {
	mockRepo := mocks.NewRepository(t)
	logger := logger.New(io.Discard, logger.DEBUG)

	core := &Core{
		Logger:     logger,
		Repository: mockRepo,
	}
	core.ProjectService = NewProjectService(core)

	return core, mockRepo
}

func TestProjectService_Create(t *testing.T) {
	tests := []struct {
		name    string
		project *models.Project
		mockFn  func(*mocks.Repository)
		wantErr bool
	}{
		{
			name: "successful creation",
			project: &models.Project{
				Name: "Test Project",
			},
			mockFn: func(m *mocks.Repository) {
				m.On("CreateProject", mock.Anything, mock.AnythingOfType("*models.Project")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			project: &models.Project{
				Name: "Test Project",
			},
			mockFn: func(m *mocks.Repository) {
				m.On("CreateProject", mock.Anything, mock.AnythingOfType("*models.Project")).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, mockRepo := setupProjectTestCore(t)
			tt.mockFn(mockRepo)

			err := core.ProjectService.Create(context.Background(), tt.project)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_Get(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		id      int
		mockFn  func(*mocks.Repository)
		want    *models.Project
		wantErr bool
		errType error
	}{
		{
			name: "existing project",
			id:   1,
			mockFn: func(m *mocks.Repository) {
				m.On("GetProject", mock.Anything, 1).Return(&models.Project{
					Base: models.Base{
						ID:        1,
						CreatedAt: null.TimeFrom(now),
						UpdatedAt: null.TimeFrom(now),
					},
					Name: "Test Project",
				}, nil)
			},
			want: &models.Project{
				Base: models.Base{
					ID:        1,
					CreatedAt: null.TimeFrom(now),
					UpdatedAt: null.TimeFrom(now),
				},
				Name: "Test Project",
			},
			wantErr: false,
		},
		{
			name: "non-existent project",
			id:   999,
			mockFn: func(m *mocks.Repository) {
				m.On("GetProject", mock.Anything, 999).Return(nil, storage.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
			errType: storage.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, mockRepo := setupProjectTestCore(t)
			tt.mockFn(mockRepo)

			got, err := core.ProjectService.Get(context.Background(), tt.id)
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

func TestProjectService_List(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		offset  int
		mockFn  func(*mocks.Repository)
		want    *models.PaginatedResponse
		wantErr bool
	}{
		{
			name:   "successful list",
			limit:  10,
			offset: 0,
			mockFn: func(m *mocks.Repository) {
				projects := []*models.Project{
					{Base: models.Base{ID: 1}, Name: "Project 1"},
					{Base: models.Base{ID: 2}, Name: "Project 2"},
				}
				m.On("ListProjects", mock.Anything, 10, 0).Return(projects, 2, nil)
			},
			want: &models.PaginatedResponse{
				Data: []*models.Project{
					{Base: models.Base{ID: 1}, Name: "Project 1"},
					{Base: models.Base{ID: 2}, Name: "Project 2"},
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
			name:   "repository error",
			limit:  10,
			offset: 0,
			mockFn: func(m *mocks.Repository) {
				m.On("ListProjects", mock.Anything, 10, 0).Return([]*models.Project(nil), 0, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, mockRepo := setupProjectTestCore(t)
			tt.mockFn(mockRepo)

			got, err := core.ProjectService.List(context.Background(), tt.limit, tt.offset)
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

func TestProjectService_Update(t *testing.T) {
	tests := []struct {
		name    string
		project *models.Project
		mockFn  func(*mocks.Repository)
		wantErr bool
	}{
		{
			name: "successful update",
			project: &models.Project{
				Base: models.Base{ID: 1},
				Name: "Updated Project",
			},
			mockFn: func(m *mocks.Repository) {
				m.On("UpdateProject", mock.Anything, mock.AnythingOfType("*models.Project")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "update non-existent project",
			project: &models.Project{
				Base: models.Base{ID: 999},
				Name: "Updated Project",
			},
			mockFn: func(m *mocks.Repository) {
				m.On("UpdateProject", mock.Anything, mock.AnythingOfType("*models.Project")).
					Return(storage.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, mockRepo := setupProjectTestCore(t)
			tt.mockFn(mockRepo)

			err := core.ProjectService.Update(context.Background(), tt.project)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		mockFn  func(*mocks.Repository)
		wantErr bool
	}{
		{
			name: "successful deletion",
			id:   1,
			mockFn: func(m *mocks.Repository) {
				m.On("DeleteProject", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "delete non-existent project",
			id:   999,
			mockFn: func(m *mocks.Repository) {
				m.On("DeleteProject", mock.Anything, 999).Return(storage.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, mockRepo := setupProjectTestCore(t)
			tt.mockFn(mockRepo)

			err := core.ProjectService.Delete(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
