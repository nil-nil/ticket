package domain_test

import (
	"testing"
	"time"

	"github.com/nil-nil/ticket/domain"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	users map[uint64]domain.User
}

func (r *mockUserRepository) Find(ID uint64) (domain.User, error) {
	user, ok := r.users[ID]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *mockUserRepository) Create(FirstName string, LastName string) (domain.User, error) {
	var lastKey uint64
	for k := range r.users {
		if k > lastKey {
			lastKey = k
		}
	}

	userID := lastKey + 1
	r.users[userID] = domain.User{
		ID:        userID,
		FirstName: FirstName,
		LastName:  LastName,
	}

	return r.users[userID], nil
}

func TestFindUser(t *testing.T) {
	domain.InitRepo(nil)
	repo := mockUserRepository{
		users: map[uint64]domain.User{
			1: {ID: 1, FirstName: "Bob", LastName: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	t.Run("uninitialized user repo", func(t *testing.T) {
		u, err := domain.GetUser(100)
		assert.EqualError(t, err, domain.ErrUninitializedRepository.Error())
		assert.Equal(t, domain.User{}, u)
	})

	domain.InitRepo(&repo)

	t.Run("get a valid user", func(t *testing.T) {
		u, err := domain.GetUser(1)
		assert.NoError(t, err, "getting a valid user should not error")
		assert.Equal(t, repo.users[1], u)
	})

	t.Run("get an invalid user", func(t *testing.T) {
		u, err := domain.GetUser(100)
		assert.EqualError(t, err, domain.ErrUserNotFound.Error())
		assert.Equal(t, domain.User{}, u)
	})
}

func TestGetUser(t *testing.T) {
	domain.InitRepo(nil)
	repo := mockUserRepository{
		users: map[uint64]domain.User{
			1: {ID: 1, FirstName: "Bob", LastName: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	t.Run("uninitialized user repo", func(t *testing.T) {
		u, err := domain.CreateUser("Barry", "Jobson")
		assert.EqualError(t, err, domain.ErrUninitializedRepository.Error())
		assert.Equal(t, domain.User{}, u)
	})

	domain.InitRepo(&repo)

	t.Run("create a valid user", func(t *testing.T) {
		first := "Barry"
		last := "Jobson"
		u, err := domain.CreateUser(first, last)
		assert.NoError(t, err, "creating a valid user should not error")
		assert.Equal(t, u.FirstName, first)
		assert.Equal(t, u.LastName, last)
		assert.Equal(t, u, repo.users[u.ID])
	})
}
