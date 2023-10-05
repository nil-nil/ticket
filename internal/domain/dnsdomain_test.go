package domain_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDNSDomain(t *testing.T) {
	repo := &mockDNSDomainRepository{domains: make([]domain.DNSDomain, 512)}
	svc, err := domain.NewDNSDomainService(repo, &mockEventBusDriver{}, mockCache)
	assert.NoError(t, err, "domain.NewDNSDomainService() should not error")

	tenant1 := uuid.New()

	t.Run("TestGetDomains", func(t *testing.T) {
		d1 := domain.DNSDomain{
			ID:     uuid.New(),
			Name:   "test.com",
			Tenant: tenant1,
		}
		d2 := domain.DNSDomain{
			ID:     uuid.New(),
			Name:   "example.com",
			Tenant: tenant1,
		}
		repo.domains = []domain.DNSDomain{
			d1,
			d2,
		}

		t.Run("CorrectTenant", func(t *testing.T) {
			got, err := svc.GetDomains(context.Background(), tenant1)
			assert.NoError(t, err, "DNSDomainService.GetDomains() should not error")
			assert.Equal(t, 2, len(got), "Expected 2 domains")
			assert.Equal(t, d1, mockCache.cache["dnsdomains."+d1.ID.String()], "Expected domain 1 to be cached")
			assert.Equal(t, d2, mockCache.cache["dnsdomains."+d2.ID.String()], "Expected domain 2 to be cached")
			expect := []domain.DNSDomain{d1, d2}
			slices.SortFunc(expect, domainSliceSortFunc)
			slices.SortFunc(got, domainSliceSortFunc)
			assert.Equal(t, expect, got, "Expected got domains to match repo")
		})

		t.Run("IncorrectTenant", func(t *testing.T) {
			got, err := svc.GetDomains(context.Background(), uuid.New())
			assert.NoError(t, err, "DNSDomainService.GetDomains() should not error")
			assert.Equal(t, 0, len(got), "expected empty list for non-matching tenant")
		})
	})

	t.Run("TestCreateDomain", func(t *testing.T) {
		d, err := svc.CreateDomain(context.Background(), tenant1, "foo.com")
		assert.NoError(t, err, "DNSDomainService.CreateDomain() should not error")
		assert.Equal(t, "foo.com", d.Name, "Created domain name should match")
		assert.Equal(t, tenant1, d.Tenant, "Created domain tenant should match")
		assert.Equal(t, d, mockCache.cache[fmt.Sprintf("dnsdomains.%s", d.ID)], "Expected domain to be cached")
	})
}

type mockDNSDomainRepository struct {
	domains []domain.DNSDomain
}

func (m *mockDNSDomainRepository) CreateDomain(ctx context.Context, tenant uuid.UUID, d domain.DNSDomain) (domain.DNSDomain, error) {
	if d.Tenant != tenant {
		return d, errors.New("tenant doesn't match")
	}
	m.domains = append(m.domains, d)

	return d, nil
}

func (m *mockDNSDomainRepository) GetDomains(ctx context.Context, tenant uuid.UUID) ([]domain.DNSDomain, error) {
	domains := make([]domain.DNSDomain, 0, len(m.domains))
	for _, v := range m.domains {
		if v.Tenant == tenant {
			domains = append(domains, v)
		}
	}

	return domains, nil
}

func domainSliceSortFunc(a, b domain.DNSDomain) int {
	if a.ID == b.ID {
		return 0
	}
	if a.ID.String() < b.ID.String() {
		return -1
	}
	return 1
}
