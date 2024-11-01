package smtp

import (
	"bytes"
	"io"
	"mercury/internal/core"
	"mercury/internal/models"

	"github.com/emersion/go-smtp"
)

type SmtpServer struct {
	core *core.Core
	smtp *smtp.Server
}

type SmtpBackend struct {
	core *core.Core
}

type SmtpSession struct {
	core *core.Core
	from string
	to   string
}

func NewServer(core *core.Core) *SmtpServer {
	be := &SmtpBackend{core: core}
	s := smtp.NewServer(be)

	s.Addr = core.Config.SMTPPort
	s.Domain = "localhost"
	s.AllowInsecureAuth = true

	return &SmtpServer{
		core: core,
		smtp: s,
	}
}

func (s *SmtpServer) ListenAndServe() error {
	return s.smtp.ListenAndServe()
}

func (be *SmtpBackend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &SmtpSession{core: be.core}, nil
}

func (s *SmtpSession) Mail(from string, _ *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *SmtpSession) Rcpt(to string, _ *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *SmtpSession) Data(r io.Reader) error {
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

func (s *SmtpSession) Reset() {}

func (s *SmtpSession) Logout() error {
	return nil
}

func (s *SmtpSession) AuthPlain(username, password string) error {
	return nil // TODO: For now, accept all auth
}

func readAll(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}
