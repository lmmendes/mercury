package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/emersion/go-smtp"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

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

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables()

	r := mux.NewRouter()
	r.HandleFunc("/accounts", createAccount).Methods("POST")
	r.HandleFunc("/accounts", getAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}", getAccount).Methods("GET")
	r.HandleFunc("/accounts/{id}", updateAccount).Methods("PUT")
	r.HandleFunc("/accounts/{id}", deleteAccount).Methods("DELETE")

	r.HandleFunc("/accounts/{accountId}/inboxes", createInbox).Methods("POST")
	r.HandleFunc("/accounts/{accountId}/inboxes", getInboxes).Methods("GET")
	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}", getInbox).Methods("GET")
	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}", updateInbox).Methods("PUT")
	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}", deleteInbox).Methods("DELETE")

	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules", createRule).Methods("POST")
	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules", getRules).Methods("GET")
	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules/{ruleId}", getRule).Methods("GET")
	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules/{ruleId}", updateRule).Methods("PUT")
	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules/{ruleId}", deleteRule).Methods("DELETE")

	r.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/messages", getMessages).Methods("GET")

	go func() {
		log.Println("Starting HTTP server at :8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	be := &Backend{}
	s := smtp.NewServer(be)

	s.Addr = ":1025"
	s.Domain = "localhost"
	s.AllowInsecureAuth = true

	log.Println("Starting SMTP server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS accounts (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)`,
		`CREATE TABLE IF NOT EXISTS inboxes (id INTEGER PRIMARY KEY AUTOINCREMENT, account_id INTEGER, email TEXT)`,
		`CREATE TABLE IF NOT EXISTS rules (id INTEGER PRIMARY KEY AUTOINCREMENT, inbox_id INTEGER, sender TEXT, receiver TEXT, subject TEXT)`,
		`CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY AUTOINCREMENT, inbox_id INTEGER, sender TEXT, receiver TEXT, subject TEXT, body TEXT)`,
		`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, account_id INTEGER, username TEXT)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO accounts (name) VALUES (?)", account.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM accounts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		accounts = append(accounts, account)
	}

	json.NewEncoder(w).Encode(accounts)
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var account Account
	if err := db.QueryRow("SELECT id, name FROM accounts WHERE id = ?", id).Scan(&account.ID, &account.Name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE accounts SET name = ? WHERE id = ?", account.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.Exec("DELETE FROM accounts WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createInbox(w http.ResponseWriter, r *http.Request) {
	accountId := mux.Vars(r)["accountId"]
	var inbox Inbox
	if err := json.NewDecoder(r.Body).Decode(&inbox); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inbox.AccountID, _ = strconv.Atoi(accountId)
	result, err := db.Exec("INSERT INTO inboxes (account_id, email) VALUES (?, ?)", inbox.AccountID, inbox.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inbox.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inbox)
}

func getInboxes(w http.ResponseWriter, r *http.Request) {
	accountId := mux.Vars(r)["accountId"]
	rows, err := db.Query("SELECT id, account_id, email FROM inboxes WHERE account_id = ?", accountId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var inboxes []Inbox
	for rows.Next() {
		var inbox Inbox
		if err := rows.Scan(&inbox.ID, &inbox.AccountID, &inbox.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		inboxes = append(inboxes, inbox)
	}

	json.NewEncoder(w).Encode(inboxes)
}

func getInbox(w http.ResponseWriter, r *http.Request) {
	inboxId := mux.Vars(r)["inboxId"]
	var inbox Inbox
	if err := db.QueryRow("SELECT id, account_id, email FROM inboxes WHERE id = ?", inboxId).Scan(&inbox.ID, &inbox.AccountID, &inbox.Email); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(inbox)
}

func updateInbox(w http.ResponseWriter, r *http.Request) {
	inboxId := mux.Vars(r)["inboxId"]
	var inbox Inbox
	if err := json.NewDecoder(r.Body).Decode(&inbox); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE inboxes SET email = ? WHERE id = ?", inbox.Email, inboxId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteInbox(w http.ResponseWriter, r *http.Request) {
	inboxId := mux.Vars(r)["inboxId"]
	_, err := db.Exec("DELETE FROM inboxes WHERE id = ?", inboxId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createRule(w http.ResponseWriter, r *http.Request) {
	inboxId := mux.Vars(r)["inboxId"]
	var rule Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rule.InboxID, _ = strconv.Atoi(inboxId)
	result, err := db.Exec("INSERT INTO rules (inbox_id, sender, receiver, subject) VALUES (?, ?, ?, ?)", rule.InboxID, rule.Sender, rule.Receiver, rule.Subject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rule.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

func getRules(w http.ResponseWriter, r *http.Request) {
	inboxId := mux.Vars(r)["inboxId"]
	rows, err := db.Query("SELECT id, inbox_id, sender, receiver, subject FROM rules WHERE inbox_id = ?", inboxId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var rule Rule
		if err := rows.Scan(&rule.ID, &rule.InboxID, &rule.Sender, &rule.Receiver, &rule.Subject); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rules = append(rules, rule)
	}

	json.NewEncoder(w).Encode(rules)
}

func getRule(w http.ResponseWriter, r *http.Request) {
	ruleId := mux.Vars(r)["ruleId"]
	var rule Rule
	if err := db.QueryRow("SELECT id, inbox_id, sender, receiver, subject FROM rules WHERE id = ?", ruleId).Scan(&rule.ID, &rule.InboxID, &rule.Sender, &rule.Receiver, &rule.Subject); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(rule)
}

func updateRule(w http.ResponseWriter, r *http.Request) {
	ruleId := mux.Vars(r)["ruleId"]
	var rule Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE rules SET sender = ?, receiver = ?, subject = ? WHERE id = ?", rule.Sender, rule.Receiver, rule.Subject, ruleId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteRule(w http.ResponseWriter, r *http.Request) {
	ruleId := mux.Vars(r)["ruleId"]
	_, err := db.Exec("DELETE FROM rules WHERE id = ?", ruleId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	inboxId := mux.Vars(r)["inboxId"]
	rows, err := db.Query("SELECT id, inbox_id, sender, receiver, subject, body FROM messages WHERE inbox_id = ?", inboxId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.InboxID, &message.Sender, &message.Receiver, &message.Subject, &message.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	json.NewEncoder(w).Encode(messages)
}

type Backend struct{}

func (be *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

type Session struct {
	from string
	to   string
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	var message Message
	message.Body = readAll(r)
	message.Sender = s.from
	message.Receiver = s.to

	// Match against rules and store in appropriate inbox
	storeMessage(message)

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func readAll(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}

func storeMessage(message Message) {
	rows, err := db.Query("SELECT id, inbox_id, sender, receiver, subject FROM rules")
	if err != nil {
		log.Println("Error querying rules:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rule Rule
		if err := rows.Scan(&rule.ID, &rule.InboxID, &rule.Sender, &rule.Receiver, &rule.Subject); err != nil {
			log.Println("Error scanning rule:", err)
			return
		}

		if rule.Sender == message.Sender && rule.Receiver == message.Receiver && rule.Subject == message.Subject {
			_, err := db.Exec("INSERT INTO messages (inbox_id, sender, receiver, subject, body) VALUES (?, ?, ?, ?, ?)", rule.InboxID, message.Sender, message.Receiver, message.Subject, message.Body)
			if err != nil {
				log.Println("Error inserting message:", err)
			}
			return
		}
	}

	// If no rule matches, store in the inbox with the matching email address
	var inboxID int
	err = db.QueryRow("SELECT id FROM inboxes WHERE email = ?", message.Receiver).Scan(&inboxID)
	if err != nil {
		log.Println("Error finding inbox:", err)
		return
	}

	_, err = db.Exec("INSERT INTO messages (inbox_id, sender, receiver, subject, body) VALUES (?, ?, ?, ?, ?)", inboxID, message.Sender, message.Receiver, message.Subject, message.Body)
	if err != nil {
		log.Println("Error inserting message:", err)
	}
}
