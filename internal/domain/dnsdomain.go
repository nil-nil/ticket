package domain

import (
	"context"
	"fmt"
)

type DNSDomain struct {
	ID   uint64
	Name string
}

type DNSDomainRepository interface {
	GetDomains(context.Context) ([]DNSDomain, error)
	CreateDomain(ctx context.Context, domain DNSDomain) (DNSDomain, error)
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

func (s *DNSDomainService) GetDomains(ctx context.Context) ([]DNSDomain, error) {
	domains, err := s.repo.GetDomains(ctx)
	if err != nil {
		return nil, err
	}
	for _, domain := range domains {
		s.domainCache.Set(fmt.Sprint(domain.ID), domain)
	}
	return domains, nil
}

func (s *DNSDomainService) CreateDomain(ctx context.Context, name string) (DNSDomain, error) {
	domain, err := s.repo.CreateDomain(ctx, DNSDomain{Name: name})
	if err != nil {
		return DNSDomain{}, err
	}
	s.domainCache.Set(fmt.Sprint(domain.ID), domain)

	return domain, nil
}
