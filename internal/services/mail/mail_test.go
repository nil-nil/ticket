package mail

import (
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestValidateSenderAddress(t *testing.T) {
	repo := &mockMailServerRepository{
		authoritativeDomains: []string{"example.com"},
	}

	server := NewServer(repo, mockCache, &mockEventBusDriver{}, func(username, password string) (domain.User, error) { return domain.User{}, nil })

	t.Run("test valid sender", func(t *testing.T) {
		err := server.ValidateSenderAddress("alan@example.com")
		assert.NoError(t, err)
	})
}

func TestValidateRecipientAddress(t *testing.T) {
	repo := &mockMailServerRepository{
		authoritativeDomains: []string{"test.com"},
		aliases: []domain.Alias{
			{User: "test", Domain: "test.com", ID: 1},
		},
	}

	server := NewServer(repo, mockCache, &mockEventBusDriver{}, func(username, password string) (domain.User, error) { return domain.User{}, nil })

	table := []struct {
		description string
		email       string
		expectErr   error
	}{
		{description: "valid non-authoritative recipient", email: "alan@example.com", expectErr: nil},
		{description: "valid authoritative recipient", email: "test@test.com", expectErr: nil},
		{description: "invalid authoritative recipient", email: "fail@test.com", expectErr: ErrAliasNotFound},
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
