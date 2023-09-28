package domain_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDNSDomain(t *testing.T) {
	repo := &mockDNSDomainRepository{domains: make(map[uint64]domain.DNSDomain, 512)}
	svc, err := domain.NewDNSDomainService(repo, &mockEventBusDriver{}, mockCache)
	assert.NoError(t, err, "domain.NewDNSDomainService() should not error")

	t.Run("TestGetDomains", func(t *testing.T) {
		d1 := domain.DNSDomain{
			ID:   1,
			Name: "test.com",
		}
		d2 := domain.DNSDomain{
			ID:   2,
			Name: "example.com",
		}
		repo.domains = map[uint64]domain.DNSDomain{
			1: d1,
			2: d2,
		}

		domains, err := svc.GetDomains(context.Background())
		assert.NoError(t, err, "DNSDomainService.GetDomains() should not error")
		assert.Equal(t, 2, len(domains), "Expected 2 domains")
		assert.Equal(t, d1, mockCache.cache["dnsdomains.1"], "Expected domain 1 to be cached")
		assert.Equal(t, d2, mockCache.cache["dnsdomains.2"], "Expected domain 2 to be cached")
		assert.Equal(t, []domain.DNSDomain{d1, d2}, domains, "Expected got domains to match repo")
	})

	t.Run("TestCreateDomain", func(t *testing.T) {
		d, err := svc.CreateDomain(context.Background(), "foo.com")
		assert.NoError(t, err, "DNSDomainService.CreateDOmain() should not error")
		assert.Equal(t, "foo.com", d.Name, "Created domain name should match")
		assert.Equal(t, d, mockCache.cache[fmt.Sprintf("dnsdomains.%d", d.ID)], "Expected domain to be cached")
	})
}

type mockDNSDomainRepository struct {
	domains map[uint64]domain.DNSDomain
}

func (m *mockDNSDomainRepository) CreateDomain(ctx context.Context, d domain.DNSDomain) (domain.DNSDomain, error) {
	d.ID = nextMapKey(m.domains)
	m.domains[d.ID] = d

	return d, nil
}

func (m *mockDNSDomainRepository) GetDomains(ctx context.Context) ([]domain.DNSDomain, error) {
	domains := make([]domain.DNSDomain, 0, len(m.domains))
	for _, v := range m.domains {
		domains = append(domains, v)
	}

	return domains, nil
}
