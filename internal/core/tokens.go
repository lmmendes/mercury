package core

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"inbox451/internal/models"
)

type TokenService struct {
	core *Core
}

func NewTokensService(core *Core) TokenService {
	return TokenService{core: core}
}

func (s *TokenService) ListByUser(ctx context.Context, userId int, limit, offset int) (*models.PaginatedResponse, error) {
	s.core.Logger.Info("Listing tokens for userId %d with limit: %d and offset: %d", userId, limit, offset)

	tokens, total, err := s.core.Repository.ListTokensByUser(ctx, userId, limit, offset)
	if err != nil {
		s.core.Logger.Error("Failed to list tokens for userId %d: %v", userId, err)
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

	s.core.Logger.Info("Successfully retrieved %d tokens for userID %d (total: %d)", len(tokens), userId, total)
	return response, nil
}

func (s *TokenService) GetByUser(ctx context.Context, tokenID int, userID int) (*models.Token, error) {
	s.core.Logger.Debug("Fetching token with ID: %d for userID: %d ", tokenID, userID)

	token, err := s.core.Repository.GetTokenByUser(ctx, tokenID, userID)
	if err != nil {
		s.core.Logger.Error("Failed to fetch token: %v", err)
		return nil, err
	}

	if token == nil {
		s.core.Logger.Info("Token not found with ID: %d for userID %d", tokenID, userID)
		return nil, ErrNotFound
	}

	return token, nil
}

func (s *TokenService) CreateForUser(ctx context.Context, userID int, token *models.Token) (*models.Token, error) {
	s.core.Logger.Debug("Creating token for userId: %d", userID)

	newToken := models.Token{}
	newToken.UserID = userID
	newToken.Token = generateSecureTokenBase64()

	err := s.core.Repository.CreateToken(ctx, &newToken)
	if err != nil {
		return nil, err
	}

	return &newToken, nil
}

func (s *TokenService) DeleteByUser(ctx context.Context, userID int, tokenID int) error {
	s.core.Logger.Debug("Deleting token with ID: %d for userID %d", tokenID, userID)

	// Check if token exists for this user
	_, err := s.GetByUser(ctx, userID, tokenID)
	if err != nil {
		return err
	}

	if err := s.core.Repository.DeleteToken(ctx, tokenID); err != nil {
		s.core.Logger.Error("Failed to delete token: %v", err)
		return err
	}

	s.core.Logger.Debug("Successfully deleted token with ID: %d for userId %d", tokenID, userID)
	return nil
}

// generateSecureTokenBase64 generates a cryptographically secure random token
// encoded in URL-safe base64. It returns a string of approximately 43 characters
// (for 32 bytes of entropy) that is safe for use in URLs and file names.
//
// The generated token uses the following format:
//   - 32 bytes of random data from crypto/rand
//   - URL-safe base64 encoding
//   - No padding characters
//
// Example output: "xJ_dwq8k-rLp5xGhq2d4mNvKzHjY3bWl1fTnM9iR0oE"
//
// If the random number generator fails, it returns an empty string and an error.
func generateSecureTokenBase64() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}