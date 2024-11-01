package storage

import (
	"database/sql"
	"errors"
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
	query := `INSERT INTO accounts (name) VALUES ($1) RETURNING id`
	return r.db.QueryRow(query, account.Name).Scan(&account.ID)
}

func (r *repository) GetAccount(id int) (*models.Account, error) {
	var account models.Account
	err := r.db.Get(&account, "SELECT id, name FROM accounts WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *repository) UpdateAccount(account *models.Account) error {
	result, err := r.db.Exec("UPDATE accounts SET name = $1 WHERE id = $2", account.Name, account.ID)
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
		"SELECT id, name FROM accounts ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (r *repository) CreateInbox(inbox *models.Inbox) error {
	query := `INSERT INTO inboxes (account_id, email) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, inbox.AccountID, inbox.Email).Scan(&inbox.ID)
}

func (r *repository) GetInbox(id int) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.db.Get(&inbox, "SELECT id, account_id, email FROM inboxes WHERE id = $1", id)
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
		"SELECT id, account_id, email FROM inboxes WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3",
		accountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return inboxes, total, nil
}

func (r *repository) CreateRule(rule *models.Rule) error {
	query := `INSERT INTO rules (inbox_id, sender, receiver, subject) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(query, rule.InboxID, rule.Sender, rule.Receiver, rule.Subject).Scan(&rule.ID)
}

func (r *repository) GetRule(id int) (*models.Rule, error) {
	var rule models.Rule
	err := r.db.Get(&rule, "SELECT id, inbox_id, sender, receiver, subject FROM rules WHERE id = $1", id)
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
		"SELECT id, inbox_id, sender, receiver, subject FROM rules WHERE inbox_id = $1 ORDER BY id LIMIT $2 OFFSET $3",
		inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) CreateMessage(message *models.Message) error {
	query := `INSERT INTO messages (inbox_id, sender, receiver, subject, body) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, message.InboxID, message.Sender, message.Receiver, message.Subject, message.Body).Scan(&message.ID)
}

func (r *repository) GetMessage(id int) (*models.Message, error) {
	var message models.Message
	err := r.db.Get(&message, "SELECT id, inbox_id, sender, receiver, subject, body FROM messages WHERE id = $1", id)
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
		"SELECT id, inbox_id, sender, receiver, subject, body FROM messages WHERE inbox_id = $1 ORDER BY id LIMIT $2 OFFSET $3",
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
		"SELECT id, inbox_id, sender, receiver, subject FROM rules ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) GetInboxByEmail(email string) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.db.Get(&inbox, "SELECT id, account_id, email FROM inboxes WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *repository) InitializeTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			name TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS inboxes (
			id SERIAL PRIMARY KEY,
			account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
			email TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS rules (
			id SERIAL PRIMARY KEY,
			inbox_id INTEGER REFERENCES inboxes(id) ON DELETE CASCADE,
			sender TEXT,
			receiver TEXT,
			subject TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			inbox_id INTEGER REFERENCES inboxes(id) ON DELETE CASCADE,
			sender TEXT,
			receiver TEXT,
			subject TEXT,
			body TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
			username TEXT
		)`,
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
