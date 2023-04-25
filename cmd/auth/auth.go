package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv" // load environment
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/db/postgresql"
	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	authProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/microservice/grpc/proto"
	authService "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/microservice/grpc/service"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"

	authRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/repository/postgresql"
	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"
)

func main() {
	logger, err := logger.NewLogger(commonHttp.GetReqIDFromRequest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger can not be defined: %v\n", err)
		return
	}

	db, tables, err := postgresql.InitPostgresDB()
	if err != nil {
		logger.Errorf("Error while connecting to database: %v", err)
		return
	}

	userRepo := userRepository.NewPostgreSQL(db, tables, logger)
	authRepo := authRepository.NewPostgreSQL(db, tables, logger)

	listener, err := net.Listen("tcp", ":"+os.Getenv(cmd.AuthPortParam))
	if err != nil {
		logger.Errorf("Cant listen port: %v", err)
		return
	}

	server := grpc.NewServer()
	authProto.RegisterAuthorizationServer(server, authService.NewAuthService(userRepo, authRepo, logger))
	if err := server.Serve(listener); err != nil {
		logger.Errorf("Auth Server error: %v", err)
		return
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error while loading environment: %v", err)
	}
}
