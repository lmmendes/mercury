package storage

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql/v2"
	goyesqlx "github.com/knadh/goyesql/v2/sqlx"
)

//go:embed queries.sql schema.sql
var queriesFS embed.FS

type Queries struct {
	// Account queries
	CreateAccount *sqlx.Stmt `query:"create-account"`
	GetAccount    *sqlx.Stmt `query:"get-account"`
	UpdateAccount *sqlx.Stmt `query:"update-account"`
	DeleteAccount *sqlx.Stmt `query:"delete-account"`
	ListAccounts  *sqlx.Stmt `query:"list-accounts"`
	CountAccounts *sqlx.Stmt `query:"count-accounts"`

	// Inbox queries
	CreateInbox           *sqlx.Stmt `query:"create-inbox"`
	GetInbox              *sqlx.Stmt `query:"get-inbox"`
	UpdateInbox           *sqlx.Stmt `query:"update-inbox"`
	DeleteInbox           *sqlx.Stmt `query:"delete-inbox"`
	ListInboxesByAccount  *sqlx.Stmt `query:"list-inboxes-by-account"`
	CountInboxesByAccount *sqlx.Stmt `query:"count-inboxes-by-account"`
	GetInboxByEmail       *sqlx.Stmt `query:"get-inbox-by-email"`

	// Rule queries
	CreateRule        *sqlx.Stmt `query:"create-rule"`
	GetRule           *sqlx.Stmt `query:"get-rule"`
	UpdateRule        *sqlx.Stmt `query:"update-rule"`
	DeleteRule        *sqlx.Stmt `query:"delete-rule"`
	ListRulesByInbox  *sqlx.Stmt `query:"list-rules-by-inbox"`
	CountRulesByInbox *sqlx.Stmt `query:"count-rules-by-inbox"`
	ListRules         *sqlx.Stmt `query:"list-rules"`
	CountRules        *sqlx.Stmt `query:"count-rules"`

	// Message queries
	CreateMessage        *sqlx.Stmt `query:"create-message"`
	GetMessage           *sqlx.Stmt `query:"get-message"`
	ListMessagesByInbox  *sqlx.Stmt `query:"list-messages-by-inbox"`
	CountMessagesByInbox *sqlx.Stmt `query:"count-messages-by-inbox"`

	// User queries
	CreateUser        *sqlx.Stmt `query:"create-user"`
	GetUser           *sqlx.Stmt `query:"get-user"`
	UpdateUser        *sqlx.Stmt `query:"update-user"`
	DeleteUser        *sqlx.Stmt `query:"delete-user"`
	GetUserByUsername *sqlx.Stmt `query:"get-user-by-username"`
}

func PrepareQueries(db *sqlx.DB) (*Queries, error) {
	// Read queries from embedded file
	queryBytes, err := queriesFS.ReadFile("queries.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read queries file: %w", err)
	}

	// Parse queries
	queries, err := goyesql.ParseBytes(queryBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse queries: %w", err)
	}

	// Prepare statements
	var q Queries
	if err := goyesqlx.ScanToStruct(&q, queries, db); err != nil {
		return nil, fmt.Errorf("failed to prepare queries: %w", err)
	}

	return &q, nil
}
