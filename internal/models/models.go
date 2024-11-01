package models

type Account struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" validate:"required,min=2,max=100"`
}

type Inbox struct {
	ID        int    `json:"id" db:"id"`
	AccountID int    `json:"account_id" db:"account_id" validate:"required"`
	Email     string `json:"email" db:"email" validate:"required,email"`
}

type Rule struct {
	ID       int    `json:"id" db:"id"`
	InboxID  int    `json:"inbox_id" db:"inbox_id" validate:"required"`
	Sender   string `json:"sender" db:"sender" validate:"omitempty,email"`
	Receiver string `json:"receiver" db:"receiver" validate:"omitempty,email"`
	Subject  string `json:"subject" db:"subject" validate:"omitempty,max=200"`
}

type Message struct {
	ID       int    `json:"id" db:"id"`
	InboxID  int    `json:"inbox_id" db:"inbox_id" validate:"required"`
	Sender   string `json:"sender" db:"sender" validate:"required,email"`
	Receiver string `json:"receiver" db:"receiver" validate:"required,email"`
	Subject  string `json:"subject" db:"subject" validate:"required,max=200"`
	Body     string `json:"body" db:"body" validate:"required"`
}

type User struct {
	ID        int    `json:"id" db:"id"`
	AccountID int    `json:"account_id" db:"account_id"`
	Username  string `json:"username" db:"username"`
}
