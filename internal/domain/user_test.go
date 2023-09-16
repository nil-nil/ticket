package domain_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	users map[uint64]domain.User
}

func (r *mockUserRepository) Find(ctx context.Context, ID uint64) (domain.User, error) {
	user, ok := r.users[ID]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (r *mockUserRepository) Create(ctx context.Context, FirstName string, LastName string) (domain.User, error) {
	var lastKey uint64
	for k := range r.users {
		if k > lastKey {
			lastKey = k
		}
	}

	userID := nextKey(r.users)
	r.users[userID] = domain.User{
		ID:        userID,
		FirstName: FirstName,
		LastName:  LastName,
	}

	return r.users[userID], nil
}

func TestGetUser(t *testing.T) {
	repo := mockUserRepository{
		users: map[uint64]domain.User{
			1: {ID: 1, FirstName: "Bob", LastName: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	eventDrv := mockEventBusDriver{}

	svc := domain.NewUserService(&repo, &eventDrv)

	t.Run("get a valid user", func(t *testing.T) {
		u, err := svc.GetUser(context.Background(), 1)
		assert.NoError(t, err, "getting a valid user should not error")
		assert.Equal(t, repo.users[1], u)
	})

	t.Run("get an invalid user", func(t *testing.T) {
		u, err := svc.GetUser(context.Background(), 100)
		assert.EqualError(t, err, domain.ErrNotFound.Error())
		assert.Equal(t, domain.User{}, u)
	})
}

func TestCreateUser(t *testing.T) {
	repo := mockUserRepository{
		users: map[uint64]domain.User{
			1: {ID: 1, FirstName: "Bob", LastName: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	eventDrv := mockEventBusDriver{}
	svc := domain.NewUserService(&repo, &eventDrv)

	t.Run("create a valid user", func(t *testing.T) {
		first := "Barry"
		last := "Jobson"
		u, err := svc.CreateUser(context.Background(), first, last)
		assert.NoError(t, err, "creating a valid user should not error")
		assert.Equal(t, u.FirstName, first)
		assert.Equal(t, u.LastName, last)
		assert.Equal(t, u, repo.users[u.ID])
		assert.Equal(t, *eventDrv.EventSubject, fmt.Sprintf("users:%d:create", u.ID), "expected event matching subject")
		assert.Equal(t, eventDrv.EventData, u, "expected matching event data")
	})
}
