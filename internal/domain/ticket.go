package domain

import (
	"fmt"
	"time"
)

type TicketRepository interface {
	Find(ID uint64) (Ticket, error)
	Open(Description string) (Ticket, error)
	Update(ID uint64, Params TicketUpdateParameters) (Ticket, error)
}

type TicketUpdateParameters struct {
	Status      TicketStatus
	OwnerID     *uint64
	Description *string
}

type TicketStatus int

const (
	TicketStatusUnknown TicketStatus = iota
	TicketStatusOpen
	TicketStatusInProgress
	TicketStatusBlocked
	TicketStatusClosed
)

func (t TicketStatus) String() string {
	switch t {
	case TicketStatusOpen:
		return "Open"
	case TicketStatusInProgress:
		return "In Progress"
	case TicketStatusBlocked:
		return "Blocked"
	case TicketStatusClosed:
		return "Closed"
	}
	return "Unset"
}

type Ticket struct {
	ID          uint64 `eventbus:"id"`
	Transitions []TicketTransition
}

type TicketTransition struct {
	Timestamp   time.Time
	Status      TicketStatus
	OwnerID     *uint64
	Description *string
}

type TicketMeta struct {
	Description string
	Status      TicketStatus
	OwnerID     *uint64
}

func (t *Ticket) Meta() TicketMeta {
	var (
		meta                 TicketMeta
		descriptionTimestamp time.Time
		statusTimestamp      time.Time
		ownerTimestamp       time.Time
	)
	for _, transition := range t.Transitions {
		if transition.Description != nil && transition.Timestamp.After(descriptionTimestamp) {
			meta.Description = *transition.Description
			descriptionTimestamp = transition.Timestamp
		}
		if transition.Status != TicketStatusUnknown && transition.Timestamp.After(statusTimestamp) {
			meta.Status = transition.Status
			statusTimestamp = transition.Timestamp
		}
		if transition.OwnerID != nil && transition.Timestamp.After(ownerTimestamp) {
			meta.OwnerID = transition.OwnerID
			ownerTimestamp = transition.Timestamp
		}
	}
	return meta
}

func NewTicketService(repo TicketRepository, eventDriver eventBusDriver, cacheDriver cacheDriver) *TicketService {
	cache, _ := NewCache[Ticket]("tickets", cacheDriver)
	eventBus, _ := NewEventBus[Ticket]("tickets", eventDriver)
	svc := &TicketService{
		repo:        repo,
		eventBus:    eventBus,
		ticketCache: cache,
	}

	eventBus.Subscribe(nil, []EventType{CreateEvent, UpdateEvent, DeleteEvent}, svc.ObserveTicketEvent)

	return svc
}

type TicketService struct {
	repo        TicketRepository
	eventBus    *EventBus[Ticket]
	ticketCache *Cache[Ticket]
}

func (s *TicketService) GetTicket(ID uint64) (Ticket, error) {
	hit, err := s.ticketCache.Get(fmt.Sprint(ID))
	if err == nil {
		return hit, nil
	}

	return s.repo.Find(ID)
}

func (s *TicketService) OpenTicket(Description string) (Ticket, error) {
	ticket, err := s.repo.Open(Description)
	if err != nil {
		return Ticket{}, err
	}

	err = s.eventBus.Publish(fmt.Sprint(ticket.ID), CreateEvent, ticket)
	if err != nil {
		return Ticket{}, err
	}
	return ticket, nil
}

func (s *TicketService) UpdateTicket(ID uint64, Params TicketUpdateParameters) (Ticket, error) {
	ticket, err := s.repo.Update(ID, Params)
	if err != nil {
		return Ticket{}, err
	}

	err = s.eventBus.Publish(fmt.Sprint(ticket.ID), UpdateEvent, ticket)
	if err != nil {
		return Ticket{}, err
	}
	return ticket, nil
}

func (s *TicketService) ObserveTicketEvent(eventType EventType, data Ticket) {
	ticket, err := s.repo.Find(data.ID)
	if err != nil {
		s.ticketCache.Forget(fmt.Sprint(data.ID))
		return
	}

	s.ticketCache.Set(fmt.Sprint(ticket.ID), ticket)
}
