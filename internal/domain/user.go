package domain

import "time"

type UserRepository interface {
	Find(ID uint64) (User, error)
	Create(FirstName string, LastName string) (User, error)
}

type User struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	FirstName string
	LastName  string
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserService struct {
	repo UserRepository
}

func (s *UserService) GetUser(ID uint64) (User, error) {
	return s.repo.Find(ID)
}

func (s *UserService) CreateUser(FirstName string, LastName string) (User, error) {
	return s.repo.Create(FirstName, LastName)
}
