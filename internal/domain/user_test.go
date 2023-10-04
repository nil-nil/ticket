package domain_test

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {

	repo := mockUserRepository{
		users: []domain.User{
			{ID: uuid.New(), FirstName: "Bob", LastName: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	eventDrv := mockEventBusDriver{}

	svc := domain.NewUserService(&repo, &eventDrv)

	t.Run("GetAValidUser", func(t *testing.T) {
		theUUID := uuid.New()
		theUser := domain.User{ID: theUUID, FirstName: "Barry", LastName: "Foo", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		repo.users = append(repo.users, theUser)
		u, err := svc.GetUser(context.Background(), theUUID)
		assert.NoError(t, err, "getting a valid user should not error")
		assert.Equal(t, theUser, u)
	})

	t.Run("GetAnInvalidUser", func(t *testing.T) {
		u, err := svc.GetUser(context.Background(), uuid.New())
		assert.EqualError(t, err, domain.ErrNotFound.Error())
		assert.Equal(t, domain.User{}, u)
	})
}

func TestCreateUser(t *testing.T) {
	repo := mockUserRepository{
		users: []domain.User{
			{ID: uuid.New(), FirstName: "Bob", LastName: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	eventDrv := mockEventBusDriver{}
	svc := domain.NewUserService(&repo, &eventDrv)

	t.Run("CreateAValidUser", func(t *testing.T) {
		first := "Barry"
		last := "Jobson"
		u, err := svc.CreateUser(context.Background(), first, last)
		assert.NoError(t, err, "creating a valid user should not error")
		assert.Equal(t, u.FirstName, first)
		assert.Equal(t, u.LastName, last)
		assert.Equal(t, *eventDrv.EventSubject, fmt.Sprintf("users:%s:create", u.ID), "expected event matching subject")
		assert.Equal(t, eventDrv.EventData, u, "expected matching event data")
	})
}

type mockUserRepository struct {
	users []domain.User
}

func (r *mockUserRepository) Find(ctx context.Context, ID uuid.UUID) (domain.User, error) {
	idx := slices.IndexFunc(r.users, func(u domain.User) bool {
		return u.ID == ID
	})
	if idx == -1 {
		return domain.User{}, domain.ErrNotFound
	}
	return r.users[idx], nil
}

func (r *mockUserRepository) Create(ctx context.Context, FirstName string, LastName string) (domain.User, error) {
	u := domain.User{
		ID:        uuid.New(),
		FirstName: FirstName,
		LastName:  LastName,
	}
	r.users = append(r.users, u)

	return u, nil
}
