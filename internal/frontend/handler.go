package frontend

import (
	"log/slog"
	"net/http"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/frontend/components"
)

type handler struct {
	authSvc *AuthService
	log     *slog.Logger
}

func (h *handler) Secure() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(UserContextKey).(domain.User)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		components.Hello(u.FirstName).Render(r.Context(), w)
	})
}
