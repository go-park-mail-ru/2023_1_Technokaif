package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/jmoiron/sqlx"

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/app/internal/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/app/internal/init/router"

	albumRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/repository/postgresql"
	artistRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/repository/postgresql"
	authRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/repository/postgresql"
	trackRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/repository/postgresql"
	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"

	albumUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/usecase"
	artistUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/usecase"
	authUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/usecase"
	tokenUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token/usecase"
	trackUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/usecase"
	userUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/usecase"

	albumDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/delivery/http"
	artistDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"
	authDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	csrfDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/csrf/delivery/http"
	trackDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"
	userDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"

	authMiddlware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
	csrfMiddlware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/csrf/delivery/http/middleware"
	userMiddlware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http/middleware"
)

func Init(db *sqlx.DB, tables postgresql.PostgreSQLTables, logger logger.Logger) *chi.Mux {
	albumRepo := albumRepository.NewPostgreSQL(db, tables, logger)
	artistRepo := artistRepository.NewPostgreSQL(db, tables, logger)
	authRepo := authRepository.NewPostgreSQL(db, tables, logger)
	trackRepo := trackRepository.NewPostgreSQL(db, tables, logger)
	userRepo := userRepository.NewPostgreSQL(db, tables, logger)

	albumUsecase := albumUsecase.NewUsecase(albumRepo, artistRepo, logger)
	artistUsecase := artistUsecase.NewUsecase(artistRepo, logger)
	authUsecase := authUsecase.NewUsecase(authRepo, userRepo, logger)
	trackUsecase := trackUsecase.NewUsecase(trackRepo, artistRepo, albumRepo, logger)
	userUsecase := userUsecase.NewUsecase(userRepo, logger)
	tokenUsecase := tokenUsecase.NewUsecase(logger)

	albumHandler := albumDelivery.NewHandler(albumUsecase, artistUsecase, logger)
	artistHandler := artistDelivery.NewHandler(artistUsecase, logger)
	authHandler := authDelivery.NewHandler(authUsecase, tokenUsecase, logger)
	trackHandler := trackDelivery.NewHandler(trackUsecase, artistUsecase, logger)
	userHandler := userDelivery.NewHandler(userUsecase, trackUsecase, albumUsecase, artistUsecase, logger)
	csrfHandler := csrfDelivery.NewHandler(tokenUsecase, logger)

	authMiddlware := authMiddlware.NewMiddleware(authUsecase, tokenUsecase, logger)
	userMiddleware := userMiddlware.NewMiddleware(logger)
	csrfMiddlware := csrfMiddlware.NewMiddleware(tokenUsecase, logger)

	return router.InitRouter(
		albumHandler,
		artistHandler,
		trackHandler,
		authHandler,
		userHandler,
		userMiddleware,
		authMiddlware,
		csrfHandler,
		csrfMiddlware,
		logger,
	)
}
