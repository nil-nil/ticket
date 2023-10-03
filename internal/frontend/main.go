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
	logMiddleware := NewLogMiddleware(log)

	// Extract the assets subfolder form the embedded assets fs
	assets, err := fs.Sub(embedAssets, "assets")
	if err != nil {
		log.Error("error embedding assets", "error", err)
		panic(err)
	}

	// Set up auth
	authProvider, err := ticketjwt.NewJwtAuthProvider(
		func(ctx context.Context, userID uint64) (user domain.User, err error) {
			return domain.User{ID: 999}, nil
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
	authSvc := NewAuthService(nil, authProvider, nil)
	auth := authSvc.AuthMiddleware()

	router := httprouter.New()
	router.ServeFiles("/assets/*filepath", http.FS(assets))
	router.Handler(http.MethodGet, "/login", templ.Handler(components.Login()))
	router.Handler(http.MethodGet, "/secure", auth(templ.Handler(components.Login())))

	return &http.Server{
		Addr: addr,
		Handler: handlers.CompressHandler(
			logMiddleware(
				router,
			),
		),
	}
}

// A Log middleware to log http requests
func NewLogMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			statusWriter := &statusResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(statusWriter, r)
			log.Info("", "status", statusWriter.statusCode, "handlerLatencyMs", float64(time.Since(start))/float64(time.Millisecond), "method", r.Method, "path", r.URL.Path, "client", r.RemoteAddr, "useragent", r.UserAgent())
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
