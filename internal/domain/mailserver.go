package domain

import (
	"context"
)

type domainList []string
type aliasList []string

func NewMailServerService(repo MailServerRepository, cacheDriver cacheDriver, eventBusDriver eventBusDriver) (*MailServerService, error) {
	domainCache, err := NewCache[domainList]("maildomains", cacheDriver)
	if err != nil {
		return nil, err
	}
	aliasCache, err := NewCache[aliasList]("mailaliases", cacheDriver)
	if err != nil {
		return nil, err
	}
	aliasEventBus, err := NewEventBus[Alias]("aliases", eventBusDriver)
	if err != nil {
		return nil, err
	}
	svc := &MailServerService{
		repo:          repo,
		domainCache:   domainCache,
		aliasCache:    aliasCache,
		aliasEventBus: aliasEventBus,
	}
	aliasEventBus.Subscribe(nil, []EventType{CreateEvent, UpdateEvent, DeleteEvent}, svc.ObserveAliasEvents)
	return svc, nil
}

type MailServerService struct {
	repo          MailServerRepository
	domainCache   *Cache[domainList]
	aliasCache    *Cache[aliasList]
	aliasEventBus *EventBus[Alias]
}

func (s *MailServerService) ObserveAliasEvents(eventType EventType, data Alias) {
	aliasList, err := s.repo.GetAliases(context.Background(), &data.Domain)
	if err != nil {
		return
	}
	s.aliasCache.Set(data.Domain, aliasList)
}

type MailServerRepository interface {
	GetAliases(ctx context.Context, domain *string) (aliasList, error)
	GetAuthoritativeDomains(ctx context.Context) (domainList, error)
}
