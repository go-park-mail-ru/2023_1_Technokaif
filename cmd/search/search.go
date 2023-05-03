package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
		log.Fatalf("Logger can not be defined: %v\n", err)
	}

	db, tables, err := postgresql.InitPostgresDB()
	if err != nil {
		logger.Errorf("Error while connecting to database: %v", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("Error while closing DB connection: %v", err)
		}
	}()

	searchRepo := searchRepository.NewPostgreSQL(db, tables)

	searchUsecase := searchUsecase.NewUsecase(searchRepo)

	listener, err := net.Listen("tcp", os.Getenv(config.SearchListenParam))
	defer listener.Close()

	if err != nil {
		logger.Errorf("Cant listen port: %v", err)
		return
	}

	server := grpc.NewServer()
	searchProto.RegisterSearchServer(server, searchGRPC.NewSearchGRPC(searchUsecase, logger))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		<-stop
		logger.Info("Server search gracefully shutting down...")

		server.GracefulStop()
		wg.Done()
	}()

	logger.Info("Starting grpc server server")
	if err := server.Serve(listener); err != nil {
		logger.Errorf("Server search error: %v", err)
		os.Exit(1)
	}
	wg.Wait()
}

func init() {
	_ = godotenv.Load()
}
