package smtpserver

import (
	"bytes"
	"io"
	"mercury/internal/core"
	"mercury/internal/models"

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
	message := &models.Message{
		Body:     readAll(r),
		Sender:   s.from,
		Receiver: s.to,
	}

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
	return nil // For now, accept all auth
}

func readAll(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}
