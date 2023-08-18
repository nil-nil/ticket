package api

import (
	"context"
	"fmt"
)

type Api struct {
}

type UserRespository interface {
	GetUser(userID uint64) (User, error)
}

type PasswordProvider interface {
	GenerateHash(password string) (hash string, err error)
	ComparePasswordAndHash(password, encodedHash string) (match bool, err error)
}

type AuthProvider interface {
	NewToken(user User) (token string, err error)
	ValidateToken(token string) (ok bool, err error)
	GetUser(token string) (ok bool, user User, err error)
}

// Make sure we conform to StrictServerInterface
var _ StrictServerInterface = (*Api)(nil)

func NewApi() *Api {
	api := Api{}
	return &api
}

func (*Api) GetUser(ctx context.Context, req GetUserRequestObject) (GetUserResponseObject, error) {
	authenticatedUser, ok := ctx.Value(userMiddlewareValue).(User)
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	return GetUser200JSONResponse{authenticatedUser}, nil
}
