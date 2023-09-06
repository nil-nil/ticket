package eventbus_test

import (
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/services/eventbus"
	"github.com/stretchr/testify/assert"
)

type mockEventBusDriver struct {
	Event                *string
	SubscriptionKey      *string
	SubscriptionCallback *func(subject string)
}

func (m *mockEventBusDriver) Publish(subject string) error {
	m.Event = &subject
	return nil
}

func (m *mockEventBusDriver) Subscribe(subject string, callback func(subject string)) error {
	m.SubscriptionKey = &subject
	m.SubscriptionCallback = &callback
	return nil
}

func (m *mockEventBusDriver) Reset() {
	m.Event = nil
	m.SubscriptionKey = nil
	m.SubscriptionCallback = nil
}

func TestPublishing(t *testing.T) {
	m := mockEventBusDriver{}
	eventBus := eventbus.NewEventBus(&m)

	t.Run("not a struct", func(t *testing.T) {
		m.Reset()
		var i int = 0
		err := eventBus.Publish(i, domain.CreateEvent)
		assert.EqualError(t, err, eventbus.ErrNotAStruct.Error(), "should give meaningful error when subject is not a struct")
		assert.Nil(t, m.Event, "no event should be published on err")
	})

	t.Run("no ID field", func(t *testing.T) {
		m.Reset()
		s := struct{ ID uint64 }{ID: 1}
		err := eventBus.Publish(s, domain.CreateEvent)
		assert.EqualError(t, err, eventbus.ErrNoIDField.Error(), "should give meaningful error when subject has no id tag")
		assert.Nil(t, m.Event, "no event should be published on err")
	})

	t.Run("valid struct", func(t *testing.T) {
		m.Reset()
		u := domain.User{ID: 1}
		err := eventBus.Publish(u, domain.CreateEvent)
		assert.NoError(t, err, "valid struct shouldn't error")
		assert.Equal(t, *m.Event, "domain.User:1:create", "expected event matching subject")
	})
}

func TestSubscribing(t *testing.T) {
	m := mockEventBusDriver{}
	eventBus := eventbus.NewEventBus(&m)

	t.Run("not a struct", func(t *testing.T) {
		m.Reset()
		var i int = 0
		err := eventBus.Subscribe(i, false, []domain.EventType{domain.CreateEvent}, func(subjectType, subjectId string, eventType domain.EventType) {})
		assert.EqualError(t, err, eventbus.ErrNotAStruct.Error(), "should give meaningful error when subject is not a struct")
		assert.Nil(t, m.SubscriptionKey, "no subscription should be made on error")
		assert.Nil(t, m.SubscriptionCallback, "no subscription should be made on error")
	})

	t.Run("no ID field", func(t *testing.T) {
		m.Reset()
		s := struct{ ID uint64 }{ID: 1}
		err := eventBus.Subscribe(s, false, []domain.EventType{domain.CreateEvent}, func(subjectType, subjectId string, eventType domain.EventType) {})
		assert.EqualError(t, err, eventbus.ErrNoIDField.Error(), "should give meaningful error when subject has no id tag")
		assert.Nil(t, m.SubscriptionKey, "no subscription should be made on error")
		assert.Nil(t, m.SubscriptionCallback, "no subscription should be made on error")
	})

	t.Run("valid struct non wildcard", func(t *testing.T) {
		m.Reset()
		u := domain.User{ID: 1}
		var gotSubjectType, gotSubjectId string
		var gotEventType domain.EventType
		err := eventBus.Subscribe(u, false, []domain.EventType{domain.CreateEvent}, func(subjectType, subjectId string, eventType domain.EventType) {
			gotSubjectType = subjectType
			gotSubjectId = subjectId
			gotEventType = eventType
		})
		assert.NoError(t, err, "valid struct shouldn't error")
		assert.Equal(t, "domain.User:1:create", *m.SubscriptionKey, "key should be generated properly and sent to driver")
		assert.NotNil(t, m.SubscriptionCallback, "callback should be sent to driver")

		(*m.SubscriptionCallback)("domain.User:5:delete")
		assert.Equal(t, gotSubjectType, "domain.User", "expected parsed subject type")
		assert.Equal(t, gotSubjectId, "5", "expected parsed subject id")
		assert.Equal(t, gotEventType, domain.DeleteEvent, "expected parsed event type")
	})

	t.Run("valid struct non wildcard", func(t *testing.T) {
		m.Reset()
		u := domain.User{ID: 1}
		err := eventBus.Subscribe(u, true, []domain.EventType{domain.CreateEvent}, func(subjectType, subjectId string, eventType domain.EventType) {})
		assert.NoError(t, err, "valid struct shouldn't error")
		assert.Equal(t, "domain.User:*:create", *m.SubscriptionKey, "key should be generated properly and sent to driver")
		assert.NotNil(t, m.SubscriptionCallback, "callback should be sent to driver")
	})
}
