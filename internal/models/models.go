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
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type Inbox struct {
	Base
	AccountID int    `json:"account_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type User struct {
	Base
	Name           string    `json:"name"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	Email          string    `json:"email"`
	Status         string    `json:"status"`
	Kind           string    `json:"kind"`
	password_login bool      `json:"password_login"`
	LoggedinAt     null.Time `json:"loggedin_at"`
}

type Rule struct {
	Base
	InboxID  int    `json:"inbox_id" validate:"required"`
	Sender   string `json:"sender" validate:"omitempty,email"`
	Receiver string `json:"receiver" validate:"omitempty,email"`
	Subject  string `json:"subject" validate:"omitempty,max=200"`
}

type Message struct {
	Base
	InboxID  int    `json:"inbox_id" validate:"required"`
	Sender   string `json:"sender" validate:"required,email"`
	Receiver string `json:"receiver" validate:"required,email"`
	Subject  string `json:"subject" validate:"required,max=200"`
	Body     string `json:"body" validate:"required"`
}
