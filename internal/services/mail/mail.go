package mail

import (
	"errors"
	"strings"

	"github.com/nil-nil/ticket/domain"
)

var (
	ErrBlockedSender       = errors.New("sender is blocked")
	ErrInvalidEmailAddress = errors.New("invalid email address")
)

type server struct {
	AuthFunc       AuthFunc
	blockedSenders map[string]struct{}
}

type AuthFunc func(username, password string) (domain.User, error)

func (s *server) ValidateSenderAddress(address string) error {
	if _, exists := s.blockedSenders[address]; exists {
		return ErrBlockedSender
	}
	return nil
}

func (s *server) ValidateRecipientAddress(address string) error {
	return nil
}

func getUserAndDomainParts(address string) (user, domain string, err error) {
	parts := strings.Split(address, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return user, domain, ErrInvalidEmailAddress
	}
	return parts[0], parts[1], nil
}
