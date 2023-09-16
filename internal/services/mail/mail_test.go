package mail_test

import (
	"context"
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/services/mail"
	"github.com/stretchr/testify/assert"
)

type mockMailServerRepository struct {
	authoritativeDomains map[string]struct{}
	blockedSenders       map[string]struct{}
}

func (m *mockMailServerRepository) IsAuthoritative(domain string) bool {
	_, ok := m.authoritativeDomains[domain]
	return ok
}

func (m *mockMailServerRepository) IsBlocked(address string) bool {
	_, ok := m.blockedSenders[address]
	return ok
}

type mockAliasRepository struct {
	aliases []domain.Alias
}

func (m mockAliasRepository) Find(ctx context.Context, params domain.FindAliasParameters) (domain.Alias, error) {
	if params.User == nil && params.Domain == nil {
		return domain.Alias{}, domain.ErrNotFound
	}
	for _, alias := range m.aliases {
		if alias.User == *params.User && alias.Domain == *params.Domain {
			return alias, nil
		}
	}
	return domain.Alias{}, domain.ErrNotFound
}

func TestValidateSenderAddress(t *testing.T) {
	repo := &mockMailServerRepository{
		blockedSenders: map[string]struct{}{"bob@example.com": {}},
	}
	server := mail.NewServer(repo, nil, nil)

	t.Run("test valid sender", func(t *testing.T) {
		err := server.ValidateSenderAddress("alan@example.com")
		assert.NoError(t, err)
	})

	t.Run("test blocked sender", func(t *testing.T) {
		err := server.ValidateSenderAddress("bob@example.com")
		assert.EqualError(t, err, mail.ErrBlockedSender.Error())
	})
}

func TestValidateRecipientAddress(t *testing.T) {
	mailServerRepo := &mockMailServerRepository{
		blockedSenders:       map[string]struct{}{"bob@example.com": {}},
		authoritativeDomains: map[string]struct{}{"test.com": {}},
	}
	aliasRepo := &mockAliasRepository{
		aliases: []domain.Alias{
			{ID: 1, User: "test", Domain: "test.com"},
		},
	}
	server := mail.NewServer(mailServerRepo, aliasRepo, nil)

	table := []struct {
		description string
		email       string
		expectErr   error
	}{
		{description: "valid non-authoritative recipient", email: "alan@example.com", expectErr: nil},
		{description: "valid authoritative recipient", email: "test@test.com", expectErr: nil},
		{description: "invalid authoritative recipient", email: "fail@test.com", expectErr: mail.ErrAliasNotFound},
	}

	for _, tc := range table {
		t.Run(tc.description, func(t *testing.T) {
			err := server.ValidateRecipientAddress(tc.email)
			if tc.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectErr.Error())
			}
		})
	}
}
