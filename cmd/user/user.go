package main

import (
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv" // load environment
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/db/postgresql"
	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	userGRPC "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/user/delivery/grpc"
	userProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/user/proto/generated"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"

	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/usecase"
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

	userRepo := userRepository.NewPostgreSQL(db, tables)

	userUsecase := userUsecase.NewUsecase(userRepo)

	listener, err := net.Listen("tcp", ":"+os.Getenv(cmd.UserPortParam))
	if err != nil {
		logger.Errorf("Cant listen port: %v", err)
		return
	}

	server := grpc.NewServer()
	userProto.RegisterUserServer(server, userGRPC.NewUserGRPC(userUsecase, logger))
	if err := server.Serve(listener); err != nil {
		logger.Errorf("Auth Server error: %v", err)
		return
	}
}

func init() {
	_ = godotenv.Load()
}
