package domain

import (
	"context"
	"fmt"
	"time"
)

type UserRepository interface {
	Find(ctx context.Context, ID uint64) (User, error)
	Create(ctx context.Context, FirstName string, LastName string) (User, error)
}

type User struct {
	ID        uint64 `eventbus:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	FirstName string
	LastName  string
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

func (s *UserService) GetUser(ctx context.Context, ID uint64) (User, error) {
	return s.repo.Find(ctx, ID)
}

func (s *UserService) CreateUser(ctx context.Context, FirstName string, LastName string) (User, error) {
	u, err := s.repo.Create(ctx, FirstName, LastName)
	if err != nil {
		return User{}, err
	}

	err = s.eventBus.Publish(fmt.Sprint(u.ID), CreateEvent, u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}
