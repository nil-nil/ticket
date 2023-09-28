package frontend

import (
	"embed"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nil-nil/ticket/internal/frontend/components"
)

var (
	//go:embed assets/css/* assets/js/*
	assets embed.FS
)

func NewServer() *http.Server {
	const addr = "localhost:8080"
	mux := http.NewServeMux()
	mux.Handle("/", logMiddleware()(templ.Handler(components.Page())))
	mux.Handle("/assets/", logMiddleware()(http.FileServer(http.FS(assets))))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func logMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				log.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
