package frontend

import (
	"context"
	"log/slog"
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
	ValidateToken(token string) (err error)

	// GetUser verifies that a token is valid and trusted by us, identifies the user it is tied to, and returns that user.
	//
	// ok tells us if the token is valid. err gives us additional information if the toke is invalid.
	GetUser(ctx context.Context, token string) (user domain.User, err error)
}

type AuthService struct {
	UsernamePasswordAuthenticator UsernamePasswordAuthenticator
	AuthProvider                  AuthProvider
	cookieName                    string
	log                           *slog.Logger
}

func NewAuthService(UsernamePasswordAuthenticator UsernamePasswordAuthenticator, AuthProvider AuthProvider, cookieName *string, logger *slog.Logger) *AuthService {
	s := AuthService{
		UsernamePasswordAuthenticator: UsernamePasswordAuthenticator,
		AuthProvider:                  AuthProvider,
		cookieName:                    "TICKET_SESSION",
		log:                           logger,
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

			u, err := a.AuthProvider.GetUser(r.Context(), cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, u)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (a *AuthService) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		email := r.Form.Get("email")
		if email == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		password := r.Form.Get("password")
		if password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u, err := a.UsernamePasswordAuthenticator.AuthenticateUsernamePassword(r.Context(), email, password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := a.AuthProvider.NewToken(u)
		if err != nil {
			a.log.Error("failed issuing a new token", "user", u, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     a.cookieName,
			Value:    token,
			Path:     "/",
			MaxAge:   604800,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}

		if len(cookie.String()) > 4096 {
			a.log.Error("jwt too long for cookie", "user", u, "error", err, "length", len(cookie.String()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &cookie)
		w.Header().Add("HX-Location", "/")
		w.WriteHeader(http.StatusOK)
	})
}
