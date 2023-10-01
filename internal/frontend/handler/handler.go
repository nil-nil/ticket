package handler

import (
	"log/slog"
	"net/http"

	"github.com/nil-nil/ticket/internal/frontend/components"
)

func New(log *slog.Logger) DefaultHandler {
	return DefaultHandler{
		log: log,
	}
}

type DefaultHandler struct {
	log *slog.Logger
}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.Post(w, r)
		return
	}
	h.Get(w, r)
}

func (h *DefaultHandler) Get(w http.ResponseWriter, r *http.Request) {
	var props ViewProps
	h.View(w, r, props)
}

func (h *DefaultHandler) Post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Display the view.
	h.View(w, r, ViewProps{
		Count: 1,
	})
}

type ViewProps struct {
	Count int
}

func (h *DefaultHandler) View(w http.ResponseWriter, r *http.Request, props ViewProps) {
	components.Page().Render(r.Context(), w)
}
