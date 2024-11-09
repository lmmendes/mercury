package core

import (
	"context"
	"mercury/internal/models"
)

type ProjectService struct {
	core *Core
}

func NewProjectService(core *Core) ProjectService {
	return ProjectService{core: core}
}

func (s *ProjectService) Create(ctx context.Context, project *models.Project) error {
	s.core.Logger.Info("Creating new project: %s", project.Name)

	if err := s.core.Repository.CreateProject(ctx, project); err != nil {
		s.core.Logger.Error("Failed to create project: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully created account with ID: %d", project.ID)
	return nil
}

func (s *ProjectService) Get(ctx context.Context, id int) (*models.Project, error) {
	s.core.Logger.Debug("Fetching project with ID: %d", id)

	account, err := s.core.Repository.GetProject(ctx, id)
	if err != nil {
		s.core.Logger.Error("Failed to fetch project: %v", err)
		return nil, err
	}

	if account == nil {
		s.core.Logger.Info("Project not found with ID: %d", id)
		return nil, ErrNotFound
	}

	return account, nil
}

func (s *ProjectService) Update(ctx context.Context, project *models.Project) error {
	s.core.Logger.Info("Updating project with ID: %d", project.ID)

	if err := s.core.Repository.UpdateProject(ctx, project); err != nil {
		s.core.Logger.Error("Failed to update project: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully updated project with ID: %d", project.ID)
	return nil
}

func (s *ProjectService) Delete(ctx context.Context, id int) error {
	s.core.Logger.Info("Deleting project with ID: %d", id)

	if err := s.core.Repository.DeleteProject(ctx, id); err != nil {
		s.core.Logger.Error("Failed to delete project: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully deleted project with ID: %d", id)
	return nil
}

func (s *ProjectService) List(ctx context.Context, limit, offset int) (*models.PaginatedResponse, error) {
	s.core.Logger.Info("Listing projects with limit: %d and offset: %d", limit, offset)

	projects, total, err := s.core.Repository.ListProjects(ctx, limit, offset)
	if err != nil {
		s.core.Logger.Error("Failed to list projects: %v", err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: projects,
	}
	response.Pagination.Total = total
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	s.core.Logger.Info("Successfully retrieved %d accounts (total: %d)", len(projects), total)
	return response, nil
}
