package domain_test

import (
	"testing"
	"time"

	"github.com/nil-nil/grow/domain"
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

func TestFindUser(t *testing.T) {
	repo := mockUserRepository{
		users: map[uint64]domain.User{
			1: {ID: 1, FirstName: "Bob", LastName: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

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
