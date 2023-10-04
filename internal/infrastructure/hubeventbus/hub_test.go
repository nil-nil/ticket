package hubeventbus_test

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
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
		t1, t2, t3          *testEvent
		uuid1, uuid2, uuid3 = uuid.New(), uuid.New(), uuid.New()
	)
	bus.Subscribe("test."+uuid1.String()+".create", func(eventKey string, data interface{}) {
		defer wg.Done()
		t1 = &testEvent{eventKey: eventKey, data: data}
	})

	bus.Subscribe("test."+uuid2.String()+".*", func(eventKey string, data interface{}) {
		defer wg.Done()
		t2 = &testEvent{eventKey: eventKey, data: data}
	})

	bus.Subscribe("test.*.create", func(eventKey string, data interface{}) {
		defer wg.Done()
		t3 = &testEvent{eventKey: eventKey, data: data}
	})

	// Give the subs time to start
	time.Sleep(500000)

	t.Run("test.uuid1.create", func(t *testing.T) {
		wg.Add(2)
		t1, t2, t3 = nil, nil, nil
		bus.Publish("test."+uuid1.String()+".create", domain.User{ID: uuid1, FirstName: "Test", LastName: "Abc"})
		wg.Wait()

		assert.Nil(t, t2, "non-matching sub 2 should be nil")
		assert.Equal(t, testEvent{eventKey: "test." + uuid1.String() + ".create", data: domain.User{ID: uuid1, FirstName: "Test", LastName: "Abc"}}, *t1, "sub 1 should have received matching event")
		assert.Equal(t, testEvent{eventKey: "test." + uuid1.String() + ".create", data: domain.User{ID: uuid1, FirstName: "Test", LastName: "Abc"}}, *t3, "sub 3 should have received wildcard matching event")
	})

	t.Run("test.uuid2.delete", func(t *testing.T) {
		wg.Add(1)
		t1, t2, t3 = nil, nil, nil
		bus.Publish("test."+uuid2.String()+".delete", "test2")
		wg.Wait()

		assert.Equal(t, testEvent{eventKey: "test." + uuid2.String() + ".delete", data: "test2"}, *t2, "sub 2 should have received wildcard matching event")
		assert.Nil(t, t1, "non-matching sub 1 should be nil")
		assert.Nil(t, t3, "non-matching sub 3 should be nil")
	})

	t.Run("test.uuid3.create", func(t *testing.T) {
		wg.Add(1)
		t1, t2, t3 = nil, nil, nil
		bus.Publish("test."+uuid3.String()+".create", "test3")
		wg.Wait()

		assert.Equal(t, testEvent{eventKey: "test." + uuid3.String() + ".create", data: "test3"}, *t3, "sub 3 should have received wildcard matching event")
		assert.Nil(t, t1, "non-matching sub 1 should be nil")
		assert.Nil(t, t2, "non-matching sub 2 should be nil")
	})

	t.Run("test.uuid4.delete", func(t *testing.T) {
		t1, t2, t3 = nil, nil, nil
		bus.Publish("test."+uuid.New().String()+".delete", "test4")

		// Sleep to allow processing time
		time.Sleep(2000)

		assert.Nil(t, t1, "non-matching sub 1 should be nil")
		assert.Nil(t, t2, "non-matching sub 2 should be nil")
		assert.Nil(t, t3, "non-matching sub 3 should be nil")
	})
}
