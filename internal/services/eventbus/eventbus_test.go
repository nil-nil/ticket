package eventbus_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/services/eventbus"
	"github.com/stretchr/testify/assert"
)

func TestPublishing(t *testing.T) {
	t.Run("NotAStruct", func(t *testing.T) {
		m := mockEventBusDriver[int]{}
		eventBus := eventbus.NewEventBus(&m)

		var i int = 0
		err := eventBus.Publish(i, domain.CreateEvent)
		assert.EqualError(t, err, eventbus.ErrNotAStruct.Error(), "should give meaningful error when subject is not a struct")
		assert.Nil(t, m.Event, "no event should be published on err")
	})

	t.Run("NoIDField", func(t *testing.T) {
		m := mockEventBusDriver[struct{ ID uint64 }]{}
		eventBus := eventbus.NewEventBus(&m)

		s := struct{ ID uint64 }{ID: 1}
		err := eventBus.Publish(s, domain.CreateEvent)
		assert.EqualError(t, err, eventbus.ErrNoIDField.Error(), "should give meaningful error when subject has no id tag")
		assert.Nil(t, m.Event, "no event should be published on err")
	})

	t.Run("ValidStruct", func(t *testing.T) {
		m := mockEventBusDriver[domain.User]{}
		eventBus := eventbus.NewEventBus(&m)

		testUUID := uuid.New()
		u := domain.User{ID: testUUID}
		err := eventBus.Publish(u, domain.CreateEvent)
		assert.NoError(t, err, "valid struct shouldn't error")
		assert.Equal(t, *m.Event, "domain.User:"+testUUID.String()+":create", "expected event matching subject")
	})
}

func TestSubscribing(t *testing.T) {
	t.Run("NotAStruct", func(t *testing.T) {
		m := mockEventBusDriver[int]{}
		eventBus := eventbus.NewEventBus(&m)

		var i int = 0
		err := eventBus.Subscribe(i, false, []domain.EventType{domain.CreateEvent}, func(data int, eventType domain.EventType) {})
		assert.EqualError(t, err, eventbus.ErrNotAStruct.Error(), "should give meaningful error when subject is not a struct")
		assert.Nil(t, m.SubscriptionKey, "no subscription should be made on error")
		assert.Nil(t, m.SubscriptionCallback, "no subscription should be made on error")
	})

	t.Run("NoIDField", func(t *testing.T) {
		m := mockEventBusDriver[struct{ ID uint64 }]{}
		eventBus := eventbus.NewEventBus(&m)

		s := struct{ ID uint64 }{ID: 1}
		err := eventBus.Subscribe(s, false, []domain.EventType{domain.CreateEvent}, func(data struct{ ID uint64 }, eventType domain.EventType) {})
		assert.EqualError(t, err, eventbus.ErrNoIDField.Error(), "should give meaningful error when subject has no id tag")
		assert.Nil(t, m.SubscriptionKey, "no subscription should be made on error")
		assert.Nil(t, m.SubscriptionCallback, "no subscription should be made on error")
	})

	t.Run("ValidStructNonWildcard", func(t *testing.T) {
		m := mockEventBusDriver[domain.User]{}
		eventBus := eventbus.NewEventBus(&m)

		testUUID := uuid.New()
		u := domain.User{ID: testUUID}
		var gotUser domain.User
		var gotEventType domain.EventType
		err := eventBus.Subscribe(u, false, []domain.EventType{domain.CreateEvent}, func(data domain.User, eventType domain.EventType) {
			t.Log("running callback")
			gotUser = data
			gotEventType = eventType
		})
		assert.NoError(t, err, "valid struct shouldn't error")
		assert.Equal(t, "domain.User:"+testUUID.String()+":create", *m.SubscriptionKey, "key should be generated properly and sent to driver")
		assert.NotNil(t, m.SubscriptionCallback, "callback should be sent to driver")

		(*m.SubscriptionCallback)("domain.User:"+testUUID.String()+":create", u)
		assert.Equal(t, u, gotUser, "expected matching data")
		assert.Equal(t, domain.CreateEvent, gotEventType, "expected parsed event type")
	})

	t.Run("ValidStructWithWildcard", func(t *testing.T) {
		m := mockEventBusDriver[domain.User]{}
		eventBus := eventbus.NewEventBus(&m)

		testUUID := uuid.New()
		u := domain.User{ID: testUUID}
		err := eventBus.Subscribe(u, true, []domain.EventType{domain.CreateEvent}, func(data domain.User, eventType domain.EventType) {})
		assert.NoError(t, err, "valid struct shouldn't error")
		assert.Equal(t, "domain.User:*:create", *m.SubscriptionKey, "key should be generated properly and sent to driver")
		assert.NotNil(t, m.SubscriptionCallback, "callback should be sent to driver")
	})
}

type mockEventBusDriver[T any] struct {
	Event                *string
	Data                 *T
	SubscriptionKey      *string
	SubscriptionCallback *func(eventKey string, data T)
}

func (m *mockEventBusDriver[T]) Publish(subject string, data T) error {
	m.Event = &subject
	m.Data = &data
	return nil
}

func (m *mockEventBusDriver[T]) Subscribe(subject string, callback func(eventKey string, data T)) error {
	m.SubscriptionKey = &subject
	m.SubscriptionCallback = &callback
	return nil
}

func (m *mockEventBusDriver[T]) Reset() {
	m.Event = nil
	m.SubscriptionKey = nil
	m.Data = nil
	m.SubscriptionCallback = nil
}
