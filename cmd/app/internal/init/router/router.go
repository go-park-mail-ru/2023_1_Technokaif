package router

import (
	"github.com/go-chi/chi/v5"
	swagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2023_1_Technokaif/docs"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http/middleware"
	album "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/delivery/http"
	artist "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"
	auth "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	authM "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
	track "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"
	user "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
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

	r.Use(middleware.Panic(loggger))
	r.Use(middleware.Logging(loggger))

	r.Get("/swagger/*", swagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.With(authM.Authorization).Get("/", user.Read)
				r.With(authM.Authorization).Post("/avatar", user.UploadAvatar)

				r.Route("/favourite", func(r chi.Router) {
					r.Use(authM.Authorization)
					r.Get("/tracks", user.ReadFavouriteTracks)
					r.Get("/albums", user.ReadFavouriteAlbums)
					r.Get("/artists", user.ReadFavouriteArtists)
				})
			})
		})

		r.Route("/albums", func(r chi.Router) {
			r.With(authM.Authorization).Post("/", album.Create)
			r.Route("/{albumID}", func(r chi.Router) {
				r.Get("/", album.Read)
				// r.Put("/", album.Update)
				r.With(authM.Authorization).Delete("/", album.Delete)

				r.With(authM.Authorization).Post("/like", album.Like)
				r.With(authM.Authorization).Post("/unlike", album.UnLike)

				r.Get("/tracks", track.ReadByAlbum)
			})
			r.Get("/feed", album.Feed)
		})

		r.Route("/artists", func(r chi.Router) {
			r.With(authM.Authorization).Post("/", artist.Create)
			r.Route("/{artistID}", func(r chi.Router) {
				r.Get("/", artist.Read)
				// r.Put("/", artist.Update)
				r.With(authM.Authorization).Delete("/", artist.Delete)

				r.With(authM.Authorization).Post("/like", artist.Like)
				r.With(authM.Authorization).Post("/unlike", artist.UnLike)

				r.Get("/tracks", track.ReadByArtist)
				r.Get("/albums", album.ReadByArtist)
			})
			r.Get("/feed", artist.Feed)
		})

		r.Route("/tracks", func(r chi.Router) {
			r.With(authM.Authorization).Post("/", track.Create)
			r.Route("/{trackID}", func(r chi.Router) {
				r.Get("/", track.Read)
				// r.Put("/", track.Update)
				r.With(authM.Authorization).Delete("/", track.Delete)

				r.With(authM.Authorization).Post("/like", track.Like)
				r.With(authM.Authorization).Post("/unlike", track.UnLike)
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
