package router

import (
	"github.com/go-chi/chi"
	swagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2023_1_Technokaif/docs"
	album "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/delivery/http"
	artist "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"
	auth "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	authM "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
	track "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"
	user "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"
)

// InitRouter describes all app's endpoints and their handlers
func InitRouter(album *album.AlbumHandler,
	artist *artist.ArtistHandler,
	track *track.TrackHandler,
	auth *auth.AuthHandler,
	user *user.UserHandler,
	authM *authM.AuthMiddleware) *chi.Mux {

	r := chi.NewRouter()

	r.Get("/swagger/*", swagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Get("/{userID}", user.GetByID)
			r.Get("/{userID}/brief", user.GetBriefByID)
		})

		r.Route("/album", func(r chi.Router) {
			r.Get("/{albumID}", album.GetByID)
			r.Get("/feed", album.Feed)
		})

		r.Route("/artist", func(r chi.Router) {
			r.Get("/{artistID}", artist.GetByID)
			r.Get("/feed", artist.Feed)
		})

		r.Route("/track", func(r chi.Router) {
			r.Get("/{trackID}", track.GetByID)
			r.Get("/feed", track.Feed)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Route("/login", func(r chi.Router) {
				r.Post("/", auth.Login)
			})

			r.Route("/signup", func(r chi.Router) {
				r.Post("/", auth.SignUp)
			})

			r.With(authM.Authorization).Get("/logout", auth.Logout)
		})
	})

	return r
}
