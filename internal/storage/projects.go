package storage

import (
	"context"
	"inbox451/internal/models"
)

func (r *repository) ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, int, error) {
	var total int
	err := r.queries.CountProjects.GetContext(ctx, &total)
	if err != nil {
		return nil, 0, err
	}

	var projects []*models.Project = make([]*models.Project, 0)
	err = r.queries.ListProjects.SelectContext(ctx, &projects, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *repository) ListProjectsByUser(ctx context.Context, userID int, limit int, offset int) ([]*models.Project, int, error) {
	var total int
	err := r.queries.CountProjectsByUser.GetContext(ctx, &total, userID)
	if err != nil {
		return nil, 0, err
	}

	var projects []*models.Project = make([]*models.Project, 0)
	err = r.queries.ListProjectsByUser.SelectContext(ctx, &projects, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *repository) GetProject(ctx context.Context, id int) (*models.Project, error) {
	var project models.Project
	err := r.queries.GetProject.GetContext(ctx, &project, id)
	return &project, handleDBError(err)
}

func (r *repository) CreateProject(ctx context.Context, project *models.Project) error {
	return r.queries.CreateProject.QueryRowContext(ctx, project.Name).
		Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

func (r *repository) UpdateProject(ctx context.Context, project *models.Project) error {
	return r.queries.UpdateProject.QueryRowContext(ctx, project.Name, project.ID).
		Scan(&project.UpdatedAt)
}

func (r *repository) ProjectAddUser(ctx context.Context, projectUser *models.ProjectUser) error {
	return r.queries.AddUserToProject.QueryRowContext(ctx, projectUser.ProjectID, projectUser.UserID, projectUser.Role).
		Scan(&projectUser.CreatedAt, &projectUser.UpdatedAt)
}

func (r *repository) DeleteProject(ctx context.Context, id int) error {
	result, err := r.queries.DeleteProject.ExecContext(ctx, id)
	if err != nil {
		return handleDBError(err)
	}
	return handleRowsAffected(result)
}

func (r *repository) ProjectRemoveUser(ctx context.Context, projectID int, userID int) error {
	result, err := r.queries.RemoveUserFromProject.ExecContext(ctx, projectID, userID)
	if err != nil {
		return handleDBError(err)
	}
	return handleRowsAffected(result)
}