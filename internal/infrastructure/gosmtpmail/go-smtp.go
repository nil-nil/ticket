package gosmtpmail

import (
	"io"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/services/mail"
)

func NewServer(mailServerRepo mail.MailServerRepository, cacheDriver domain.CacheDriver, eventBusDriver domain.EventBusDriver, authFunc mail.AuthFunc) *smtp.Server {
	mailServer := mail.NewServer(mailServerRepo, cacheDriver, eventBusDriver, authFunc)
	be := backend{server: mailServer}
	server := smtp.NewServer(&be)
	server.Addr = ":25"
	server.Domain = "localhost"
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	return server
}

type backend struct {
	server *mail.Server
}

func (b *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{server: b.server}, nil
}

type session struct {
	server *mail.Server
	user   *domain.User
	from   string
	to     []string
}

func (s *session) AuthPlain(username, password string) error {
	if s.server.AuthFunc == nil {
		return nil
	}

	user, err := s.server.AuthFunc(username, password)
	if err != nil {
		return err
	}

	s.user = &user
	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	err := s.server.ValidateSenderAddress(from)
	if err != nil {
		return err
	}

	s.from = from
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	err := s.server.ValidateRecipientAddress(to)
	if err != nil {
		return err
	}

	s.to = append(s.to, to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	// b, err := io.ReadAll(r)
	// if err != nil {
	// 	return err
	// }
	// if err = HandleMessage(&b); err != nil {
	// 	return err
	// }
	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
