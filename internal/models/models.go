package models

import (
	null "github.com/volatiletech/null/v9"
)

type Base struct {
	ID        int       `json:"id"`
	CreatedAt null.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
}

type Account struct {
	Base

	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Inbox struct {
	Base

	AccountID int    `json:"account_id"`
	Email     string `json:"email"`
}

type User struct {
	Base

	AccountID      int    `json:"account_id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	Status         string `json:"status"`
	Kind           string `json:"kind"`
	password_login bool   `json:"password_login"`
}

type Rule struct {
	Base

	InboxID  int    `json:"inbox_id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Subject  string `json:"subject"`
}

type Message struct {
	Base

	InboxID  int    `json:"inbox_id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}
