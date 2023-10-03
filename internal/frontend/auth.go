package frontend

import (
	"context"
	"net/http"

	"github.com/nil-nil/ticket/internal/domain"
)

type UsernamePasswordAuthenticator interface {
	AuthenticateUsernamePassword(ctx context.Context, username string, password string) (domain.User, error)
}

type AuthProvider interface {
	// NewToken creates a new token tied to the specfied user
	NewToken(user domain.User) (token string, err error)

	// ValidateToken verifies that a token is valid and trusted by us
	ValidateToken(token string) (ok bool, err error)

	// GetUser verifies that a token is valid and trusted by us, identifies the user it is tied to, and returns that user.
	//
	// ok tells us if the token is valid. err gives us additional information if the toke is invalid.
	GetUser(token string) (user domain.User, err error)
}

type AuthService struct {
	UsernamePasswordAuthenticator UsernamePasswordAuthenticator
	AuthProvider                  AuthProvider
	cookieName                    string
}

func NewAuthService(UsernamePasswordAuthenticator UsernamePasswordAuthenticator, AuthProvider AuthProvider, cookieName *string) *AuthService {
	s := AuthService{
		UsernamePasswordAuthenticator: UsernamePasswordAuthenticator,
		AuthProvider:                  AuthProvider,
		cookieName:                    "TICKET_SESSION",
	}
	if cookieName != nil {
		s.cookieName = *cookieName
	}
	return &s
}

type userContextKeyType struct{}

var UserContextKey = userContextKeyType{}

// The middleware handler extracts the user's token from the cookie name set on the AuthService, parses it, gets the associated user and sets that use on the context for the request.
//
// The context key for the user is UserContextKey.
func (a *AuthService) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(a.cookieName)
			if err != nil || cookie == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			u, err := a.AuthProvider.GetUser(cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, u)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
