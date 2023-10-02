package api

import (
	"context"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/nil-nil/ticket/internal/domain"
)

type Api struct {
}

type UserRespository interface {
	GetUser(userID uint64) (domain.User, error)
}

type PasswordProvider interface {
	GenerateHash(password string) (hash string, err error)
	ComparePasswordAndHash(password, encodedHash string) (match bool, err error)
}

type AuthProvider interface {
	NewToken(user domain.User) (token string, err error)
	ValidateToken(token string) (ok bool, err error)
	GetUser(token string) (ok bool, user domain.User, err error)
}

// Make sure we conform to StrictServerInterface
var _ StrictServerInterface = (*Api)(nil)

func NewApi() *Api {
	api := Api{}
	return &api
}

func (*Api) GetUser(ctx context.Context, req GetUserRequestObject) (GetUserResponseObject, error) {
	authenticatedUser, ok := ctx.Value(userMiddlewareValue).(domain.User)
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	u := User{
		Id:        authenticatedUser.ID,
		CreatedAt: types.Date{Time: authenticatedUser.CreatedAt},
		UpdatedAt: types.Date{Time: authenticatedUser.UpdatedAt},
		FirstName: authenticatedUser.FirstName,
		LastName:  authenticatedUser.LastName,
	}
	if authenticatedUser.DeletedAt != nil {
		u.DeletedAt = &types.Date{Time: *authenticatedUser.DeletedAt}
	}

	return GetUser200JSONResponse{u}, nil
}
