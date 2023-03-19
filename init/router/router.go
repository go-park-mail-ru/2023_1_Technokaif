package router

import (
	"github.com/go-chi/chi"
	swagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2023_1_Technokaif/docs"
	auth "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	authM "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
)

// InitRouter describes all app's endpoints and their handlers
func InitRouter(auth *auth.AuthHandler, authM *authM.AuthMiddleware) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", swagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		// r.With(h.authorization).Get("/feed", h.feed) // Auth middleware

		r.Route("/auth", func(r chi.Router) {
			r.Route("/login", func(r chi.Router) {
				r.Post("/", auth.Login)
			})

			r.Route("/signup", func(r chi.Router) {
				r.Post("/", auth.SignUp)
			})

			r.With(authM.Authorization).Get("/logout", auth.Logout) // Auth middleware
		})
	})

	return r
}