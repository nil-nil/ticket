package email

import (
	"context"
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestIsAuthoritative(t *testing.T) {
	repo := &mockMailServerRepository{
		authoritativeDomains: []string{"example.com"},
	}
	svc, err := NewMailServerService(repo, mockCache, &mockEventBusDriver{})
	assert.NoError(t, err, "NewMailServerService shoudln't error")

	t.Run("AuthoritativeDomain", func(t *testing.T) {
		result := svc.IsAuthoritative("example.com")
		assert.True(t, result, "existing domain should be authoritative")
		assert.Equal(t, repo.authoritativeDomains, *svc.domainCache, "should be cached now")
	})

	t.Run("NonAuthoritativeDomain", func(t *testing.T) {
		result := svc.IsAuthoritative("test.com")
		assert.False(t, result, "non existing domain should not be authoritative")
	})
}

func TestGetAlias(t *testing.T) {
	repo := &mockMailServerRepository{
		aliases: []domain.Alias{{Domain: "example.com", User: "test", ID: 1}},
	}
	svc, err := NewMailServerService(repo, mockCache, &mockEventBusDriver{})
	assert.NoError(t, err, "NewMailServerService shoudln't error")

	t.Run("ExistingAlias", func(t *testing.T) {
		alias, err := svc.GetAlias(context.Background(), "test", "example.com")
		assert.NoError(t, err)
		assert.Equal(t, domain.Alias{Domain: "example.com", User: "test", ID: 1}, alias)
		assert.Equal(t, repo.aliases, mockCache.cache["mailaliases.example.com"], "should be cached now")
	})

	t.Run("NotExistingAlias", func(t *testing.T) {
		alias, err := svc.GetAlias(context.Background(), "notexist", "notexist.com")
		assert.EqualError(t, err, ErrAliasNotFound.Error())
		assert.Equal(t, domain.Alias{}, alias, "expect zero value alias when erroring")
	})
}

func TestObserver(t *testing.T) {
	repo := &mockMailServerRepository{
		aliases: []domain.Alias{{Domain: "test.com", User: "test", ID: 1}},
	}
	svc, err := NewMailServerService(repo, mockCache, &mockEventBusDriver{})
	assert.NoError(t, err, "NewMailServerService shoudln't error")

	svc.ObserveAliasEvents(domain.CreateEvent, domain.Alias{ID: 2, User: "test2", Domain: "test.com"})
	assert.Equal(t, []domain.Alias{{Domain: "test.com", User: "test", ID: 1}}, mockCache.cache["mailaliases.test.com"], "cache should be refreshed")
}

type mockMailServerRepository struct {
	authoritativeDomains []string
	aliases              []domain.Alias
}

func (m *mockMailServerRepository) GetAuthoritativeDomains(ctx context.Context) ([]string, error) {
	return m.authoritativeDomains, nil
}

func (m *mockMailServerRepository) GetAliases(ctx context.Context, mailDomain *string) ([]domain.Alias, error) {
	if mailDomain != nil {
		matches := make([]domain.Alias, 0)
		for _, alias := range m.aliases {
			if alias.Domain == *mailDomain {
				matches = append(matches, alias)
			}
		}
		return matches, nil
	}
	return m.aliases, nil
}

type mockCacheDriver struct {
	cache map[string]interface{}
}

func (m *mockCacheDriver) Get(key string) (interface{}, error) {
	val, ok := m.cache[key]
	if !ok {
		return nil, domain.ErrNotFoundInCache
	}
	return val, nil
}

func (m *mockCacheDriver) Set(key string, value interface{}) error {
	m.cache[key] = value
	return nil
}

func (m *mockCacheDriver) Forget(key string) error {
	delete(m.cache, key)
	return nil
}

var mockCache = &mockCacheDriver{
	cache: make(map[string]interface{}, 0),
}

type mockEventBusDriver struct {
	EventSubject         *string
	EventData            interface{}
	SubscriptionKey      *string
	SubscriptionCallback *func(eventKey string, data interface{})
}

func (m *mockEventBusDriver) Publish(subject string, data interface{}) error {
	m.EventSubject = &subject
	m.EventData = data
	return nil
}

func (m *mockEventBusDriver) Subscribe(subject string, callback func(eventKey string, data interface{})) error {
	m.SubscriptionKey = &subject
	m.SubscriptionCallback = &callback
	return nil
}

func (m *mockEventBusDriver) Reset() {
	m.EventSubject = nil
	m.EventData = nil
	m.SubscriptionKey = nil
	m.SubscriptionCallback = nil
}
