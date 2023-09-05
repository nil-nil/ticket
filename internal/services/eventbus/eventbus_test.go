package eventbus_test

import (
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/services/eventbus"
	"github.com/stretchr/testify/assert"
)

type mockEventBusDriver struct {
	Event *string
}

func (m *mockEventBusDriver) Publish(subject string) error {
	m.Event = &subject
	return nil
}

func (m *mockEventBusDriver) Reset() {
	m.Event = nil
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
		assert.Equal(t, *m.Event, "domain.User:1", "expected event matching subject")
	})
}
