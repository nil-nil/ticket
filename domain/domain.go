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
