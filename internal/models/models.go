package models

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type Inbox struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type Rule struct {
	ID       int    `json:"id"`
	InboxID  int    `json:"inbox_id" validate:"required"`
	Sender   string `json:"sender" validate:"omitempty,email"`
	Receiver string `json:"receiver" validate:"omitempty,email"`
	Subject  string `json:"subject" validate:"omitempty,max=200"`
}

type Message struct {
	ID       int    `json:"id"`
	InboxID  int    `json:"inbox_id" validate:"required"`
	Sender   string `json:"sender" validate:"required,email"`
	Receiver string `json:"receiver" validate:"required,email"`
	Subject  string `json:"subject" validate:"required,max=200"`
	Body     string `json:"body" validate:"required"`
}

type User struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Username  string `json:"username"`
}
