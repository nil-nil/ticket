package api

import (
	"context"
	"net/http"
	"regexp"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

type userMiddlewareValueType struct{}

var userMiddlewareValue = userMiddlewareValueType{}

var tokenRegex = regexp.MustCompile("Bearer (.*)")

func AuthMiddleware(authProvider AuthProvider) runtime.StrictEchoMiddlewareFunc {
	return func(f runtime.StrictEchoHandlerFunc, operationID string) runtime.StrictEchoHandlerFunc {
		return func(echoCtx echo.Context, request interface{}) (response interface{}, err error) {
			authHeader := echoCtx.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echoCtx.NoContent(http.StatusUnauthorized), nil
			}
			submatch := tokenRegex.FindStringSubmatch(authHeader)
			if len(submatch) != 2 {
				return echoCtx.NoContent(http.StatusUnauthorized), nil
			}
			authToken := submatch[1]

			user, err := authProvider.GetUser(authToken)
			if err != nil {
				return echoCtx.NoContent(http.StatusUnauthorized), nil
			}

			ctxWithUser := context.WithValue(echoCtx.Request().Context(), userMiddlewareValue, user)

			requestWithUser := echoCtx.Request().WithContext(ctxWithUser)

			echoCtx.SetRequest(requestWithUser)

			return f(echoCtx, request)
		}
	}
}
