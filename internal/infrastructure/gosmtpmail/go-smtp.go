package gosmtpmail

import (
	"errors"
	"io"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/services/email"
)

var (
	ErrMailboxNotFound = &smtp.SMTPError{
		Code:         550,
		EnhancedCode: smtp.EnhancedCode{5, 1, 1},
		Message:      "The email account that you tried to reach does not exist. Please try double-checking the recipient's email address for typos or unnecessary spaces.",
	}
)

func NewServer(mailServerRepo email.MailServerRepository, cacheDriver domain.CacheDriver, eventBusDriver domain.EventBusDriver, authFunc email.AuthFunc) *smtp.Server {
	mailServer := email.NewServer(mailServerRepo, cacheDriver, eventBusDriver, authFunc)
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
	server *email.Server
}

func (b *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{server: b.server}, nil
}

type session struct {
	server *email.Server
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
	if errors.Is(err, email.ErrAliasNotFound) {
		return ErrMailboxNotFound
	}
	if err != nil {
		return err
	}

	s.to = append(s.to, to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	return s.server.ReceiveData(r)
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
