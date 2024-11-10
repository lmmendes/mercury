package storage

import (
	"context"
	"fmt"
	"inbox451/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository interface {
	// Project operations
	CreateProject(ctx context.Context, project *models.Project) error
	GetProject(ctx context.Context, id int) (*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id int) error
	ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, int, error)

	// Inbox operations
	CreateInbox(ctx context.Context, inbox *models.Inbox) error
	GetInbox(ctx context.Context, id int) (*models.Inbox, error)
	UpdateInbox(ctx context.Context, inbox *models.Inbox) error
	DeleteInbox(ctx context.Context, id int) error
	ListInboxesByProject(ctx context.Context, projectID, limit, offset int) ([]*models.Inbox, int, error)

	// Rule operations
	CreateRule(ctx context.Context, rule *models.ForwardRule) error
	GetRule(ctx context.Context, id int) (*models.ForwardRule, error)
	UpdateRule(ctx context.Context, rule *models.ForwardRule) error
	DeleteRule(ctx context.Context, id int) error
	ListRulesByInbox(ctx context.Context, inboxID, limit, offset int) ([]*models.ForwardRule, int, error)

	// Message operations
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessage(ctx context.Context, id int) (*models.Message, error)
	ListMessagesByInbox(ctx context.Context, inboxID, limit, offset int) ([]*models.Message, int, error)
	ListRules(ctx context.Context, limit, offset int) ([]*models.ForwardRule, int, error)
	GetInboxByEmail(ctx context.Context, email string) (*models.Inbox, error)

	// User operations
	ListUsers(ctx context.Context, limit, offset int) ([]*models.User, int, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userId int) error

	// Tokens
	ListTokensByUser(ctx context.Context, userID int, limit, offset int) ([]*models.Token, int, error)
	GetTokenByUser(ctx context.Context, userID int, tokenID int) (*models.Token, error)
	CreateToken(ctx context.Context, token *models.Token) error
	DeleteToken(ctx context.Context, tokenID int) error
}

type repository struct {
	db      *sqlx.DB
	queries *Queries
}

func NewRepository(db *sqlx.DB) (Repository, error) {
	// Then prepare queries
	queries, err := PrepareQueries(db)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare queries: %w", err)
	}

	return &repository{
		db:      db,
		queries: queries,
	}, nil
}
