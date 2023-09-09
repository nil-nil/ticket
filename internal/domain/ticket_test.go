package domain_test

import (
	"testing"
	"time"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/services/eventbus"
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

func TestGetTicket(t *testing.T) {
	eventDrv := mockEventBusDriver[domain.Ticket]{}
	mockEventBus := eventbus.NewEventBus(&eventDrv)

	svc := domain.NewTicketService(&repo, mockEventBus, mockCache)

	ticket, err := svc.GetTicket(10)
	assert.Equal(t, domain.Ticket{}, ticket, "ticket should be empty")
	assert.EqualError(t, domain.ErrNotFound, err.Error(), "expected not found error")

	ticket, err = svc.GetTicket(3)
	assert.Equal(t, domain.Ticket{ID: 3, Transitions: repo.transitions[3]}, ticket, "ticket should not be empty")
	assert.NoError(t, err, "error should be nil")
}

func TestOpenTicket(t *testing.T) {
	eventDrv := mockEventBusDriver[domain.Ticket]{}
	mockEventBus := eventbus.NewEventBus(&eventDrv)

	svc := domain.NewTicketService(&repo, mockEventBus, mockCache)

	ticket, err := svc.OpenTicket("test")
	assert.Equal(t, uint64(4), ticket.ID, "ticket should have next ID")
	assert.NoError(t, err, "error should be nil")
	assert.Equal(t, *eventDrv.Event, "domain.Ticket:4:create", "expected event matching subject")

	meta := ticket.Meta()
	assert.Equal(t, domain.TicketStatusOpen, meta.Status, "ticket status should be open")
	assert.Equal(t, "test", meta.Description, "ticket should have description provided")
	assert.Nil(t, meta.OwnerID, "ticket owner id should be nil")
}

func TestUpdateTicket(t *testing.T) {
	eventDrv := mockEventBusDriver[domain.Ticket]{}
	mockEventBus := eventbus.NewEventBus(&eventDrv)

	svc := domain.NewTicketService(&repo, mockEventBus, mockCache)

	ticket, err := svc.UpdateTicket(3, domain.TicketUpdateParameters{
		Description: ptr.To("Expected New Description"),
		Status:      domain.TicketStatusBlocked,
		OwnerID:     ptr.To(uint64(99)),
	})
	assert.Equal(t, uint64(3), ticket.ID, "ticket should have same ID")
	assert.NoError(t, err, "error should be nil")
	assert.Equal(t, *eventDrv.Event, "domain.Ticket:3:update", "expected event matching subject")

	meta := ticket.Meta()
	assert.Equal(t, domain.TicketStatusBlocked, meta.Status, "ticket status should have status provided")
	assert.Equal(t, "Expected New Description", meta.Description, "ticket should have description provided")
	assert.NotNil(t, meta.OwnerID, "ticket owner id should no longer be nil")
	assert.Equal(t, uint64(99), *meta.OwnerID, "ticket should have owner id provided")
}

func TestTicketObserver(t *testing.T) {
	eventDrv := mockEventBusDriver[domain.Ticket]{}
	mockEventBus := eventbus.NewEventBus(&eventDrv)

	svc := domain.NewTicketService(&repo, mockEventBus, mockCache)

	t.Run("valid ticket", func(t *testing.T) {
		mockCache.cache["tickets.3"] = "value"
		svc.ObserveTicketEvent(domain.Ticket{ID: 3}, domain.DeleteEvent)
		assert.Equal(t, mockCache.cache["tickets.3"], domain.Ticket{ID: 3, Transitions: repo.transitions[3]}, "cached ticket should be set")
	})
}

type mockTicketRepo struct {
	transitions map[uint64][]domain.TicketTransition
}

func (m *mockTicketRepo) Find(ID uint64) (domain.Ticket, error) {
	transitions, ok := m.transitions[ID]
	if !ok {
		return domain.Ticket{}, domain.ErrNotFound
	}
	ticket := domain.Ticket{
		ID:          ID,
		Transitions: transitions,
	}
	return ticket, nil
}

func (m *mockTicketRepo) Open(Description string) (domain.Ticket, error) {
	ticketId := nextKey(m.transitions)

	m.transitions[ticketId] = []domain.TicketTransition{
		{
			Timestamp:   time.Now(),
			Status:      domain.TicketStatusOpen,
			Description: &Description,
		},
	}

	return domain.Ticket{
		ID:          ticketId,
		Transitions: m.transitions[ticketId],
	}, nil
}

func (m *mockTicketRepo) Update(ID uint64, Params domain.TicketUpdateParameters) (domain.Ticket, error) {
	m.transitions[ID] = append(m.transitions[ID], domain.TicketTransition{
		Timestamp:   time.Now(),
		Status:      Params.Status,
		Description: Params.Description,
		OwnerID:     Params.OwnerID,
	})

	return domain.Ticket{
		ID:          ID,
		Transitions: m.transitions[ID],
	}, nil
}

var repo = mockTicketRepo{
	transitions: map[uint64][]domain.TicketTransition{
		3: {
			{
				Timestamp:   time.Now().Add(-2 * 24 * time.Hour),
				Status:      domain.TicketStatusOpen,
				Description: ptr.To("Opening Description"),
			},
			{
				Timestamp:   time.Now().Add(-1 * 24 * time.Hour),
				Status:      domain.TicketStatusBlocked,
				Description: ptr.To("Expected Description"),
			},
		},
	},
}
