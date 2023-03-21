package router

import (
	"github.com/go-chi/chi"
	// "github.com/go-chi/chi/v5/middleware"
	swagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2023_1_Technokaif/docs"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http/middleware"
	album "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/delivery/http"
	artist "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"
	auth "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	authM "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
	track "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"
	user "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"
)

// InitRouter describes all app's endpoints and their handlers
func InitRouter(
	album *album.Handler,
	artist *artist.Handler,
	track *track.Handler,
	auth *auth.Handler,
	user *user.Handler,
	authM *authM.Middleware,
	loggger logger.Logger) *chi.Mux {

	r := chi.NewRouter()
	
	// r.Use(middleware.Logger)  // DEBUG
	r.Use(middleware.Logging(loggger))

	r.Get("/swagger/*", swagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.Get("/{userID}", user.Read)
			// TODO add endpoints
		})

		r.Route("/albums", func(r chi.Router) {
			// r.Post("/", album.Create)
			r.Route("/{albumID}", func(r chi.Router) {
				r.Get("/", album.Read)
				// r.Put("/", album.Update)
				// r.Delete("/", album.Delete)

				r.Get("/tracks", album.Tracks)
			})
			r.Get("/feed", album.Feed)
		})

		r.Route("/artists", func(r chi.Router) {
			// r.Post("/", artist.Create)
			r.Route("/{artistID}", func(r chi.Router) {
				r.Get("/", artist.Read)
				// r.Put("/", artist.Update)
				// r.Delete("/", artist.Delete)

				r.Get("/tracks", artist.Tracks)
				r.Get("/albums", artist.Albums)
			})
			r.Get("/feed", artist.Feed)
		})

		r.Route("/tracks", func(r chi.Router) {
			// r.Post("/", track.Create)
			r.Route("/{trackID}", func(r chi.Router) {
				r.Get("/", track.Read)
				// r.Put("/", track.Update)
				// r.Delete("/", track.Delete)
			})
			r.Get("/feed", track.Feed)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", auth.Login)
			r.Post("/signup", auth.SignUp)
			r.With(authM.Authorization).Get("/logout", auth.Logout)
		})

		chi.NewRouteContext()
	})

	return r
}
