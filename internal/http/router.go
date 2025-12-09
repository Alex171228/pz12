package httpx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"example.com/notes-api/internal/http/handlers"
)

// NewRouter создаёт и настраивает HTTP роутер.
func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// базовые middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Swagger UI на /docs (использует сгенерированную документацию из пакета docs)
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"), // URL к swagger.json
	))

	// основное API
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/notes", func(r chi.Router) {
			r.Post("/", h.CreateNote)       // POST /api/v1/notes
			r.Get("/", h.ListNotes)         // GET  /api/v1/notes
			r.Get("/{id}", h.GetNote)       // GET  /api/v1/notes/{id}
			r.Patch("/{id}", h.UpdateNote)  // PATCH /api/v1/notes/{id}
			r.Delete("/{id}", h.DeleteNote) // DELETE /api/v1/notes/{id}
		})
	})

	return r
}
