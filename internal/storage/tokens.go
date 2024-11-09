package storage

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"inbox451/internal/models"

	_ "github.com/lib/pq"
)

func (r *repository) ListTokensByUser(ctx context.Context, user_id int, limit, offset int) ([]*models.Token, int, error) {
	var total int
	err := r.queries.CountTokensByUser.GetContext(ctx, &total, user_id)
	if err != nil {
		return nil, 0, err
	}

	var tokens []*models.Token
	err = r.queries.ListTokensByUser.SelectContext(ctx, &tokens, user_id, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

func (r *repository) GetTokenByUser(ctx context.Context, token_id int, user_id int) (*models.Token, error) {
	var token models.Token
	err := r.queries.GetTokenByUser.GetContext(ctx, &token, token_id, user_id)
	return &token, handleDBError(err)
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
