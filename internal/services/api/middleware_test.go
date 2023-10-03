package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/services/api"

	"github.com/stretchr/testify/assert"
)

func mockHandlerFunc(ctx echo.Context, request interface{}) (response interface{}, err error) {
	return ctx.JSON(http.StatusOK, nil), nil
}

type mockAuthProvider struct{}

func (p mockAuthProvider) NewToken(_ domain.User) (token string, err error) {
	return "", nil
}

func (p mockAuthProvider) ValidateToken(_ string) (err error) {
	return nil
}

func (p mockAuthProvider) GetUser(_ string) (user domain.User, err error) {
	return domain.User{}, nil
}

var table = []struct {
	Description  string
	Auth         bool
	Token        string
	ExpectStatus int
}{
	{Description: "Proper Bearer token should succeed", Auth: true, Token: "Bearer 8723082470245709425", ExpectStatus: http.StatusOK},
	{Description: "Token without Bearer should fail", Auth: true, Token: "8723082470245709425", ExpectStatus: http.StatusUnauthorized},
	{Description: "Request without Autorization header should fail", Auth: false, ExpectStatus: http.StatusUnauthorized},
}

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()
	for _, testCase := range table {
		t.Run(testCase.Description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			res := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if testCase.Auth {
				req.Header.Set("Authorization", testCase.Token)
			}

			c := e.NewContext(req, res)

			_, err := api.AuthMiddleware(mockAuthProvider{})(mockHandlerFunc, "TestAuthMiddleware")(c, nil)

			if assert.NoError(t, err) {
				assert.Equal(t, testCase.ExpectStatus, res.Code)
			}
		})
	}
}
