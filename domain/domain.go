package domain

import (
	"errors"
	"time"
)

type UserRepository interface {
	Find(ID uint64) (User, error)
	Create(FirstName string, LastName string) (User, error)
}

type repository struct {
	users UserRepository
}

var (
	repo                       repository
	ErrUninitializedRepository = errors.New("repository has not been initialized")
	ErrUserNotFound            = errors.New("user not found")
)

func InitRepo(userRepository UserRepository) {
	repo = repository{
		users: userRepository,
	}
}

type User struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	FirstName string
	LastName  string
}

func GetUser(ID uint64) (User, error) {
	if repo.users == nil {
		return User{}, ErrUninitializedRepository
	}
	return repo.users.Find(ID)
}

func CreateUser(FirstName string, LastName string) (User, error) {
	if repo.users == nil {
		return User{}, ErrUninitializedRepository
	}
	return repo.users.Create(FirstName, LastName)
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
	ID          uint64
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
