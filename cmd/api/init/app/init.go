package app

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/db/postgresql"

	albumRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/repository/postgresql"
	artistRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/repository/postgresql"
	playlistRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/repository/postgresql"
	searchRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search/repository/postgresql"
	trackRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/repository/postgresql"
	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"

	authAgent "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/client/grpc"
	authProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/proto/generated"

	albumUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/usecase"
	artistUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/usecase"
	playlistUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/usecase"
	searchUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search/usecase"
	tokenUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token/usecase"
	trackUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/usecase"
	userUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/usecase"

	albumDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/delivery/http"
	artistDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"
	authDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http"
	csrfDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/csrf/delivery/http"
	playlistDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/delivery/http"
	searchDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search/delivery/http"
	trackDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/delivery/http"
	userDelivery "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"

	authMiddlware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
	csrfMiddlware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/csrf/delivery/http/middleware"
	userMiddlware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http/middleware"
)

type Agents struct {
	*authAgent.AuthAgent
}

func Init(db *sqlx.DB, tables postgresql.PostgreSQLTables, logger logger.Logger) (*chi.Mux, error) {
	albumRepo := albumRepository.NewPostgreSQL(db, tables)
	playlistRepo := playlistRepository.NewPostgreSQL(db, tables)
	artistRepo := artistRepository.NewPostgreSQL(db, tables)
	trackRepo := trackRepository.NewPostgreSQL(db, tables)
	userRepo := userRepository.NewPostgreSQL(db, tables)
	searchRepo := searchRepository.NewPostgreSQL(db, tables)

	agents, err := makeAgents()
	if err != nil {
		return nil, err
	}

	albumUsecase := albumUsecase.NewUsecase(albumRepo, artistRepo)
	playlistUsecase := playlistUsecase.NewUsecase(playlistRepo, trackRepo, userRepo)
	artistUsecase := artistUsecase.NewUsecase(artistRepo)
	trackUsecase := trackUsecase.NewUsecase(trackRepo, artistRepo, albumRepo, playlistRepo)
	userUsecase := userUsecase.NewUsecase(userRepo)
	searchUsecase := searchUsecase.NewUsecase(searchRepo)
	tokenUsecase := tokenUsecase.NewUsecase()

	albumHandler := albumDelivery.NewHandler(albumUsecase, artistUsecase, logger)
	playlistHandler := playlistDelivery.NewHandler(playlistUsecase, trackUsecase, userUsecase, logger)
	artistHandler := artistDelivery.NewHandler(artistUsecase, logger)
	authHandler := authDelivery.NewHandler(agents.AuthAgent, tokenUsecase, logger)
	trackHandler := trackDelivery.NewHandler(trackUsecase, artistUsecase, logger)
	userHandler := userDelivery.NewHandler(userUsecase, logger)
	searchHandler := searchDelivery.NewHandler(searchUsecase,
		albumUsecase, artistUsecase, trackUsecase, playlistUsecase, userUsecase, logger)
	csrfHandler := csrfDelivery.NewHandler(tokenUsecase, logger)

	authMiddlware := authMiddlware.NewMiddleware(agents.AuthAgent, tokenUsecase, logger)
	userMiddleware := userMiddlware.NewMiddleware(logger)
	csrfMiddlware := csrfMiddlware.NewMiddleware(tokenUsecase, logger)

	return router.InitRouter(
		albumHandler,
		playlistHandler,
		artistHandler,
		trackHandler,
		authHandler,
		userHandler,
		userMiddleware,
		authMiddlware,
		csrfHandler,
		csrfMiddlware,
		searchHandler,
		logger,
	), nil
}

func makeAgents() (*Agents, error) {
	grpcAuthConn, err := grpc.Dial(os.Getenv(cmd.AuthHostParam)+":"+os.Getenv(cmd.AuthPortParam),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("can't connect to auth service: %v", err)
	}

	return &Agents{
		AuthAgent: authAgent.NewAuthAgent(authProto.NewAuthorizationClient(grpcAuthConn)),
	}, nil
}
