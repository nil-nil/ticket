package email

import (
	"context"
	"errors"
	"io"
	"net/mail"
	"strings"

	"github.com/nil-nil/ticket/internal/domain"
)

var (
	ErrBlockedSender       = errors.New("sender is blocked")
	ErrInvalidEmailAddress = errors.New("invalid email address")
	ErrAliasNotFound       = errors.New("alias is not found")
)

type AuthFunc func(username, password string) (domain.User, error)

func NewServer(mailServerRepo MailServerRepository, cacheDriver domain.CacheDriver, eventBusDriver domain.EventBusDriver, authFunc AuthFunc) *Server {
	svc, _ := NewMailServerService(mailServerRepo, cacheDriver, eventBusDriver)
	return &Server{
		AuthFunc:    authFunc,
		mailService: svc,
	}
}

type Server struct {
	AuthFunc    AuthFunc
	mailService *MailServerService
}

func (s *Server) ValidateSenderAddress(address string) error {
	return nil
}

func (s *Server) ValidateRecipientAddress(address string) error {
	user, mailDomain, err := getUserAndDomainParts(address)
	if err != nil {
		return ErrInvalidEmailAddress
	}

	if authoritative := s.mailService.IsAuthoritative(mailDomain); !authoritative {
		return nil
	}

	alias, err := s.mailService.GetAlias(context.Background(), user, mailDomain)
	if errors.Is(err, domain.ErrNotFound) {
		return ErrAliasNotFound
	} else if err != nil {
		return err
	}
	if alias.DeletedAt != nil {
		return ErrAliasNotFound
	}

	return nil
}

func (s *Server) ReceiveData(reader io.Reader) error {
	msg, err := mail.ReadMessage(reader)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("invalid message")
	}

	_, err = s.mailService.CreateEmail(context.Background(), *msg)

	return err
}

func getUserAndDomainParts(address string) (user, domain string, err error) {
	parts := strings.Split(address, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return user, domain, ErrInvalidEmailAddress
	}
	return parts[0], parts[1], nil
}
