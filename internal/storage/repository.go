package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"mercury/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository interface {
	// Account operations
	CreateAccount(account *models.Account) error
	GetAccount(id int) (*models.Account, error)
	UpdateAccount(account *models.Account) error
	DeleteAccount(id int) error
	ListAccounts(limit, offset int) ([]*models.Account, int, error)

	// Inbox operations
	CreateInbox(inbox *models.Inbox) error
	GetInbox(id int) (*models.Inbox, error)
	UpdateInbox(inbox *models.Inbox) error
	DeleteInbox(id int) error
	ListInboxesByAccount(accountID, limit, offset int) ([]*models.Inbox, int, error)

	// Rule operations
	CreateRule(rule *models.Rule) error
	GetRule(id int) (*models.Rule, error)
	UpdateRule(rule *models.Rule) error
	DeleteRule(id int) error
	ListRulesByInbox(inboxID, limit, offset int) ([]*models.Rule, int, error)

	// Message operations
	CreateMessage(message *models.Message) error
	GetMessage(id int) (*models.Message, error)
	ListMessagesByInbox(inboxID, limit, offset int) ([]*models.Message, int, error)
	ListRules(limit, offset int) ([]*models.Rule, int, error)
	GetInboxByEmail(email string) (*models.Inbox, error)

	// User operations
	CreateUser(user *models.User) error
	GetUser(id int) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	GetUserByUsername(username string) (*models.User, error)

	// Initialize tables
	InitializeTables() error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAccount(account *models.Account) error {
	query := `
		INSERT INTO accounts (name, created_at, updated_at)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, account.Name).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
}

func (r *repository) GetAccount(id int) (*models.Account, error) {
	var account models.Account
	err := r.db.Get(&account,
		"SELECT id, name, created_at, updated_at FROM accounts WHERE id = $1",
		id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *repository) UpdateAccount(account *models.Account) error {
	query := `
		UPDATE accounts
		SET name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at`
	result := r.db.QueryRow(query, account.Name, account.ID)
	if err := result.Scan(&account.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("account not found")
		}
		return err
	}
	return nil
}

func (r *repository) DeleteAccount(id int) error {
	result, err := r.db.Exec("DELETE FROM accounts WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (r *repository) ListAccounts(limit, offset int) ([]*models.Account, int, error) {
	var total int
	err := r.db.Get(&total, "SELECT COUNT(*) FROM accounts")
	if err != nil {
		return nil, 0, err
	}

	var accounts []*models.Account
	err = r.db.Select(&accounts,
		"SELECT id, name, created_at, updated_at FROM accounts ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (r *repository) CreateInbox(inbox *models.Inbox) error {
	query := `
		INSERT INTO inboxes (account_id, email, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, inbox.AccountID, inbox.Email).Scan(
		&inbox.ID, &inbox.CreatedAt, &inbox.UpdatedAt)
}

