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
	searchGRPC "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/search/delivery/grpc"
	searchProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/search/proto/generated"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"

	searchRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search/repository/postgresql"
	searchUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search/usecase"
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

	searchRepo := searchRepository.NewPostgreSQL(db, tables)

	searchUsecase := searchUsecase.NewUsecase(searchRepo)

	listener, err := net.Listen("tcp", os.Getenv(config.SearchListenParam))
	if err != nil {
		logger.Errorf("Cant listen port: %v", err)
		return
	}

	server := grpc.NewServer()
	searchProto.RegisterSearchServer(server, searchGRPC.NewSearchGRPC(searchUsecase, logger))
	if err := server.Serve(listener); err != nil {
		logger.Errorf("Auth Server error: %v", err)
		return
	}
}

func init() {
	_ = godotenv.Load()
}
