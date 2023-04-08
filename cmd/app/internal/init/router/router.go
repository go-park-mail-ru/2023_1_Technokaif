package router

import (
	"github.com/go-chi/chi/v5"
	swagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2023_1_Technokaif/docs"
	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http/middleware"
	album "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/delivery/http"
	artist "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"
	auth "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	authM "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
	csrf "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/csrf/delivery/http"
	csrfM "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/csrf/delivery/http/middleware"
	track "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"
	user "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const (
	userIdRoute   = "/{" + commonHttp.UserIdUrlParam + "}"
	albumIdRoute  = "/{" + commonHttp.AlbumIdUrlParam + "}"
	artistIdRoute = "/{" + commonHttp.ArtistIdUrlParam + "}"
	trackIdRoute  = "/{" + commonHttp.TrackIdUrlParam + "}"
)

// InitRouter describes all app's endpoints and their handlers
func InitRouter(
	album *album.Handler,
	artist *artist.Handler,
	track *track.Handler,
	auth *auth.Handler,
	user *user.Handler,
	authM *authM.Middleware,
	csrf *csrf.Handler,
	csrfM *csrfM.Middleware,
	loggger logger.Logger) *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.Panic(loggger))
	r.Use(middleware.Logging(loggger))

	r.Get("/swagger/*", swagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.With(authM.Authorization).Route(userIdRoute, func(r chi.Router) {
				r.Get("/", user.Get)

				r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
					r.Post("/update", user.UpdateInfo)
					r.Post("/avatar", user.UploadAvatar)
				})

				r.Route("/favorite", func(r chi.Router) {
					r.Get("/tracks", user.GetFavouriteTracks)
					r.Get("/albums", user.GetFavouriteAlbums)
					r.Get("/artists", user.GetFavouriteArtists)
				})
			})
		})

		r.Route("/albums", func(r chi.Router) {
			r.With(authM.Authorization).Post("/", album.Create)
			r.Route(albumIdRoute, func(r chi.Router) {
				r.Get("/", album.Get)

				r.With(authM.Authorization).Group(func(r chi.Router) {
					r.Get("/tracks", track.GetByAlbum)

					r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
						r.Delete("/", album.Delete)
						r.Post("/like", album.Like)
						r.Post("/unlike", album.UnLike)
					})
				})
			})
			r.Get("/feed", album.Feed)
		})

		r.Route("/artists", func(r chi.Router) {
			r.With(authM.Authorization).Post("/", artist.Create)
			r.Route(artistIdRoute, func(r chi.Router) {
				r.Get("/", artist.Get)

				r.With(authM.Authorization).Group(func(r chi.Router) {
					r.Get("/tracks", track.GetByArtist)

					r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
						r.Delete("/", artist.Delete)
						r.Post("/like", artist.Like)
						r.Post("/unlike", artist.UnLike)
					})
				})
				r.Get("/albums", album.GetByArtist)
			})
			r.Get("/feed", artist.Feed)
		})

		r.With(authM.Authorization).Route("/tracks", func(r chi.Router) {
			r.Post("/", track.Create)
			r.Route(trackIdRoute, func(r chi.Router) {
				r.Get("/", track.Get)
				r.Get("/record", track.GetRecord)

				r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
					r.Delete("/", track.Delete)
					r.Post("/like", track.Like)
					r.Post("/unlike", track.UnLike)
				})
			})
			r.Get("/feed", track.Feed)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", auth.Login)
			r.Post("/signup", auth.SignUp)

			r.With(authM.Authorization).Group(func(r chi.Router) {
				r.Get("/check", auth.IsAuthenticated)
				r.Get("/logout", auth.Logout)
				r.With(csrfM.CheckCSRFToken).Post("/changepass", auth.ChangePassword)
			})
		})

		r.With(authM.Authorization).Get("/csrf", csrf.GetCSRF)
	})

	return r
}
