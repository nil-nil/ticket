package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type DNSDomain struct {
	ID     uuid.UUID
	Tenant uuid.UUID
	Name   string
}

type DNSDomainRepository interface {
	GetDomains(ctx context.Context, tenant uuid.UUID) ([]DNSDomain, error)
	CreateDomain(ctx context.Context, tenant uuid.UUID, domain DNSDomain) (DNSDomain, error)
}

type DNSDomainService struct {
	repo        DNSDomainRepository
	eventBus    *EventBus[DNSDomain]
	domainCache *Cache[DNSDomain]
}

func NewDNSDomainService(repo DNSDomainRepository, eventDriver EventBusDriver, cacheDriver CacheDriver) (*DNSDomainService, error) {
	cache, err := NewCache[DNSDomain]("dnsdomains", cacheDriver)
	if err != nil {
		return nil, fmt.Errorf("error creating cache instance: %w", err)
	}
	evt, err := NewEventBus[DNSDomain]("dnsdomains", eventDriver)
	if err != nil {
		return nil, fmt.Errorf("error creating event bus instance: %w", err)
	}
	return &DNSDomainService{repo: repo, eventBus: evt, domainCache: cache}, nil
}

func (s *DNSDomainService) GetDomains(ctx context.Context, tenant uuid.UUID) ([]DNSDomain, error) {
	domains, err := s.repo.GetDomains(ctx, tenant)
	if err != nil {
		return nil, err
	}
	for _, domain := range domains {
		s.domainCache.Set(domain.ID.String(), domain)
	}
	return domains, nil
}

func (s *DNSDomainService) CreateDomain(ctx context.Context, tenant uuid.UUID, name string) (DNSDomain, error) {
	domain, err := s.repo.CreateDomain(ctx, tenant, DNSDomain{ID: uuid.New(), Tenant: tenant, Name: name})
	if err != nil {
		return DNSDomain{}, err
	}
	s.domainCache.Set(domain.ID.String(), domain)

	return domain, nil
}
