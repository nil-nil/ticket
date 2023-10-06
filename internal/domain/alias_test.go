package domain_test

import (
	"context"
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
		aliases: []domain.Alias{
			{ID: uuid.New(), Tenant: tenant1, User: "test", Domain: "test.com"},
			{ID: uuid.New(), Tenant: tenant1, User: "sample", Domain: "example.com"},
		},
	}

	svc := domain.NewAliasService(&repo)

	t.Run("NotFoundAtAll", func(t *testing.T) {
		alias, err := svc.Find(context.Background(), tenant1, domain.FindAliasParameters{User: ptr.To("bob"), Domain: ptr.To("sample.com")})
		assert.Equal(t, domain.Alias{}, alias, "alias should be empty")
		assert.EqualError(t, domain.ErrNotFound, err.Error(), "expected not found error")
	})

	t.Run("Success", func(t *testing.T) {
		expect := domain.Alias{ID: uuid.New(), Tenant: tenant1, User: "qux", Domain: "foo.com"}
		repo.aliases = append(repo.aliases, expect)
		alias, err := svc.Find(context.Background(), tenant1, domain.FindAliasParameters{User: ptr.To("qux"), Domain: ptr.To("foo.com")})
		assert.Equal(t, expect, alias, "alias should not be empty")
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
		aliases: []domain.Alias{
			{ID: uuid.New(), User: "test", Domain: "test.com"},
			{ID: uuid.New(), User: "sample", Domain: "example.com"},
		},
	}

	svc := domain.NewAliasService(&repo)

	alias, err := svc.Create(context.Background(), tenant1, "bob", "sample.com")
	assert.Equal(t, tenant1, alias.Tenant, "alias tenant should match")
	assert.Equal(t, "bob", alias.User, "alias user should match")
	assert.Equal(t, "sample.com", alias.Domain, "alias domain should match")
	assert.NoError(t, err, "error should be nil")
}

func TestDeleteAlias(t *testing.T) {
	tenant1 := uuid.New()
	aliasId := uuid.New()
	var repo = mockAliasRepo{
		aliases: []domain.Alias{
			{ID: aliasId, Tenant: tenant1, User: "test", Domain: "test.com"},
		},
	}

	svc := domain.NewAliasService(&repo)

	alias, err := svc.Delete(context.Background(), tenant1, aliasId)
	assert.Equal(t, repo.aliases[0], alias, "alias should not be empty")
	assert.NotNil(t, alias.DeletedAt, "DeletedAt should be set now")
	assert.NoError(t, err, "error should be nil")
}

type mockAliasRepo struct {
	aliases []domain.Alias
}

func (m *mockAliasRepo) Find(ctx context.Context, tenant uuid.UUID, params domain.FindAliasParameters) (domain.Alias, error) {
	if params.User == nil || params.Domain == nil {
		return domain.Alias{}, domain.ErrNotFound
	}
	for _, alias := range m.aliases {
		if alias.Domain == *params.Domain && alias.User == *params.User && alias.Tenant == tenant {
			return alias, nil
		}
	}

	return domain.Alias{}, domain.ErrNotFound
}

func (m *mockAliasRepo) Create(ctx context.Context, alias domain.Alias) error {
	m.aliases = append(m.aliases, alias)
	return nil
}

func (m *mockAliasRepo) Delete(ctx context.Context, tenant uuid.UUID, ID uuid.UUID) (domain.Alias, error) {
	var alias domain.Alias
	for idx, alias := range m.aliases {
		if alias.ID == ID && alias.Tenant == tenant {
			alias.DeletedAt = ptr.To(time.Now())
			m.aliases[idx] = alias
			return alias, nil
		}
	}
	return alias, nil
}
