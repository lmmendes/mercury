package models

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Inbox struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Email     string `json:"email"`
}

type Rule struct {
	ID       int    `json:"id"`
	InboxID  int    `json:"inbox_id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Subject  string `json:"subject"`
}

type Message struct {
	ID       int    `json:"id"`
	InboxID  int    `json:"inbox_id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

type User struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Username  string `json:"username"`
}
