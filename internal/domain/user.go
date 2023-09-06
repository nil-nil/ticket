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

func NewUserService(repo UserRepository, eventBus EventBus) *UserService {
	return &UserService{
		repo:     repo,
		eventBus: eventBus,
	}
}

type UserService struct {
	repo     UserRepository
	eventBus EventBus
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
