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
