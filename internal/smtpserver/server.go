package smtpserver

import (
	"bytes"
	"fmt"
	"io"
	"mercury/internal/core"
	"mercury/internal/models"

	"github.com/emersion/go-message"
	"github.com/emersion/go-smtp"
)

type Server struct {
	core *core.Core
	smtp *smtp.Server
}

func NewServer(core *core.Core) *Server {
	be := &Backend{core: core}
	s := smtp.NewServer(be)

	s.Addr = core.Config.SMTPPort
	s.Domain = "localhost"
	s.AllowInsecureAuth = true

	return &Server{
		core: core,
		smtp: s,
	}
}

func (s *Server) ListenAndServe() error {
	return s.smtp.ListenAndServe()
}

type Backend struct {
	core *core.Core
}

func (be *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{core: be.core}, nil
}

type Session struct {
	core *core.Core
	from string
	to   string
}

func (s *Session) Mail(from string, _ *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, _ *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	// Parse the email to get the subject
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	// Parse email headers and body
	msg, err := message.Read(bytes.NewReader(buf.Bytes()))
	if err != nil {
		s.core.Logger.Error("Failed to parse email: %v", err)
		return err
	}

	header := msg.Header
	body := new(bytes.Buffer)
	if _, err := io.Copy(body, msg.Body); err != nil {
		s.core.Logger.Error("Failed to read message body: %v", err)
		return err
	}

	message := &models.Message{
		Body:     body.String(),
		Sender:   s.from,
		Receiver: s.to,
		Subject:  header.Get("Subject"),
		InboxID:  0, // We need to look up the inbox ID based on the recipient email
	}

	// Look up the inbox ID based on the recipient email
	inbox, err := s.core.Repository.GetInboxByEmail(s.to)
	if err != nil {
		s.core.Logger.Error("Failed to find inbox for email %s: %v", s.to, err)
		return err
	}
	if inbox == nil {
		s.core.Logger.Error("No inbox found for email %s", s.to)
		return fmt.Errorf("no inbox found for recipient")
	}

	message.InboxID = inbox.ID

	s.core.Logger.Info("Received email from %s to %s", s.from, s.to)

	if err := s.core.StoreMessage(message); err != nil {
		s.core.Logger.Error("Failed to store message: %v", err)
		return err
	}

	s.core.Logger.Debug("Stored message successfully")
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func (s *Session) AuthPlain(username, password string) error {
	return nil // TODO: For now, accept all auth
}
