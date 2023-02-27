package delivery

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase"
)

type Handler struct {
	services *usecase.Usecase
	log      *logger.Logger
}

func NewHandler(u *usecase.Usecase, l *logger.Logger) *Handler {
	return &Handler{
		services: u,
		log:      l,
	}
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I didn't hit her! I did not... Oh, hi, Mark!"))
}

// InitRouter describes all app's endpoints and their handlers
func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", h.index)

	r.With(h.Authorization).Get("/feed", h.feed)

	r.Route("/auth", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			r.Post("/", h.login)
		})

		r.Route("/signup", func(r chi.Router) {
			r.Post("/", h.signUp)
		})

		r.Get("/logout", h.logout)
	})

	return r
}
