package mail

import (
	"context"
	"errors"
	"strings"

	"github.com/nil-nil/ticket/internal/domain"
)

var (
	ErrBlockedSender       = errors.New("sender is blocked")
	ErrInvalidEmailAddress = errors.New("invalid email address")
	ErrAliasNotFound       = errors.New("alias is not found")
)

type MailServerRepository interface {
	IsAuthoritative(domain string) bool
	IsBlocked(address string) bool
}

type AuthFunc func(username, password string) (domain.User, error)

func NewServer(mailServerRepo MailServerRepository, aliasRepo domain.AliasRepository, authFunc AuthFunc) *Server {
	return &Server{
		AuthFunc:   authFunc,
		repository: mailServerRepo,
		aliases:    domain.NewAliasService(aliasRepo),
	}
}

type Server struct {
	AuthFunc   AuthFunc
	repository MailServerRepository
	aliases    *domain.AliasService
}

func (s *Server) ValidateSenderAddress(address string) error {
	if s.repository.IsBlocked(address) {
		return ErrBlockedSender
	}
	return nil
}

func (s *Server) ValidateRecipientAddress(address string) error {
	user, mailDomain, err := getUserAndDomainParts(address)
	if err != nil {
		return ErrInvalidEmailAddress
	}

	if authoritative := s.repository.IsAuthoritative(mailDomain); !authoritative {
		return nil
	}

	_, err = s.aliases.Find(context.Background(), domain.FindAliasParameters{User: &user, Domain: &mailDomain})
	if errors.Is(err, domain.ErrNotFound) {
		return ErrAliasNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func getUserAndDomainParts(address string) (user, domain string, err error) {
	parts := strings.Split(address, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return user, domain, ErrInvalidEmailAddress
	}
	return parts[0], parts[1], nil
}
