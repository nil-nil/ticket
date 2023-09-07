package domain

import (
	"fmt"
	"strconv"
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

func NewTicketService(repo TicketRepository, eventBus EventBus, cache CacheDriver) *TicketService {
	svc := &TicketService{
		repo:     repo,
		eventBus: eventBus,
		cache:    cache,
	}

	eventBus.Subscribe(Ticket{}, true, []EventType{CreateEvent, UpdateEvent, DeleteEvent}, svc.observeTicketEvents)

	return svc
}

type TicketService struct {
	repo     TicketRepository
	eventBus EventBus
	cache    CacheDriver
}

func (s *TicketService) GetTicket(ID uint64) (Ticket, error) {
	hit, err := s.cache.Get(fmt.Sprintf("tickets.%d", ID))
	if err == nil {
		ticket, ok := hit.(Ticket)
		if ok {
			return ticket, nil
		}
	}

	return s.repo.Find(ID)
}

func (s *TicketService) OpenTicket(Description string) (Ticket, error) {
	ticket, err := s.repo.Open(Description)
	if err != nil {
		return Ticket{}, err
	}

	err = s.eventBus.Publish(ticket, CreateEvent)
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

	err = s.eventBus.Publish(ticket, UpdateEvent)
	if err != nil {
		return Ticket{}, err
	}
	return ticket, nil
}

func (s *TicketService) observeTicketEvents(subjectType string, subjectId string, eventType EventType) {
	if subjectType != "domain.Ticket" {
		return
	}

	id, err := strconv.ParseUint(subjectId, 10, 64)
	if err != nil {
		return
	}

	ticket, err := s.repo.Find(id)
	if err != nil {
		s.cache.Forget(fmt.Sprintf("tickets.%d", id))
		return
	}

	s.cache.Set(fmt.Sprintf("tickets.%d", id), ticket)
}
