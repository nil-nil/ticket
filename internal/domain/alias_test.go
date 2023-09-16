package domain_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestAliasEmail(t *testing.T) {
	table := []struct {
		description string
		user        string
		domain      string
		expect      string
	}{
		{description: "User email", user: "bob", domain: "test.com", expect: "bob@test.com"},
	}

	for _, tc := range table {
		t.Run(tc.description, func(t *testing.T) {
			alias := domain.Alias{User: tc.user, Domain: tc.domain}
			email := alias.GetEmail()
			assert.Equal(t, tc.expect, email)
		})
	}
}

type mockAliasRepo struct {
	aliases map[string]domain.Alias
}

func (m *mockAliasRepo) Find(ctx context.Context, params domain.FindAliasParameters) (domain.Alias, error) {
	if params.User == nil || params.Domain == nil {
		return domain.Alias{}, domain.ErrNotFound
	}
	email := fmt.Sprintf("%s@%s", *params.User, *params.Domain)

	alias, ok := m.aliases[email]
	if !ok {
		return domain.Alias{}, domain.ErrNotFound
	}

	return alias, nil
}

func TestGetAlias(t *testing.T) {
	var repo = mockAliasRepo{
		aliases: map[string]domain.Alias{
			"test@test.com":      {ID: 1, User: "test", Domain: "test.com"},
			"sample@example.com": {ID: 1, User: "sample", Domain: "example.com"},
		},
	}

	svc := domain.NewAliasService(&repo)

	alias, err := svc.Find(context.Background(), domain.FindAliasParameters{User: ptr.To("bob"), Domain: ptr.To("sample.com")})
	assert.Equal(t, domain.Alias{}, alias, "alias should be empty")
	assert.EqualError(t, domain.ErrNotFound, err.Error(), "expected not found error")

	alias, err = svc.Find(context.Background(), domain.FindAliasParameters{User: ptr.To("test"), Domain: ptr.To("test.com")})
	assert.Equal(t, domain.Alias{ID: 1, User: "test", Domain: "test.com"}, alias, "alias should not be empty")
	assert.NoError(t, err, "error should be nil")
}
