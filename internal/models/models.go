package models

import (
	null "github.com/volatiletech/null/v9"
)

type Base struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt null.Time `json:"created_at" db:"created_at"`
	UpdatedAt null.Time `json:"updated_at" db:"updated_at"`
}

type Project struct {
	Base
	Name string `json:"name" db:"name" validate:"required,min=2,max=100"`
}

type Inbox struct {
	Base
	ProjectID int    `json:"project_id" db:"project_id" validate:"required"`
	Email     string `json:"email" db:"email" validate:"required,email"`
}

type User struct {
	Base
	Name          string    `json:"name" db:"name"`
	Username      string    `json:"username" db:"username"`
	Password      string    `json:"password" db:"password"`
	Email         string    `json:"email" db:"email"`
	Status        string    `json:"status" db:"status"`
	Role          string    `json:"role" db:"role"`
	PasswordLogin bool      `json:"password_login" db:"password_login"`
	LoggedinAt    null.Time `json:"loggedin_at" db:"loggedin_at"`
}

type ProjectUser struct {
	Base
	ProjectID int    `json:"project_id" db:"project_id" validate:"required"`
	UserID    int    `json:"user_id" db:"user_id" validate:"required"`
	Role      string `json:"role" db:"role" validate:"required"`
}

type UserToken struct {
	Base
	UserID int    `json:"user_id" db:"user_id" validate:"required"`
	Token  string `json:"token" db:"token" validate:"required"`
}

type ForwardRule struct {
	Base
	InboxID  int    `json:"inbox_id" db:"inbox_id" validate:"required"`
	Sender   string `json:"sender" db:"sender" validate:"omitempty,email"`
	Receiver string `json:"receiver" db:"receiver" validate:"omitempty,email"`
	Subject  string `json:"subject" db:"subject" validate:"omitempty,max=200"`
}

type Message struct {
	Base
	InboxID  int    `json:"inbox_id" db:"inbox_id" validate:"required"`
	Sender   string `json:"sender" db:"sender" validate:"required,email"`
	Receiver string `json:"receiver" db:"receiver" validate:"required,email"`
	Subject  string `json:"subject" db:"subject" validate:"required,max=200"`
	Body     string `json:"body" db:"body" validate:"required"`
}
