package domain_test

import (
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

var eventStringTable = []struct {
	event       domain.EventType
	eventString string
	description string
}{
	{
		event:       domain.CreateEvent,
		eventString: "create",
		description: "StringCreateEvent",
	},
	{
		event:       domain.UpdateEvent,
		eventString: "update",
		description: "UpdateEventString",
	},
	{
		event:       domain.DeleteEvent,
		eventString: "delete",
		description: "DeleteEventString",
	},
	{
		event:       domain.UnknownEvent,
		eventString: "unknown",
		description: "UnknownEventString",
	},
}

func TestEventString(t *testing.T) {
	for _, tc := range eventStringTable {
		t.Run(tc.description, func(t *testing.T) {
			v := tc.event.String()
			assert.Equal(t, tc.eventString, v)
		})
	}
}

func TestEventParse(t *testing.T) {
	for _, tc := range eventStringTable {
		t.Run(tc.description, func(t *testing.T) {
			event := domain.ParseEventString(tc.eventString)
			assert.Equal(t, event, tc.event)
		})
	}

	t.Run("ParseInvalidEventString", func(t *testing.T) {
		event := domain.ParseEventString("this Event Doesn't Exist")
		assert.Equal(t, event, domain.UnknownEvent)
	})
}

func TestPublishing(t *testing.T) {
	t.Run("NoPrefix", func(t *testing.T) {
		m := mockEventBusDriver{}
		eventBus, err := domain.NewEventBus[int]("", &m)
		assert.EqualError(t, err, domain.ErrEventPrefixInvalid.Error(), "NewEventBus() should give meaningful error when no prefix given")
		assert.Nil(t, eventBus, "no event bus should be returned on err")
	})

	t.Run("PublishMissingID", func(t *testing.T) {
		m := mockEventBusDriver{}
		eventBus, err := domain.NewEventBus[struct{ ID uint64 }]("structs", &m)
		assert.NoError(t, err, "NewEventBus() should not error when a prefix is given")

		s := struct{ ID uint64 }{ID: 1}
		err = eventBus.Publish("", domain.CreateEvent, s)
		assert.EqualError(t, err, domain.ErrEventKeyInvalid.Error(), "should give meaningful error when no id is provided in publishing")
		assert.Nil(t, m.EventSubject, "no event should be published on err")
	})

	t.Run("PublishSuccess", func(t *testing.T) {
		m := mockEventBusDriver{}
		eventBus, err := domain.NewEventBus[domain.User]("users", &m)
		assert.NoError(t, err, "NewEventBus() should not error when a prefix is given")

		u := domain.User{ID: 1}
		err = eventBus.Publish("1", domain.CreateEvent, u)
		assert.NoError(t, err, "valid publish shouldn't error")
		assert.Equal(t, *m.EventSubject, "users:1:create", "expected event matching subject")
		assert.Equal(t, u, m.EventData, "expected given event data")
	})
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
