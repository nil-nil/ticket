package hubeventbus_test

import (
	"sync"
	"testing"
	"time"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/infrastructure/hubeventbus"
	"github.com/stretchr/testify/assert"
)

type testEvent struct {
	eventKey string
	data     interface{}
}

func TestHub(t *testing.T) {
	var wg sync.WaitGroup

	bus, _ := hubeventbus.NewHubEventBusDriver()

	var (
		t1, t2, t3 *testEvent
	)
	bus.Subscribe("test.1.create", func(eventKey string, data interface{}) {
		defer wg.Done()
		t1 = &testEvent{eventKey: eventKey, data: data}
	})

	bus.Subscribe("test.2.*", func(eventKey string, data interface{}) {
		defer wg.Done()
		t2 = &testEvent{eventKey: eventKey, data: data}
	})

	bus.Subscribe("test.*.create", func(eventKey string, data interface{}) {
		defer wg.Done()
		t3 = &testEvent{eventKey: eventKey, data: data}
	})

	// Give the subs time to start
	time.Sleep(500000)

	t.Run("test.1.create", func(t *testing.T) {
		wg.Add(2)
		t1, t2, t3 = nil, nil, nil
		bus.Publish("test.1.create", domain.User{ID: 1, FirstName: "Test", LastName: "Abc"})
		wg.Wait()

		assert.Nil(t, t2, "non-matching sub 2 should be nil")
		assert.Equal(t, testEvent{eventKey: "test.1.create", data: domain.User{ID: 1, FirstName: "Test", LastName: "Abc"}}, *t1, "sub 1 should have received matching event")
		assert.Equal(t, testEvent{eventKey: "test.1.create", data: domain.User{ID: 1, FirstName: "Test", LastName: "Abc"}}, *t3, "sub 3 should have received wildcard matching event")
	})

	t.Run("test.2.delete", func(t *testing.T) {
		wg.Add(1)
		t1, t2, t3 = nil, nil, nil
		bus.Publish("test.2.delete", "test2")
		wg.Wait()

		assert.Equal(t, testEvent{eventKey: "test.2.delete", data: "test2"}, *t2, "sub 2 should have received wildcard matching event")
		assert.Nil(t, t1, "non-matching sub 1 should be nil")
		assert.Nil(t, t3, "non-matching sub 3 should be nil")
	})

	t.Run("test.3.create", func(t *testing.T) {
		wg.Add(1)
		t1, t2, t3 = nil, nil, nil
		bus.Publish("test.3.create", "test3")
		wg.Wait()

		assert.Equal(t, testEvent{eventKey: "test.3.create", data: "test3"}, *t3, "sub 3 should have received wildcard matching event")
		assert.Nil(t, t1, "non-matching sub 1 should be nil")
		assert.Nil(t, t2, "non-matching sub 2 should be nil")
	})

	t.Run("test.4.delete", func(t *testing.T) {
		t1, t2, t3 = nil, nil, nil
		bus.Publish("test.4.delete", "test4")

		// Sleep to allow processing time
		time.Sleep(2000)

		assert.Nil(t, t1, "non-matching sub 1 should be nil")
		assert.Nil(t, t2, "non-matching sub 2 should be nil")
		assert.Nil(t, t3, "non-matching sub 3 should be nil")
	})
}
