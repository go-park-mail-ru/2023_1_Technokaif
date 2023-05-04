package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/joho/godotenv" // load environment
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

var (
	reg         = prometheus.NewRegistry()
	grpcMetrics = grpcPrometheus.NewServerMetrics()
)

func main() {
	logger, err := logger.NewLogger(commonHttp.GetReqIDFromContext)
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
	defer func() {
		if err := listener.Close(); err != nil {
			logger.Errorf("Error while closing search tcp listener: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("Cant listen port: %v", err)
		return
	}

	reg.MustRegister(grpcMetrics)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
	)

	grpcMetrics.InitializeMetrics(server)

	httpMetricsServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: os.Getenv(config.SearchExporterListenParam)}
	go func() {
		if err := httpMetricsServer.ListenAndServe(); err != nil {
			logger.Errorf("Unable to start a http search metrics server:", err)
		}
	}()
	searchProto.RegisterSearchServer(server, searchGRPC.NewSearchGRPC(searchUsecase, logger))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-stop
		logger.Info("Server search gracefully shutting down...")

		server.GracefulStop()
	}()

	logger.Info("Starting grpc server search")
	if err := server.Serve(listener); err != nil {
		logger.Errorf("Server search error: %v", err)
		os.Exit(1)
	}
	wg.Wait()
}

func init() {
	_ = godotenv.Load()
}
