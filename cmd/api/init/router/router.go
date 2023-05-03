package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	playlist "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/delivery/http"
	search "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search/delivery/http"
	track "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"
	user "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"
	userM "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http/middleware"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const (
	userIdRoute     = "/{" + commonHttp.UserIdUrlParam + "}"
	albumIdRoute    = "/{" + commonHttp.AlbumIdUrlParam + "}"
	playlistIdRoute = "/{" + commonHttp.PlaylistIdUrlParam + "}"
	artistIdRoute   = "/{" + commonHttp.ArtistIdUrlParam + "}"
	trackIdRoute    = "/{" + commonHttp.TrackIdUrlParam + "}"
)

// InitRouter describes all app's endpoints and their handlers
func InitRouter(
	albumH *album.Handler,
	playlistH *playlist.Handler,
	artistH *artist.Handler,
	trackH *track.Handler,
	authH *auth.Handler,
	userH *user.Handler,
	userM *userM.Middleware,
	authM *authM.Middleware,
	csrfH *csrf.Handler,
	csrfM *csrfM.Middleware,
	searchH *search.Handler,
	loggger logger.Logger) *chi.Mux {

	r := chi.NewRouter()
	
	r.Use(middleware.Panic(loggger))
	
	r.Use(middleware.SetReqId)
	r.Use(middleware.Logging(loggger))
	r.Use(middleware.Metrics())

	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/swagger/*", swagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.With(authM.Authorization, userM.CheckUserAuthAndResponce).Route(userIdRoute, func(r chi.Router) {
				r.Get("/playlists", playlistH.GetByUser)

				r.Get("/", userH.Get)
				r.Get("/playlists", playlistH.GetByUser)

				r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
					r.Post("/update", userH.UpdateInfo)
					r.With(middleware.RequestBodyMaxSize(user.MaxAvatarMemory)).Post("/avatar", userH.UploadAvatar)
				})

				r.Route("/favorite", func(r chi.Router) {
					r.Get("/tracks", trackH.GetFavorite)
					r.Get("/albums", albumH.GetFavorite)
					r.Get("/playlists", playlistH.GetFavorite)
					r.Get("/artists", artistH.GetFavorite)
				})
			})
		})

		r.With(authM.Authorization).Route("/albums", func(r chi.Router) {
			r.Post("/search", searchH.FindAlbums)

			r.Post("/", albumH.Create)
			r.Route(albumIdRoute, func(r chi.Router) {
				r.Get("/", albumH.Get)

				r.Group(func(r chi.Router) {
					r.Get("/tracks", trackH.GetByAlbum)

					r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
						r.Delete("/", albumH.Delete)
						r.Post("/like", albumH.Like)
						r.Post("/unlike", albumH.UnLike)
					})
				})
			})
			r.Get("/feed", albumH.Feed)
		})

		r.With(authM.Authorization).Route("/playlists", func(r chi.Router) {
			r.Post("/search", searchH.FindPlaylists)

			r.With(csrfM.CheckCSRFToken).Post("/", playlistH.Create)
			r.Route(playlistIdRoute, func(r chi.Router) {
				r.Get("/", playlistH.Get)

				r.Group(func(r chi.Router) {
					r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
						r.Post("/update", playlistH.Update)
						r.Post("/cover", playlistH.UploadCover)
						r.Delete("/", playlistH.Delete)

						r.Post("/like", playlistH.Like)
						r.Post("/unlike", playlistH.UnLike)
					})

					r.Route("/tracks", func(r chi.Router) {
						r.Get("/", trackH.GetByPlaylist)
						r.Route(trackIdRoute, func(r chi.Router) {
							r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
								r.Post("/", playlistH.AddTrack)
								r.Delete("/", playlistH.DeleteTrack)
							})
						})
					})

				})
			})
			r.Get("/feed", playlistH.Feed)
		})

		r.With(authM.Authorization).Route("/artists", func(r chi.Router) {
			r.Post("/search", searchH.FindArtists)

			r.Post("/", artistH.Create)
			r.Route(artistIdRoute, func(r chi.Router) {
				r.Get("/", artistH.Get)

				r.Group(func(r chi.Router) {
					r.Get("/tracks", trackH.GetByArtist)

					r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
						r.Delete("/", artistH.Delete)
						r.Post("/like", artistH.Like)
						r.Post("/unlike", artistH.UnLike)
					})
				})
				r.Get("/albums", albumH.GetByArtist)
			})
			r.Get("/feed", artistH.Feed)
		})

		r.With(authM.Authorization).Route("/tracks", func(r chi.Router) {
			r.Post("/search", searchH.FindTracks)

			r.Post("/", trackH.Create)
			r.Route(trackIdRoute, func(r chi.Router) {
				r.Get("/", trackH.Get)

				r.With(csrfM.CheckCSRFToken).Group(func(r chi.Router) {
					r.Delete("/", trackH.Delete)
					r.Post("/like", trackH.Like)
					r.Post("/unlike", trackH.UnLike)
				})
			})
			r.Get("/feed", trackH.Feed)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authH.Login)
			r.Post("/signup", authH.SignUp)

			r.With(authM.Authorization).Group(func(r chi.Router) {
				r.Get("/", authH.Auth)
				r.Get("/check", authH.IsAuthenticated)
				r.Get("/logout", authH.Logout)
				r.With(csrfM.CheckCSRFToken).Post("/changepass", authH.ChangePassword)
			})
		})

		r.With(authM.Authorization).Get("/csrf", csrfH.GetCSRF)
	})

	return r
}
