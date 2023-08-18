package domain

import (
	"errors"
	"time"
)

type UserRepository interface {
	Find(ID uint64) (User, error)
}

type repository struct {
	users UserRepository
}

var (
	repo            repository
	ErrUserNotFound = errors.New("user not found")
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
	return repo.users.Find(ID)
}
