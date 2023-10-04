package frontend

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/frontend/components"
	"github.com/nil-nil/ticket/internal/infrastructure/ticketjwt"
	"github.com/nil-nil/ticket/internal/services/config"

	"github.com/a-h/templ"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
)

var (
	//go:embed assets/css assets/js
	embedAssets embed.FS
)

func NewServer(config config.Config) *http.Server {
	addr := fmt.Sprintf("%s:%d", config.HTTP.ListenAddress, config.HTTP.Port)

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Extract the assets subfolder form the embedded assets fs
	assets, err := fs.Sub(embedAssets, "assets")
	if err != nil {
		log.Error("error embedding assets", "error", err)
		panic(err)
	}

	// Set up auth
	// TODO: replace placeholder func with real function
	authProvider, err := ticketjwt.NewJwtAuthProvider(
		func(ctx context.Context, userID uuid.UUID) (user domain.User, err error) {
			return domain.User{
				ID:        uuid.New(),
				FirstName: "Tom",
				LastName:  "Salmon",
			}, nil
		},
		[]byte(config.Auth.JWT.PublicKey),
		[]byte(config.Auth.JWT.PrivateKey),
		ticketjwt.GetJWTProtocol(config.Auth.JWT.SigningMethod),
		config.Auth.JWT.TokenLifetime,
	)
	if err != nil {
		log.Error("error creating auth provider", "error", err)
		panic(err)
	}
	// TODO: replace placeholder func with real function
	authSvc := NewAuthService(placeholderAuthenticator, authProvider, nil, log)

	router := httprouter.New()
	logMiddleware := NewLogMiddleware(log, "base")
	router.ServeFiles("/assets/*filepath", http.FS(assets))
	router.Handler(http.MethodGet, "/login", logMiddleware(templ.Handler(components.Login())))
	router.Handler(http.MethodPost, "/login", logMiddleware(authSvc.Login()))

	authRouter := NewHandler(authSvc, log)
	router.HandleMethodNotAllowed = false
	router.NotFound = authRouter

	return &http.Server{
		Addr: addr,
		Handler: handlers.CompressHandler(
			router,
		),
	}
}

// A Log middleware to log http requests
func NewLogMiddleware(log *slog.Logger, router string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			statusWriter := &statusResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(statusWriter, r)
			msg := log
			if user, ok := r.Context().Value(UserContextKey).(domain.User); ok {
				msg = log.With("user", user)
			}
			msg.Info("", "status", statusWriter.statusCode, "handlerLatencyMs", float64(time.Since(start))/float64(time.Millisecond), "method", r.Method, "path", r.URL.Path, "client", r.RemoteAddr, "router", router, "useragent", r.UserAgent())
		})
	}
}

// Embed the http.ResponseWriter and override the WriteHeader to capture the status code
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

type placeholderAuth struct{}

func (p *placeholderAuth) AuthenticateUsernamePassword(_ context.Context, username string, password string) (domain.User, error) {
	return domain.User{
		ID:        uuid.New(),
		FirstName: "Tom",
		LastName:  "Salmon",
	}, nil
}

var placeholderAuthenticator = &placeholderAuth{}
