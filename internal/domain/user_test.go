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
			{ID: uuid.New(), Tenant: uuid.New(), FirstName: "Bob", LastName: "Test", Email: "bob@test.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	eventDrv := mockEventBusDriver{}

	svc := domain.NewUserService(&repo, &eventDrv)

	t.Run("GetAValidUser", func(t *testing.T) {
		expect := domain.User{ID: uuid.New(), Tenant: uuid.New(), FirstName: "Barry", LastName: "Foo", Email: "barryfoo@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		repo.users = append(repo.users, expect)
		u, err := svc.GetUser(context.Background(), expect.ID)
		assert.NoError(t, err, "getting a valid user should not error")
		assert.Equal(t, expect, u)
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
			{ID: uuid.New(), Tenant: uuid.New(), FirstName: "Bob", LastName: "Test", Email: "bobtest@qux.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	eventDrv := mockEventBusDriver{}
	svc := domain.NewUserService(&repo, &eventDrv)

	t.Run("CreateAValidUser", func(t *testing.T) {
		expect := domain.User{
			FirstName: "Barry",
			LastName:  "Jobson",
			Email:     "barry@bar.com",
			Tenant:    uuid.New(),
		}
		u, err := svc.CreateUser(context.Background(), expect.Tenant, expect.FirstName, expect.LastName, expect.Email)
		assert.NoError(t, err, "creating a valid user should not error")
		assert.Equal(t, expect.FirstName, u.FirstName, "expect correct value")
		assert.Equal(t, expect.LastName, u.LastName, "expect correct value")
		assert.Equal(t, expect.Email, u.Email, "expect correct value")
		assert.Equal(t, expect.Tenant, u.Tenant, "expect correct value")
		assert.NotEqual(t, u.ID, uuid.Nil, "expect non-nil UUID")
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

func (r *mockUserRepository) Create(ctx context.Context, Tenant uuid.UUID, FirstName string, LastName string, Email string) (domain.User, error) {
	u := domain.User{
		ID:        uuid.New(),
		Tenant:    Tenant,
		FirstName: FirstName,
		LastName:  LastName,
		Email:     Email,
	}
	r.users = append(r.users, u)

	return u, nil
}
