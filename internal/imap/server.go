package imap

import (
	"context"
	"errors"
	"time"

	"inbox451/internal/core"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/server"
)

// ImapBackend implements go-imap/backend interface
type ImapBackend struct {
	core *core.Core
}

// Login handles user authentication
func (be *ImapBackend) Login(connInfo *imap.ConnInfo, username string, password string) (backend.User, error) {
	be.core.Logger.Info("Login attempt starting")
	be.core.Logger.Info("Login username=[%s] password=[%s]", username, password)
	if username == "test@example.com" && password == "password" {
		return &ImapUser{username: username, core: be.core}, nil
	}
	return nil, errors.New("invalid credentials")
}

// ImapUser implements go-imap/backend.User interface
type ImapUser struct {
	core     *core.Core
	username string
}

// Username returns user's email
func (u *ImapUser) Username() string {
	return u.username
}

// ListMailboxes returns a list of mailboxes for the user
func (u *ImapUser) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	// Return list of available mailboxes (folders)
	return []backend.Mailbox{
		&ImapMailbox{name: "INBOX", user: u},
	}, nil
}

// GetMailbox returns a specific mailbox
func (u *ImapUser) GetMailbox(name string) (backend.Mailbox, error) {
	if name == "INBOX" {
		return &ImapMailbox{name: name, user: u}, nil
	}
	return nil, errors.New("mailbox not found")
}

func (u *ImapUser) CreateMailbox(name string) error {
	return errors.New("mailbox creation not supported")
}

func (u *ImapUser) DeleteMailbox(name string) error {
	return errors.New("mailbox deletion not supported")
}

func (u *ImapUser) RenameMailbox(existingName, newName string) error {
	return errors.New("mailbox renaming not supported")
}

// ImapMailbox implements go-imap/backend.Mailbox interface
type ImapMailbox struct {
	name     string
	user     *ImapUser
	messages []*imap.Message
}

// Name returns mailbox name
func (m *ImapMailbox) Name() string {
	return m.name
}

// Info returns mailbox info
func (m *ImapMailbox) Info() (*imap.MailboxInfo, error) {
	info := &imap.MailboxInfo{
		Attributes: []string{},
		Delimiter:  "/",
		Name:       m.name,
	}
	return info, nil
}

func (u *ImapUser) Logout() error {
	return nil
}

// Status returns mailbox status
func (m *ImapMailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	status := imap.NewMailboxStatus(m.name, items)
	status.Messages = uint32(len(m.messages))
	status.Unseen = 0
	return status, nil
}

// ListMessages returns a list of messages
func (m *ImapMailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)
	for _, msg := range m.messages {
		if seqSet.Contains(msg.SeqNum) {
			ch <- msg
		}
	}
	return nil
}

// SearchMessages searches for messages matching the given criteria
func (m *ImapMailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	// Implement search logic
	return []uint32{}, nil
}

func (m *ImapMailbox) Check() error {
	return nil
}

func (m *ImapMailbox) ExpungeMessages(uids []uint32) error {
	return errors.New("expunge not supported")
}

func (m *ImapMailbox) CopyMessages(uid bool, seqSet *imap.SeqSet, dest string) error {
	return errors.New("copy not supported")
}

func (m *ImapMailbox) MoveMessages(uid bool, seqSet *imap.SeqSet, dest string) error {
	return errors.New("move not supported")
}

func (m *ImapMailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return errors.New("message creation not supported")
}

func (m *ImapMailbox) UpdateMessagesFlags(uid bool, seqSet *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	return errors.New("flag updates not supported")
}

func (m *ImapMailbox) Expunge() error {
	return errors.New("expunge not supported")
}

func (m *ImapMailbox) SetSubscribed(subscribed bool) error {
	return errors.New("subscription changes not supported")
}

type ImapServer struct {
	core *core.Core
	imap *server.Server
}

func (s *ImapServer) ListenAndServe() error {
	return s.imap.ListenAndServe()
}

// Add Shutdown method to ImapServer struct
func (s *ImapServer) Shutdown(ctx context.Context) error {
	return s.imap.Close()
}

func NewServer(core *core.Core) *ImapServer {
	core.Logger.Info("IMAP Server initializing")

	be := &ImapBackend{core: core}
	s := server.New(be)
	s.Addr = ":1143"

	s.Debug = core.Logger.Writer()

	// Allow unencrypted plain text authentication
	s.AllowInsecureAuth = true

	return &ImapServer{
		core: core,
		imap: s,
	}
}
