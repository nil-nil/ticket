package email

import (
	"context"
	"net/mail"
	"slices"

	"github.com/nil-nil/ticket/internal/domain"
)

func NewMailServerService(repo MailServerRepository, cacheDriver domain.CacheDriver, eventBusDriver domain.EventBusDriver) (*MailServerService, error) {
	aliasCache, err := domain.NewCache[[]domain.Alias]("mailaliases", cacheDriver)
	if err != nil {
		return nil, err
	}
	aliasEventBus, err := domain.NewEventBus[domain.Alias]("aliases", eventBusDriver)
	if err != nil {
		return nil, err
	}
	svc := &MailServerService{
		repo:          repo,
		aliasCache:    aliasCache,
		aliasEventBus: aliasEventBus,
	}
	aliasEventBus.Subscribe(nil, []domain.EventType{domain.CreateEvent, domain.UpdateEvent, domain.DeleteEvent}, svc.ObserveAliasEvents)
	return svc, nil
}

type MailServerService struct {
	repo          MailServerRepository
	domainCache   *[]string
	aliasCache    *domain.Cache[[]domain.Alias]
	aliasEventBus *domain.EventBus[domain.Alias]
}

func (s *MailServerService) ObserveAliasEvents(eventType domain.EventType, data domain.Alias) {
	aliases, err := s.repo.GetAliases(context.Background(), &data.Domain)
	if err != nil {
		return
	}
	s.aliasCache.Set(data.Domain, aliases)
}

func (s *MailServerService) IsAuthoritative(domain string) bool {
	if s.domainCache == nil {
		domains, err := s.repo.GetAuthoritativeDomains(context.Background())
		if err != nil {
			return false
		}
		s.domainCache = &domains
	}
	return slices.Contains[[]string](*s.domainCache, domain)
}

func (s *MailServerService) GetAlias(ctx context.Context, user string, mailDomain string) (domain.Alias, error) {
	aliasList, err := s.aliasCache.Get(mailDomain)
	if err != nil {
		aliasList, err = s.repo.GetAliases(ctx, &mailDomain)
		if err != nil {
			return domain.Alias{}, err
		}
		s.aliasCache.Set(mailDomain, aliasList)
	}
	for _, alias := range aliasList {
		if alias.Domain == mailDomain && alias.User == user {
			return alias, nil
		}
	}
	return domain.Alias{}, ErrAliasNotFound
}

func (s *MailServerService) CreateEmail(ctx context.Context, msg mail.Message) (domain.Email, error) {
	return domain.CreateEmail(ctx, s.repo, msg)
}

type MailServerRepository interface {
	GetAliases(ctx context.Context, domain *string) ([]domain.Alias, error)
	GetAuthoritativeDomains(ctx context.Context) ([]string, error)
	domain.CreateEmailRepository
}
