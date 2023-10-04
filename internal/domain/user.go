package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Find(ctx context.Context, ID uuid.UUID) (User, error)
	Create(ctx context.Context, Tenant uuid.UUID, FirstName string, LastName string, Email string) (User, error)
}

type User struct {
	ID        uuid.UUID `eventbus:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Tenant    uuid.UUID

	FirstName string
	LastName  string
	Email     string
}

func NewUserService(repo UserRepository, eventBusDriver EventBusDriver) *UserService {
	eventBus, _ := NewEventBus[User]("users", eventBusDriver)
	return &UserService{
		repo:     repo,
		eventBus: eventBus,
	}
}

type UserService struct {
	repo     UserRepository
	eventBus *EventBus[User]
}

func (s *UserService) GetUser(ctx context.Context, ID uuid.UUID) (User, error) {
	return s.repo.Find(ctx, ID)
}

func (s *UserService) CreateUser(ctx context.Context, Tenant uuid.UUID, FirstName string, LastName string, Email string) (User, error) {
	u, err := s.repo.Create(ctx, Tenant, FirstName, LastName, Email)
	if err != nil {
		return User{}, err
	}

	err = s.eventBus.Publish(fmt.Sprint(u.ID), CreateEvent, u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}
