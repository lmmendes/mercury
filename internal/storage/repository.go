package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"mercury/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql/v2"
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
	db      *sqlx.DB
	queries *Queries
}

func NewRepository(db *sqlx.DB) (Repository, error) {
	queries, err := PrepareQueries(db)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare queries: %w", err)
	}

	return &repository{
		db:      db,
		queries: queries,
	}, nil
}

func (r *repository) CreateAccount(account *models.Account) error {
	return r.queries.CreateAccount.QueryRow(account.Name).
		Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
}

func (r *repository) GetAccount(id int) (*models.Account, error) {
	var account models.Account
	err := r.queries.GetAccount.Get(&account, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *repository) UpdateAccount(account *models.Account) error {
	return r.queries.UpdateAccount.QueryRow(account.Name, account.ID).
		Scan(&account.UpdatedAt)
}

func (r *repository) DeleteAccount(id int) error {
	result, err := r.queries.DeleteAccount.Exec(id)
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
	err := r.queries.CountAccounts.Get(&total)
	if err != nil {
		return nil, 0, err
	}

	var accounts []*models.Account
	err = r.queries.ListAccounts.Select(&accounts, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (r *repository) CreateInbox(inbox *models.Inbox) error {
	return r.queries.CreateInbox.QueryRow(inbox.AccountID, inbox.Email).
		Scan(&inbox.ID, &inbox.CreatedAt, &inbox.UpdatedAt)
}

func (r *repository) GetInbox(id int) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.queries.GetInbox.Get(&inbox, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *repository) UpdateInbox(inbox *models.Inbox) error {
	result, err := r.queries.UpdateInbox.Exec(inbox.Email, inbox.ID)
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
	result, err := r.queries.DeleteInbox.Exec(id)
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
	err := r.queries.CountInboxesByAccount.Get(&total, accountID)
	if err != nil {
		return nil, 0, err
	}

	var inboxes []*models.Inbox
	err = r.queries.ListInboxesByAccount.Select(&inboxes, accountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return inboxes, total, nil
}

func (r *repository) CreateRule(rule *models.Rule) error {
	return r.queries.CreateRule.QueryRow(rule.InboxID, rule.Sender, rule.Receiver, rule.Subject).
		Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

func (r *repository) GetRule(id int) (*models.Rule, error) {
	var rule models.Rule
	err := r.queries.GetRule.Get(&rule, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *repository) UpdateRule(rule *models.Rule) error {
	result, err := r.queries.UpdateRule.Exec(rule.Sender, rule.Receiver, rule.Subject, rule.ID)
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
	result, err := r.queries.DeleteRule.Exec(id)
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
	err := r.queries.CountRulesByInbox.Get(&total, inboxID)
	if err != nil {
		return nil, 0, err
	}

	var rules []*models.Rule
	err = r.queries.ListRulesByInbox.Select(&rules, inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) CreateMessage(message *models.Message) error {
	return r.queries.CreateMessage.QueryRow(
		message.InboxID, message.Sender, message.Receiver, message.Subject, message.Body).
		Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)
}

func (r *repository) GetMessage(id int) (*models.Message, error) {
	var message models.Message
	err := r.queries.GetMessage.Get(&message, id)
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
	err := r.queries.CountMessagesByInbox.Get(&total, inboxID)
	if err != nil {
		return nil, 0, err
	}

	var messages []*models.Message
	err = r.queries.ListMessagesByInbox.Select(&messages, inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *repository) ListRules(limit, offset int) ([]*models.Rule, int, error) {
	var total int
	err := r.queries.CountRules.Get(&total)
	if err != nil {
		return nil, 0, err
	}

	var rules []*models.Rule
	err = r.queries.ListRules.Select(&rules, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) GetInboxByEmail(email string) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.queries.GetInboxByEmail.Get(&inbox, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *repository) CreateUser(user *models.User) error {
	return r.queries.CreateUser.QueryRow(
		user.Name,
		user.Username,
		user.Password,
		user.Email,
		user.Status,
		user.Kind,
		user.PasswordLogin).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *repository) GetUser(id int) (*models.User, error) {
	var user models.User
	err := r.queries.GetUser.Get(&user, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(user *models.User) error {
	return r.queries.UpdateUser.QueryRow(
		user.Name,
		user.Username,
		user.Password,
		user.Email,
		user.Status,
		user.Kind,
		user.PasswordLogin,
		user.ID).
		Scan(&user.UpdatedAt)
}

func (r *repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.queries.GetUserByUsername.Get(&user, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) DeleteUser(id int) error {
	result, err := r.queries.DeleteUser.Exec(id)
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
	// Read the initialization query directly from the database connection
	// since it's a one-time operation and doesn't need to be prepared
	queryBytes, err := queriesFS.ReadFile("queries.sql")
	if err != nil {
		return fmt.Errorf("failed to read queries file: %w", err)
	}

	// Get the initialization query
	queries, err := goyesql.ParseBytes(queryBytes)
	if err != nil {
		return fmt.Errorf("failed to parse queries: %w", err)
	}

	// Start a transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute the initialization query
	if _, err := tx.Exec(queries["initialize-tables"].Query); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to execute initialization query: %w, failed to rollback: %v", err, rbErr)
		}
		return fmt.Errorf("failed to execute initialization query: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
