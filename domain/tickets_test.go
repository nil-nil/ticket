package domain_test

import (
	"testing"
	"time"

	"github.com/nil-nil/ticket/domain"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestTicketMeta(t *testing.T) {
	transitions := []domain.TicketTransition{
		{
			Timestamp: time.Now().Add(-1 * 12 * time.Hour),
			Status:    domain.TicketStatusUnknown,
			OwnerID:   ptr.To(uint64(99)),
		},
		{
			Timestamp: time.Now().Add(-1 * 24 * time.Hour),
			Status:    domain.TicketStatusOpen,
		},
		{
			Timestamp:   time.Now().Add(-4 * 24 * time.Hour),
			Status:      domain.TicketStatusUnknown,
			Description: ptr.To("Test 1"),
		},
		{
			Timestamp:   time.Now().Add(-2 * 24 * time.Hour),
			Status:      domain.TicketStatusClosed,
			Description: ptr.To("Test 2"),
			OwnerID:     ptr.To(uint64(100)),
		},
	}

	ticket := domain.Ticket{
		ID:          1,
		Transitions: transitions,
	}

	meta := ticket.Meta()

	assert.Equal(t, domain.TicketStatusOpen.String(), meta.Status.String(), "Wrong status")
	assert.NotNil(t, meta.OwnerID, "Missing Owner ID")
	assert.Equal(t, uint64(99), *meta.OwnerID, "Wrong Owner ID")
	assert.Equal(t, "Test 2", meta.Description, "Wrong Description")
}
