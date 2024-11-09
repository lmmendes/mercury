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

func (s *TokensService) List(ctx context.Context, limit, offset int) (*models.PaginatedResponse, error) {
	s.core.Logger.Info("Listing tokens with limit: %d and offset: %d", limit, offset)

	tokens, total, err := s.core.Repository.ListTokens(ctx, limit, offset)
	if err != nil {
		s.core.Logger.Error("Failed to list projects: %v", err)
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

	s.core.Logger.Info("Successfully retrieved %d tokens (total: %d)", len(tokens), total)
	return response, nil
}

func (s *TokensService) Get(ctx context.Context, id int) (*models.Token, error) {
	s.core.Logger.Debug("Fetching token with ID: %d", id)

	token, err := s.core.Repository.GetToken(ctx, id)
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
