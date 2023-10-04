package domain_test

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestGetAlias(t *testing.T) {
	tenant1 := uuid.New()
	tenant2 := uuid.New()
	var repo = mockAliasRepo{
		aliases: map[string]domain.Alias{
			"test@test.com":      {ID: 1, Tenant: tenant1, User: "test", Domain: "test.com"},
			"sample@example.com": {ID: 2, Tenant: tenant1, User: "sample", Domain: "example.com"},
		},
	}

	svc := domain.NewAliasService(&repo)

	t.Run("NotFoundAtAll", func(t *testing.T) {
		alias, err := svc.Find(context.Background(), tenant1, domain.FindAliasParameters{User: ptr.To("bob"), Domain: ptr.To("sample.com")})
		assert.Equal(t, domain.Alias{}, alias, "alias should be empty")
		assert.EqualError(t, domain.ErrNotFound, err.Error(), "expected not found error")
	})

	t.Run("Success", func(t *testing.T) {
		alias, err := svc.Find(context.Background(), tenant1, domain.FindAliasParameters{User: ptr.To("test"), Domain: ptr.To("test.com")})
		assert.Equal(t, domain.Alias{ID: 1, Tenant: tenant1, User: "test", Domain: "test.com"}, alias, "alias should not be empty")
		assert.NoError(t, err, "error should be nil")
	})

	t.Run("WrongTenant", func(t *testing.T) {
		alias, err := svc.Find(context.Background(), tenant2, domain.FindAliasParameters{User: ptr.To("test"), Domain: ptr.To("test.com")})
		assert.Equal(t, domain.Alias{}, alias, "alias should be empty")
		assert.EqualError(t, domain.ErrNotFound, err.Error(), "expected not found error")
	})
}

func TestCreateAlias(t *testing.T) {
	tenant1 := uuid.New()
	var repo = mockAliasRepo{
		aliases: map[string]domain.Alias{
			"test@test.com":      {ID: 1, User: "test", Domain: "test.com"},
			"sample@example.com": {ID: 2, User: "sample", Domain: "example.com"},
		},
	}

	svc := domain.NewAliasService(&repo)

	alias, err := svc.Create(context.Background(), tenant1, "bob", "sample.com")
	assert.Equal(t, domain.Alias{ID: 3, Tenant: tenant1, User: "bob", Domain: "sample.com"}, alias, "alias should not be empty")
	assert.NoError(t, err, "error should be nil")
}

func TestDeleteAlias(t *testing.T) {
	tenant1 := uuid.New()
	var repo = mockAliasRepo{
		aliases: map[string]domain.Alias{
			"test@test.com":      {ID: 1, Tenant: tenant1, User: "test", Domain: "test.com"},
			"sample@example.com": {ID: 2, Tenant: tenant1, User: "sample", Domain: "example.com"},
		},
	}

	svc := domain.NewAliasService(&repo)

	alias, err := svc.Delete(context.Background(), tenant1, 2)
	assert.Equal(t, repo.aliases["sample@example.com"], alias, "alias should not be empty")
	assert.NotNil(t, alias.DeletedAt, "DeletedAt should be set now")
	assert.NoError(t, err, "error should be nil")
}

type mockAliasRepo struct {
	aliases map[string]domain.Alias
}

func (m *mockAliasRepo) Find(ctx context.Context, tenant uuid.UUID, params domain.FindAliasParameters) (domain.Alias, error) {
	if params.User == nil || params.Domain == nil {
		return domain.Alias{}, domain.ErrNotFound
	}
	email := fmt.Sprintf("%s@%s", *params.User, *params.Domain)

	alias, ok := m.aliases[email]
	if !ok || alias.Tenant != tenant {
		return domain.Alias{}, domain.ErrNotFound
	}

	return alias, nil
}

func (m *mockAliasRepo) Create(ctx context.Context, tenant uuid.UUID, user string, mailDomain string) (domain.Alias, error) {
	email := fmt.Sprintf("%s@%s", user, mailDomain)
	m.aliases[email] = domain.Alias{Tenant: tenant, Domain: mailDomain, User: user, ID: m.getNextId()}
	return m.aliases[email], nil
}

func (m *mockAliasRepo) Delete(ctx context.Context, tenant uuid.UUID, ID uint64) (domain.Alias, error) {
	var alias domain.Alias
	for k := range m.aliases {
		if m.aliases[k].ID == ID && m.aliases[k].Tenant == tenant {
			alias = m.aliases[k]
			now := time.Now()
			alias.DeletedAt = &now
			m.aliases[k] = alias
		}
	}
	if alias == (domain.Alias{}) {
		return alias, domain.ErrNotFound
	}
	return alias, nil
}

func (m *mockAliasRepo) getNextId() uint64 {
	if len(m.aliases) == 0 {
		return 1
	}
	aliases := make([]domain.Alias, 0, len(m.aliases))
	for _, alias := range m.aliases {
		aliases = append(aliases, alias)
	}
	lastAlias := slices.MaxFunc(aliases, func(a domain.Alias, b domain.Alias) int {
		if a.ID > b.ID {
			return 1
		}
		return 0
	})
	return lastAlias.ID + 1
}

func TestGetNextId(t *testing.T) {
	t.Run("EmptySlice", func(t *testing.T) {
		repo := mockAliasRepo{aliases: map[string]domain.Alias{}}
		nextId := repo.getNextId()
		assert.Equal(t, uint64(1), nextId)
	})

	t.Run("ExistingElements", func(t *testing.T) {
		repo := mockAliasRepo{
			aliases: map[string]domain.Alias{
				"test@test.com":      {ID: 1, User: "test", Domain: "test.com"},
				"sample@example.com": {ID: 5, User: "sample", Domain: "example.com"},
			},
		}
		nextId := repo.getNextId()
		assert.Equal(t, uint64(6), nextId)
	})
}
