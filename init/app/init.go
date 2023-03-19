package app

import (
	"github.com/go-chi/chi"
	"github.com/go-park-mail-ru/2023_1_Technokaif/init/router"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/jmoiron/sqlx"

	albumRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/repository/postgresql"
	artistRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/repository/postgresql"
	authRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/repository/postgresql"
	trackRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/repository/postgresql"

	// userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"

	albumUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/usecase"
	artistUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/usecase"
	authUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/usecase"
	trackUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/usecase"

	// userUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/usecase"

	albumDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/delivery/http"
	artistDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"
	authDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	trackDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"

	// userDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"

	authMiddlware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
)

func Init(db *sqlx.DB, logger logger.Logger) *chi.Mux {
	albumRepo := albumRepository.NewAlbumPostgres(db, logger)
	artistRepo := artistRepository.NewArtistPostgres(db, logger)
	authRepo := authRepository.NewAuthPostgres(db, logger)
	trackRepo := trackRepository.NewTrackPostgres(db, logger)
	// userRepo := userRepository.NewUserPostgres(db, logger)

	albumUsecase := albumUsecase.NewAlbumUsecase(albumRepo, logger)
	artistUsecase := artistUsecase.NewArtistUsecase(artistRepo, logger)
	authUsecase := authUsecase.NewAuthUsecase(authRepo, logger)
	trackUsecase := trackUsecase.NewTrackUsecase(trackRepo, logger)
	// userUsecase := userUsecase.NewUserUsecase(userRepo, logger)

	albumHandler := albumDelivery.NewAlbumHandler(albumUsecase, artistUsecase, logger)
	ArtistHandler := artistDelivery.NewArtistHandler(artistUsecase, logger)
	authHandler := authDelivery.NewAuthHandler(authUsecase, logger)
	trackHandler := trackDelivery.NewTrackHandler(trackUsecase, artistUsecase, logger)
	// userHandler := userDelivery.NewUserHandler(userUsecase, logger)

	authMiddlware := authMiddlware.NewAuthMiddleware(authUsecase, logger)

	return router.InitRouter(albumHandler,
		ArtistHandler,
		trackHandler,
		authHandler,
		// userHandler,
		authMiddlware)
}
