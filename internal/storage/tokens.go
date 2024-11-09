package storage

import (
	"context"
	"inbox451/internal/models"

	_ "github.com/lib/pq"
)

func (r *repository) ListTokens(ctx context.Context, limit, offset int) ([]*models.Token, int, error) {
	var total int
	err := r.queries.CountTokens.GetContext(ctx, &total)
	if err != nil {
		return nil, 0, err
	}

	var tokens []*models.Token
	err = r.queries.ListTokens.SelectContext(ctx, &tokens, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

func (r *repository) GetToken(ctx context.Context, id int) (*models.Token, error) {
	var token models.Token
	err := r.queries.GetToken.GetContext(ctx, &token, id)
	return &token, handleDBError(err)
}
