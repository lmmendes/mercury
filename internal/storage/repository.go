package storage

import (
	"database/sql"
	"errors"
	"mercury/internal/models"
)

type Repository interface {
	// Account operations
	CreateAccount(account *models.Account) error
	GetAccount(id int) (*models.Account, error)
	UpdateAccount(account *models.Account) error
	DeleteAccount(id int) error
	ListAccounts() ([]*models.Account, error)

	// Inbox operations
	CreateInbox(inbox *models.Inbox) error
	GetInbox(id int) (*models.Inbox, error)
	UpdateInbox(inbox *models.Inbox) error
	DeleteInbox(id int) error
	ListInboxesByAccount(accountID int) ([]*models.Inbox, error)

	// Rule operations
	CreateRule(rule *models.Rule) error
	GetRule(id int) (*models.Rule, error)
	UpdateRule(rule *models.Rule) error
	DeleteRule(id int) error
	ListRulesByInbox(inboxID int) ([]*models.Rule, error)

	// Message operations
	CreateMessage(message *models.Message) error
	GetMessage(id int) (*models.Message, error)
	ListMessagesByInbox(inboxID int) ([]*models.Message, error)
	ListRules() ([]*models.Rule, error)
	GetInboxByEmail(email string) (*models.Inbox, error)

	// Initialize tables
	InitializeTables() error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAccount(account *models.Account) error {
	result, err := r.db.Exec("INSERT INTO accounts (name) VALUES (?)", account.Name)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	account.ID = int(id)
	return nil
}

func (r *repository) GetAccount(id int) (*models.Account, error) {
	var account models.Account
	err := r.db.QueryRow("SELECT id, name FROM accounts WHERE id = ?", id).Scan(&account.ID, &account.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *repository) UpdateAccount(account *models.Account) error {
	result, err := r.db.Exec("UPDATE accounts SET name = ? WHERE id = ?", account.Name, account.ID)
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
	result, err := r.db.Exec("DELETE FROM accounts WHERE id = ?", id)
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

func (r *repository) ListAccounts() ([]*models.Account, error) {
	rows, err := r.db.Query("SELECT id, name FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		account := &models.Account{}
		if err := rows.Scan(&account.ID, &account.Name); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *repository) CreateInbox(inbox *models.Inbox) error {
	result, err := r.db.Exec("INSERT INTO inboxes (account_id, email) VALUES (?, ?)",
		inbox.AccountID, inbox.Email)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	inbox.ID = int(id)
	return nil
}

func (r *repository) GetInbox(id int) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.db.QueryRow("SELECT id, account_id, email FROM inboxes WHERE id = ?", id).
		Scan(&inbox.ID, &inbox.AccountID, &inbox.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *repository) UpdateInbox(inbox *models.Inbox) error {
	result, err := r.db.Exec("UPDATE inboxes SET email = ? WHERE id = ?",
		inbox.Email, inbox.ID)
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
	result, err := r.db.Exec("DELETE FROM inboxes WHERE id = ?", id)
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

func (r *repository) ListInboxesByAccount(accountID int) ([]*models.Inbox, error) {
	rows, err := r.db.Query("SELECT id, account_id, email FROM inboxes WHERE account_id = ?", accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inboxes []*models.Inbox
	for rows.Next() {
		inbox := &models.Inbox{}
		if err := rows.Scan(&inbox.ID, &inbox.AccountID, &inbox.Email); err != nil {
			return nil, err
		}
		inboxes = append(inboxes, inbox)
	}
	return inboxes, rows.Err()
}

func (r *repository) CreateRule(rule *models.Rule) error {
	result, err := r.db.Exec("INSERT INTO rules (inbox_id, sender, receiver, subject) VALUES (?, ?, ?, ?)",
		rule.InboxID, rule.Sender, rule.Receiver, rule.Subject)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	rule.ID = int(id)
	return nil
}

func (r *repository) GetRule(id int) (*models.Rule, error) {
	var rule models.Rule
	err := r.db.QueryRow("SELECT id, inbox_id, sender, receiver, subject FROM rules WHERE id = ?", id).
		Scan(&rule.ID, &rule.InboxID, &rule.Sender, &rule.Receiver, &rule.Subject)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *repository) UpdateRule(rule *models.Rule) error {
	result, err := r.db.Exec("UPDATE rules SET sender = ?, receiver = ?, subject = ? WHERE id = ?",
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
	result, err := r.db.Exec("DELETE FROM rules WHERE id = ?", id)
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

func (r *repository) ListRulesByInbox(inboxID int) ([]*models.Rule, error) {
	rows, err := r.db.Query("SELECT id, inbox_id, sender, receiver, subject FROM rules WHERE inbox_id = ?", inboxID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*models.Rule
	for rows.Next() {
		rule := &models.Rule{}
		if err := rows.Scan(&rule.ID, &rule.InboxID, &rule.Sender, &rule.Receiver, &rule.Subject); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func (r *repository) CreateMessage(message *models.Message) error {
	result, err := r.db.Exec("INSERT INTO messages (inbox_id, sender, receiver, subject, body) VALUES (?, ?, ?, ?, ?)",
		message.InboxID, message.Sender, message.Receiver, message.Subject, message.Body)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	message.ID = int(id)
	return nil
}

func (r *repository) GetMessage(id int) (*models.Message, error) {
	var message models.Message
	err := r.db.QueryRow("SELECT id, inbox_id, sender, receiver, subject, body FROM messages WHERE id = ?", id).
		Scan(&message.ID, &message.InboxID, &message.Sender, &message.Receiver, &message.Subject, &message.Body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *repository) ListMessagesByInbox(inboxID int) ([]*models.Message, error) {
	rows, err := r.db.Query("SELECT id, inbox_id, sender, receiver, subject, body FROM messages WHERE inbox_id = ?", inboxID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		if err := rows.Scan(&message.ID, &message.InboxID, &message.Sender, &message.Receiver, &message.Subject, &message.Body); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r *repository) ListRules() ([]*models.Rule, error) {
	rows, err := r.db.Query("SELECT id, inbox_id, sender, receiver, subject FROM rules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*models.Rule
	for rows.Next() {
		rule := &models.Rule{}
		if err := rows.Scan(&rule.ID, &rule.InboxID, &rule.Sender, &rule.Receiver, &rule.Subject); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func (r *repository) GetInboxByEmail(email string) (*models.Inbox, error) {
	var inbox models.Inbox
	err := r.db.QueryRow("SELECT id, account_id, email FROM inboxes WHERE email = ?", email).
		Scan(&inbox.ID, &inbox.AccountID, &inbox.Email)
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS inboxes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id INTEGER,
			email TEXT UNIQUE,
			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS rules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			inbox_id INTEGER,
			sender TEXT,
			receiver TEXT,
			subject TEXT,
			FOREIGN KEY (inbox_id) REFERENCES inboxes(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			inbox_id INTEGER,
			sender TEXT,
			receiver TEXT,
			subject TEXT,
			body TEXT,
			FOREIGN KEY (inbox_id) REFERENCES inboxes(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id INTEGER,
			username TEXT NOT NULL UNIQUE,
			password TEXT NULL,
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			status TEXT CHECK( status IN ('enabled', 'disabled') ),
			kind TEXT CHECK( kind IN ('user', 'api') ),
			password_login BOOLEAN NOT NULL DEFAULT false,
			loggedin_at      TIMESTAMP NULL,
			created_at       DEFAULT CURRENT_TIMESTAMP,
			updated_at       DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		if _, err := r.db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
