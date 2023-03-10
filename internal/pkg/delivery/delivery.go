package delivery

import (
	"github.com/go-chi/chi"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase"

	_ "github.com/go-park-mail-ru/2023_1_Technokaif/docs"
	swagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services *usecase.Usecase
	logger   logger.Logger
}

func NewHandler(u *usecase.Usecase, l logger.Logger) *Handler {
	return &Handler{
		services: u,
		logger:   l,
	}
}

// InitRouter describes all app's endpoints and their handlers
func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", swagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		r.With(h.authorization).Get("/feed", h.feed) // Auth middleware

		r.Route("/auth", func(r chi.Router) {
			r.Route("/login", func(r chi.Router) {
				r.Post("/", h.login)
			})

			r.Route("/signup", func(r chi.Router) {
				r.Post("/", h.signUp)
			})

			r.With(h.authorization).Get("/logout", h.logout) // Auth middleware
		})
	})

	return r
}
