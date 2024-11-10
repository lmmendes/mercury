package storage

import (
	"context"
	"database/sql"
	"errors"
	"inbox451/internal/models"

	_ "github.com/lib/pq"
)

func (r *repository) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, int, error) {
	var total int
	err := r.queries.CountUsers.GetContext(ctx, &total)
	if err != nil {
		return nil, 0, err
	}

	var users []*models.User
	err = r.queries.ListUsers.SelectContext(ctx, &users, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *repository) GetUser(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := r.queries.GetUser.GetContext(ctx, &user, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.queries.GetUserByUsername.GetContext(ctx, &user, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	return r.queries.CreateUser.QueryRowContext(ctx,
		user.Name,
		user.Username,
		user.Password,
		user.Email,
		user.Status,
		user.Role,
		user.PasswordLogin).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *repository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.queries.UpdateUser.QueryRowContext(ctx,
		user.Name,
		user.Username,
		user.Password,
		user.Email,
		user.Status,
		user.Role,
		user.PasswordLogin,
		user.ID).
		Scan(&user.UpdatedAt)
}

func (r *repository) DeleteUser(ctx context.Context, id int) error {
	result, err := r.queries.DeleteUser.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}
