package core

import (
	"context"
	"inbox451/internal/models"
)

type TokensService struct {
	core *Core
}

func NewTokensService(core *Core) TokensService {
	return TokensService{core: core}
}

func (s *TokensService) ListByUser(ctx context.Context, user_id int, limit, offset int) (*models.PaginatedResponse, error) {
	s.core.Logger.Info("Listing tokens for user_id=%d with limit: %d and offset: %d", user_id, limit, offset)

	tokens, total, err := s.core.Repository.ListTokensByUser(ctx, user_id, limit, offset)
	if err != nil {
		s.core.Logger.Error("Failed to list tokens for user_id=%d: %v", user_id, err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: tokens,
		Pagination: models.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}

	s.core.Logger.Info("Successfully retrieved %d tokens for user_id=%d (total: %d)", len(tokens), user_id, total)
	return response, nil
}

func (s *TokensService) Get(ctx context.Context, id int, user_id int) (*models.Token, error) {
	s.core.Logger.Debug("Fetching token with ID: %d", id)

	token, err := s.core.Repository.GetTokenByUser(ctx, id, user_id)
	if err != nil {
		s.core.Logger.Error("Failed to fetch token: %v", err)
		return nil, err
	}

	if token == nil {
		s.core.Logger.Info("Token not found with ID: %d", id)
		return nil, ErrNotFound
	}

	return token, nil
}
