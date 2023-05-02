package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv" // load environment
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/config"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/db/postgresql"
	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	authGRPC "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/delivery/grpc"
	authProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/proto/generated"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"

	authRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/repository/postgresql"
	authUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/usecase"
	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"
)

func main() {
	logger, err := logger.NewLogger(commonHttp.GetReqIDFromRequest)
	if err != nil {
		log.Fatalf("logger can not be defined: %v\n", err)
	}

	db, tables, err := postgresql.InitPostgresDB()
	if err != nil {
		logger.Errorf("Error while connecting to database: %v", err)
		return
	}

	userRepo := userRepository.NewPostgreSQL(db, tables)
	authRepo := authRepository.NewPostgreSQL(db, tables)

	authUsecase := authUsecase.NewUsecase(authRepo, userRepo)

	listener, err := net.Listen("tcp", os.Getenv(config.AuthListenParam))
	if err != nil {
		logger.Errorf("Cant listen port: %v", err)
		return
	}

	server := grpc.NewServer()
	authProto.RegisterAuthorizationServer(server, authGRPC.NewAuthGRPC(authUsecase, logger))
	if err := server.Serve(listener); err != nil {
		logger.Errorf("Auth Server error: %v", err)
		return
	}
}

func init() {
	_ = godotenv.Load()
}