func (r *repository) GetInbox(id int) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.db.Get(&inbox, "SELECT id, account_id, email, created_at, updated_at FROM inboxes WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *repository) UpdateInbox(inbox *models.Inbox) error {
	result, err := r.db.Exec("UPDATE inboxes SET email = $1 WHERE id = $2", inbox.Email, inbox.ID)
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

func (r *repository) DeleteInbox(id int) error {
	result, err := r.db.Exec("DELETE FROM inboxes WHERE id = $1", id)
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

func (r *repository) ListInboxesByAccount(accountID, limit, offset int) ([]*models.Inbox, int, error) {
	var total int
	err := r.db.Get(&total, "SELECT COUNT(*) FROM inboxes WHERE account_id = $1", accountID)
	if err != nil {
		return nil, 0, err
	}

	var inboxes []*models.Inbox
	err = r.db.Select(&inboxes,
		"SELECT id, account_id, email, created_at, updated_at FROM inboxes WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3",
		accountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return inboxes, total, nil
}

func (r *repository) CreateRule(rule *models.Rule) error {
	query := `
		INSERT INTO rules (inbox_id, sender, receiver, subject, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, rule.InboxID, rule.Sender, rule.Receiver, rule.Subject).Scan(
		&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

func (r *repository) GetRule(id int) (*models.Rule, error) {
	var rule models.Rule
	err := r.db.Get(&rule, "SELECT id, inbox_id, sender, receiver, subject, created_at, updated_at FROM rules WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *repository) UpdateRule(rule *models.Rule) error {
	result, err := r.db.Exec("UPDATE rules SET sender = $1, receiver = $2, subject = $3 WHERE id = $4",
		rule.Sender, rule.Receiver, rule.Subject, rule.ID)
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

func (r *repository) DeleteRule(id int) error {
	result, err := r.db.Exec("DELETE FROM rules WHERE id = $1", id)
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

func (r *repository) ListRulesByInbox(inboxID, limit, offset int) ([]*models.Rule, int, error) {
	var total int
	err := r.db.Get(&total, "SELECT COUNT(*) FROM rules WHERE inbox_id = $1", inboxID)
	if err != nil {
		return nil, 0, err
	}

	var rules []*models.Rule
	err = r.db.Select(&rules,
		"SELECT id, inbox_id, sender, receiver, subject, created_at, updated_at FROM rules WHERE inbox_id = $1 ORDER BY id LIMIT $2 OFFSET $3",
		inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) CreateMessage(message *models.Message) error {
	query := `
		INSERT INTO messages (inbox_id, sender, receiver, subject, body, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query,
		message.InboxID, message.Sender, message.Receiver, message.Subject, message.Body).Scan(
		&message.ID, &message.CreatedAt, &message.UpdatedAt)
}

func (r *repository) GetMessage(id int) (*models.Message, error) {
	var message models.Message
	err := r.db.Get(&message, "SELECT id, inbox_id, sender, receiver, subject, body, created_at, updated_at FROM messages WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *repository) ListMessagesByInbox(inboxID, limit, offset int) ([]*models.Message, int, error) {
	var total int
	err := r.db.Get(&total, "SELECT COUNT(*) FROM messages WHERE inbox_id = $1", inboxID)
	if err != nil {
		return nil, 0, err
	}

	var messages []*models.Message
	err = r.db.Select(&messages,
		"SELECT id, inbox_id, sender, receiver, subject, body, created_at, updated_at FROM messages WHERE inbox_id = $1 ORDER BY id LIMIT $2 OFFSET $3",
		inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *repository) ListRules(limit, offset int) ([]*models.Rule, int, error) {
	var total int
	err := r.db.Get(&total, "SELECT COUNT(*) FROM rules")
	if err != nil {
		return nil, 0, err
	}

	var rules []*models.Rule
	err = r.db.Select(&rules,
		"SELECT id, inbox_id, sender, receiver, subject, created_at, updated_at FROM rules ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) GetInboxByEmail(email string) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.db.Get(&inbox, "SELECT id, account_id, email, created_at, updated_at FROM inboxes WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *repository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (
			name, username, password, email, status, kind,
			password_login, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		user.Name,
		user.Username,
		user.Password,
		user.Email,
		user.Status,
		user.Kind,
		user.PasswordLogin,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *repository) GetUser(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, `
		SELECT id, name, username, password, email, status, kind,
			   password_login, loggedin_at, created_at, updated_at
		FROM users WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, username = $2, password = $3, email = $4,
			status = $5, kind = $6, password_login = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at`

	result := r.db.QueryRow(
		query,
		user.Name,
		user.Username,
		user.Password,
		user.Email,
		user.Status,
		user.Kind,
		user.PasswordLogin,
		user.ID,
	)

	if err := result.Scan(&user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}

func (r *repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, `
		SELECT id, name, username, password, email, status, kind,
			   password_login, loggedin_at, created_at, updated_at
		FROM users WHERE username = $1`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) DeleteUser(id int) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
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

func (r *repository) InitializeTables() error {
	queries := []string{
		`
		-- Ruby on Rails inspired schema_migrations table
		CREATE TABLE IF NOT EXISTS schema_migrations (
		    version VARCHAR(255) PRIMARY KEY,
		)`,
		`
		-- Create custom types for roles if they don't exist
		DO $$
		BEGIN
		    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
		        CREATE TYPE user_role AS ENUM ('user', 'admin');
		    END IF;

		    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'project_role') THEN
		        CREATE TYPE project_role AS ENUM ('user', 'admin');
		    END IF;
		END
		$$;`,
		`-- Projects table
		CREATE TABLE IF NOT EXISTS projects (
		    id SERIAL PRIMARY KEY,
		    name VARCHAR(100) NOT NULL CHECK (LENGTH(name) >= 2),
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`-- Users table
		CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY,
		    name VARCHAR(255) NOT NULL,
		    username VARCHAR(255) NOT NULL UNIQUE,
		    password VARCHAR(255) NOT NULL,
		    email VARCHAR(255) NOT NULL UNIQUE,
		    status VARCHAR(50) NOT NULL,
		    role user_role NOT NULL DEFAULT 'user',
		    password_login BOOLEAN NOT NULL DEFAULT true,
		    loggedin_at TIMESTAMP WITH TIME ZONE,
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`-- Project Users junction table
		CREATE TABLE IF NOT EXISTS project_users (
		    id SERIAL PRIMARY KEY,
		    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
		    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    role project_role NOT NULL DEFAULT 'user',
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    UNIQUE(project_id, user_id)
		)`,
		`-- Inboxes table
		CREATE TABLE IF NOT EXISTS inboxes (
		    id SERIAL PRIMARY KEY,
		    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
		    email VARCHAR(255) NOT NULL UNIQUE,
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`-- User Tokens table
		CREATE TABLE IF NOT EXISTS user_tokens (
		    id SERIAL PRIMARY KEY,
		    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    token VARCHAR(255) NOT NULL UNIQUE,
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`-- Forward Rules table
		CREATE TABLE IF NOT EXISTS forward_rules (
		    id SERIAL PRIMARY KEY,
		    inbox_id INTEGER NOT NULL REFERENCES inboxes(id) ON DELETE CASCADE,
		    sender VARCHAR(255),
		    receiver VARCHAR(255),
		    subject VARCHAR(200),
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    CHECK (sender IS NOT NULL OR receiver IS NOT NULL OR subject IS NOT NULL)
		)`,
		`-- Messages table
		CREATE TABLE IF NOT EXISTS messages (
		    id SERIAL PRIMARY KEY,
		    inbox_id INTEGER NOT NULL REFERENCES inboxes(id) ON DELETE CASCADE,
		    sender VARCHAR(255) NOT NULL,
		    receiver VARCHAR(255) NOT NULL,
		    subject VARCHAR(200) NOT NULL,
		    body TEXT NOT NULL,
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				return fmt.Errorf("failed to execute query: %w, failed to rollback transaction: %v", err, rbErr)
			}
			return err
		}
	}

	return tx.Commit()
}
