package core

import (
	"context"
	"inbox451/internal/models"
)

type UserService struct {
	core *Core
}

func NewUserService(core *Core) UserService {
	return UserService{core: core}
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	s.core.Logger.Info("Creating new project: %s", user.Name)

	if err := s.core.Repository.CreateUser(ctx, user); err != nil {
		s.core.Logger.Error("Failed to create user: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully created user with ID: %d", user.ID)
	return nil
}

func (s *UserService) Get(ctx context.Context, userID int) (*models.User, error) {
	s.core.Logger.Debug("Fetching user with ID: %d", userID)

	project, err := s.core.Repository.GetUser(ctx, userID)
	if err != nil {
		s.core.Logger.Error("Failed to fetch user: %v", err)
		return nil, err
	}

	if project == nil {
		s.core.Logger.Info("User not found with ID: %d", userID)
		return nil, ErrNotFound
	}

	return project, nil
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	s.core.Logger.Info("Updating user with ID: %d", user.ID)

	if err := s.core.Repository.UpdateUser(ctx, user); err != nil {
		s.core.Logger.Error("Failed to user project: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully updated user with ID: %d", user.ID)
	return nil
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	s.core.Logger.Info("Deleting user with ID: %d", id)

	if err := s.core.Repository.DeleteUser(ctx, id); err != nil {
		s.core.Logger.Error("Failed to delete user: %v", err)
		return err
	}

	s.core.Logger.Info("Successfully deleted user with ID: %d", id)
	return nil
}

func (s *UserService) List(ctx context.Context, limit, offset int) (*models.PaginatedResponse, error) {
	s.core.Logger.Info("Listing users with limit: %d and offset: %d", limit, offset)

	users, total, err := s.core.Repository.ListUsers(ctx, limit, offset)
	if err != nil {
		s.core.Logger.Error("Failed to list users: %v", err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: users,
		Pagination: models.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}

	s.core.Logger.Info("Successfully retrieved %d users (total: %d)", len(users), total)
	return response, nil
}
