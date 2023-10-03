package frontend

import (
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/frontend/components"
)

type handler struct {
	router         *httprouter.Router
	authSvc        *AuthService
	log            *slog.Logger
	authMiddleware func(http.Handler) http.Handler
	logMiddleware  func(http.Handler) http.Handler
}

func NewHandler(authSvc *AuthService, log *slog.Logger) *handler {
	h := handler{
		router:        httprouter.New(),
		authSvc:       authSvc,
		log:           log,
		logMiddleware: NewLogMiddleware(log, "auth"),
	}

	// Register routes
	h.router.GET("/", h.secure)

	// Set the auth middleware
	h.authMiddleware = h.authSvc.AuthMiddleware()

	return &h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.authMiddleware(h.logMiddleware(h.router)).ServeHTTP(w, r)
}

func (h *handler) secure(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, ok := r.Context().Value(UserContextKey).(domain.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	components.Hello(u.FirstName).Render(r.Context(), w)
}
