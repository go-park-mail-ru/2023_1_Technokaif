package init

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2023_1_Technokaif/init/router"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/jmoiron/sqlx"

	albumRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/repository/postgresql"
	artistRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/repository/postgresql"
	authRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/repository/postgresql"
	trackRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/repository/postgresql"
	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"

	albumUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/usecase"
	artistUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/usecase"
	authUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/usecase"
	trackUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/usecase"
	userUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/usecase"

	albumDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/delivery/http"
	artistDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"
	authDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	trackDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"
	userDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"

	authMiddlware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
)

func Init(db *sqlx.DB, logger logger.Logger) *chi.Mux {
	albumRepo := albumRepository.NewPostgreSQL(db, logger)
	artistRepo := artistRepository.NewPostgreSQL(db, logger)
	authRepo := authRepository.NewPostgreSQL(db, logger)
	trackRepo := trackRepository.NewPostgreSQL(db, logger)
	userRepo := userRepository.NewPostgreSQL(db, logger)

	albumUsecase := albumUsecase.NewUsecase(albumRepo, logger)
	artistUsecase := artistUsecase.NewUsecase(artistRepo, logger)
	authUsecase := authUsecase.NewUsecase(authRepo, userRepo, logger)
	trackUsecase := trackUsecase.NewUsecase(trackRepo, artistRepo, logger)
	userUsecase := userUsecase.NewUsecase(userRepo, logger)

	albumHandler := albumDelivery.NewHandler(albumUsecase, artistUsecase, logger)
	ArtistHandler := artistDelivery.NewHandler(artistUsecase, logger)
	authHandler := authDelivery.NewHandler(authUsecase, logger)
	trackHandler := trackDelivery.NewHandler(trackUsecase, artistUsecase, logger)
	userHandler := userDelivery.NewHandler(userUsecase, logger)

	authMiddlware := authMiddlware.NewMiddleware(authUsecase, logger)

	return router.InitRouter(albumHandler,
		ArtistHandler,
		trackHandler,
		authHandler,
		userHandler,
		authMiddlware,
		logger)
}
