package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mercury/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository interface {
	// Account operations
	CreateProject(ctx context.Context, account *models.Project) error
	GetProject(ctx context.Context, id int) (*models.Project, error)
	UpdateProject(ctx context.Context, account *models.Project) error
	DeleteProject(ctx context.Context, id int) error
	ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, int, error)

	// Inbox operations
	CreateInbox(ctx context.Context, inbox *models.Inbox) error
	GetInbox(ctx context.Context, id int) (*models.Inbox, error)
	UpdateInbox(ctx context.Context, inbox *models.Inbox) error
	DeleteInbox(ctx context.Context, id int) error
	ListInboxesByAccount(ctx context.Context, accountID, limit, offset int) ([]*models.Inbox, int, error)

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
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type repository struct {
	db      *sqlx.DB
	queries *Queries
}

func NewRepository(db *sqlx.DB) (Repository, error) {
	// First initialize tables
	if err := initializeTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

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

func (r *repository) CreateProject(ctx context.Context, project *models.Project) error {
	return r.queries.CreateAccount.QueryRowContext(ctx, project.Name).
		Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

func (r *repository) GetProject(ctx context.Context, id int) (*models.Project, error) {
	var project models.Project
	err := r.queries.GetProject.GetContext(ctx, &project, id)
	return &project, handleDBError(err)
}

func (r *repository) UpdateAccount(ctx context.Context, project *models.Project) error {
	return r.queries.UpdateProject.QueryRowContext(ctx, project.Name, project.ID).
		Scan(&project.UpdatedAt)
}

func (r *repository) DeleteAccount(ctx context.Context, id int) error {
	result, err := r.queries.DeleteAccount.ExecContext(ctx, id)
	if err != nil {
		return handleDBError(err)
	}
	return handleRowsAffected(result)
}

func (r *repository) ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, int, error) {
	var total int
	err := r.queries.CountProjects.GetContext(ctx, &total)
	if err != nil {
		return nil, 0, err
	}

	var projects []*models.Project
	err = r.queries.ListProjects.SelectContext(ctx, &projects, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *repository) CreateInbox(ctx context.Context, inbox *models.Inbox) error {
	return r.queries.CreateInbox.QueryRowContext(ctx, inbox.projectID, inbox.Email).
		Scan(&inbox.ID, &inbox.CreatedAt, &inbox.UpdatedAt)
}

func (r *repository) GetInbox(ctx context.Context, id int) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.queries.GetInbox.GetContext(ctx, &inbox, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *repository) UpdateInbox(ctx context.Context, inbox *models.Inbox) error {
	result, err := r.queries.UpdateInbox.ExecContext(ctx, inbox.Email, inbox.ID)
	if err != nil {
		return handleDBError(err)
	}
	return handleRowsAffected(result)
}

func (r *repository) DeleteInbox(ctx context.Context, id int) error {
	result, err := r.queries.DeleteInbox.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("inbox not found")
	}
	return nil
}

func (r *repository) ListInboxesByAccount(ctx context.Context, accountID, limit, offset int) ([]*models.Inbox, int, error) {
	var total int
	err := r.queries.CountInboxesByAccount.GetContext(ctx, &total, accountID)
	if err != nil {
		return nil, 0, err
	}

	var inboxes []*models.Inbox
	err = r.queries.ListInboxesByAccount.SelectContext(ctx, &inboxes, accountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return inboxes, total, nil
}

func (r *repository) CreateRule(ctx context.Context, rule *models.ForwardRule) error {
	return r.queries.CreateRule.QueryRowContext(ctx, rule.InboxID, rule.Sender, rule.Receiver, rule.Subject).
		Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

func (r *repository) GetRule(ctx context.Context, id int) (*models.ForwardRule, error) {
	var rule models.ForwardRule
	err := r.queries.GetRule.GetContext(ctx, &rule, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *repository) UpdateRule(ctx context.Context, rule *models.ForwardRule) error {
	result, err := r.queries.UpdateRule.ExecContext(ctx, rule.Sender, rule.Receiver, rule.Subject, rule.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("rule not found")
	}
	return nil
}

func (r *repository) DeleteRule(ctx context.Context, id int) error {
	result, err := r.queries.DeleteRule.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("rule not found")
	}
	return nil
}

func (r *repository) ListRulesByInbox(ctx context.Context, inboxID, limit, offset int) ([]*models.ForwardRule, int, error) {
	var total int
	err := r.queries.CountRulesByInbox.GetContext(ctx, &total, inboxID)
	if err != nil {
		return nil, 0, err
	}

	var rules []*models.ForwardRule
	err = r.queries.ListRulesByInbox.SelectContext(ctx, &rules, inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) CreateMessage(ctx context.Context, message *models.Message) error {
	return r.queries.CreateMessage.QueryRowContext(ctx,
		message.InboxID, message.Sender, message.Receiver, message.Subject, message.Body).
		Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)
}

func (r *repository) GetMessage(ctx context.Context, id int) (*models.Message, error) {
	var message models.Message
	err := r.queries.GetMessage.GetContext(ctx, &message, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *repository) ListMessagesByInbox(ctx context.Context, inboxID, limit, offset int) ([]*models.Message, int, error) {
	var total int
	err := r.queries.CountMessagesByInbox.GetContext(ctx, &total, inboxID)
	if err != nil {
		return nil, 0, err
	}

	var messages []*models.Message
	err = r.queries.ListMessagesByInbox.SelectContext(ctx, &messages, inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *repository) ListRules(ctx context.Context, limit, offset int) ([]*models.ForwardRule, int, error) {
	var total int
	err := r.queries.CountRules.GetContext(ctx, &total)
	if err != nil {
		return nil, 0, err
	}

	var rules []*models.ForwardRule
	err = r.queries.ListRules.SelectContext(ctx, &rules, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) GetInboxByEmail(ctx context.Context, email string) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.queries.GetInboxByEmail.GetContext(ctx, &inbox, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
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

func (r *repository) GetUser(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := r.queries.GetUser.GetContext(ctx, &user, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
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

// Update the initializeTables function
func initializeTables(db *sqlx.DB) error {
	// Read the schema file
	schemaBytes, err := queriesFS.ReadFile("schema.sql")

	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Start a transaction
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure proper rollback handling
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			// We can only log this error since we can't return it
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	// Execute the schema
	if _, err := tx.Exec(string(schemaBytes)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
