package app

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/config"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/s3"

	albumRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/repository/postgresql"
	artistRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/repository/postgresql"
	playlistRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/repository/postgresql"
	trackRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/repository/postgresql"
	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"

	playlistS3 "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/client/s3"

	authAgent "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/client/grpc"
	authProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/proto/generated"

	searchProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/search/proto/generated"
	searchAgent "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search/client/grpc"

	userProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/user/proto/generated"
	userAgent "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/client/grpc"

	albumUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/usecase"
	artistUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/usecase"
	playlistUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/usecase"
	tokenUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token/usecase"
	trackUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/usecase"

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
	*searchAgent.SearchAgent
	*userAgent.UserAgent
}

func Init(db *sqlx.DB, tables postgresql.PostgreSQLTables, c *cron.Cron, logger logger.Logger) (*chi.Mux, error) {
	albumRepo := albumRepository.NewPostgreSQL(db, tables)
	playlistRepo := playlistRepository.NewPostgreSQL(db, tables)
	artistRepo := artistRepository.NewPostgreSQL(db, tables)
	trackRepo := trackRepository.NewPostgreSQL(db, tables)
	userRepo := userRepository.NewPostgreSQL(db, tables)

	agents, err := makeAgents()
	if err != nil {
		return nil, err
	}

	s3Client, err := s3.MakeS3MinioClient(os.Getenv(config.S3HostParam), os.Getenv(config.S3AccessKeyParam), os.Getenv(config.S3SecretKeyParam))
	if err != nil {
		return nil, fmt.Errorf("error while connecting to S3: %v", err)
	}
	playlistS3 := playlistS3.NewS3PlaylistCoverSaver(os.Getenv(config.S3BucketParam), os.Getenv(config.S3PlaylistCoversFolderParam), s3Client)

	albumUsecase := albumUsecase.NewUsecase(albumRepo, artistRepo)
	playlistUsecase := playlistUsecase.NewUsecase(playlistRepo, trackRepo, userRepo, playlistS3)
	artistUsecase := artistUsecase.NewUsecase(artistRepo)
	trackUsecase := trackUsecase.NewUsecase(trackRepo, artistRepo, albumRepo, playlistRepo)
	tokenUsecase := tokenUsecase.NewUsecase()

	c.AddFunc("@every 10m", func() {
		logger.Info("Count all listens started..")
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("PANIC (recovered): %s\n stacktrace:\n%s", err, string(debug.Stack()))
			}
		}()

		if err := trackUsecase.UpdateAllListens(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("Count all listens fail: %v", err))
			return
		}
		logger.Info("Count all listens succeeded")
	})

	albumHandler := albumDelivery.NewHandler(albumUsecase, artistUsecase, logger)
	playlistHandler := playlistDelivery.NewHandler(playlistUsecase, trackUsecase, agents.UserAgent, logger)
	artistHandler := artistDelivery.NewHandler(artistUsecase, logger)
	authHandler := authDelivery.NewHandler(agents.AuthAgent, tokenUsecase, logger)
	trackHandler := trackDelivery.NewHandler(trackUsecase, artistUsecase, logger)
	userHandler := userDelivery.NewHandler(agents.UserAgent, logger)
	searchHandler := searchDelivery.NewHandler(agents.SearchAgent,
		albumUsecase, artistUsecase, trackUsecase, playlistUsecase, agents.UserAgent, logger)
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
	grpcAuthConn, err := grpc.Dial(os.Getenv(config.AuthConnectParam),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("can't connect to auth service: %v", err)
	}

	grpcSearchConn, err := grpc.Dial(os.Getenv(config.SearchConnectParam),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("can't connect to search service: %v", err)
	}

	grpcUserConn, err := grpc.Dial(os.Getenv(config.UserConnectParam),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("can't connect to user service: %v", err)
	}

	return &Agents{
		AuthAgent:   authAgent.NewAuthAgent(authProto.NewAuthorizationClient(grpcAuthConn)),
		SearchAgent: searchAgent.NewAuthAgent(searchProto.NewSearchClient(grpcSearchConn)),
		UserAgent:   userAgent.NewUserAgent(userProto.NewUserClient(grpcUserConn)),
	}, nil
}
