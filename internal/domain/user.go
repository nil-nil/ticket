package domain

import "time"

type UserRepository interface {
	Find(ID uint64) (User, error)
	Create(FirstName string, LastName string) (User, error)
}

type User struct {
	ID        uint64 `eventbus:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	FirstName string
	LastName  string
}

func NewUserService(repo UserRepository, eventBus UserEventBus) *UserService {
	return &UserService{
		repo:     repo,
		eventBus: eventBus,
	}
}

type UserEventBus interface {
	Publish(data User, eventType EventType) error
	Subscribe(subject User, wildcardID bool, eventTypes []EventType, callback func(data User, eventType EventType)) error
}

type UserService struct {
	repo     UserRepository
	eventBus UserEventBus
}

func (s *UserService) GetUser(ID uint64) (User, error) {
	return s.repo.Find(ID)
}

func (s *UserService) CreateUser(FirstName string, LastName string) (User, error) {
	u, err := s.repo.Create(FirstName, LastName)
	if err != nil {
		return User{}, err
	}

	err = s.eventBus.Publish(u, CreateEvent)
	if err != nil {
		return User{}, err
	}

	return u, nil
}
