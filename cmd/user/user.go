package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/joho/godotenv" // load environment
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/config"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/s3"
	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	userGRPC "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/user/delivery/grpc"
	userProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/user/proto/generated"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"

	userS3 "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/client/s3"
	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/usecase"
)

const (
	maxHeaderBytesHTTP = 1 << 20
	readTimeoutHTTP    = 10 * time.Second
	writeTimeoutHTTP   = 10 * time.Second
)

var (
	reg         = prometheus.NewRegistry()
	grpcMetrics = grpcPrometheus.NewServerMetrics()
)

func main() {
	logger, err := logger.NewLogger(commonHttp.GetReqIDFromContext)
	if err != nil {
		log.Fatalf("logger can not be defined: %v\n", err)
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

	userRepo := userRepository.NewPostgreSQL(db, tables)

	s3Client, err := s3.MakeS3MinioClient(os.Getenv(config.S3HostParam), os.Getenv(config.S3AccessKeyParam), os.Getenv(config.S3SecretKeyParam))
	if err != nil {
		logger.Errorf("Error while connecting to S3: %v", err)
		return
	}
	userS3 := userS3.NewS3AvatarSaver(os.Getenv(config.S3BucketParam), os.Getenv(config.S3AvatarFolderParam), s3Client) 

	userUsecase := userUsecase.NewUsecase(userRepo, userS3)

	listener, err := net.Listen("tcp", os.Getenv(config.UserListenParam))
	defer func() {
		if err := listener.Close(); err != nil {
			logger.Errorf("Error while closing user tcp listener: %v", err)
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

	httpMetricsServer := &http.Server{
		Handler:        promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:           os.Getenv(config.UserExporterListenParam),
		MaxHeaderBytes: maxHeaderBytesHTTP,
		ReadTimeout:    readTimeoutHTTP,
		WriteTimeout:   writeTimeoutHTTP,
	}

	go func() {
		if err := httpMetricsServer.ListenAndServe(); err != nil {
			logger.Errorf("Unable to start a http user metrics server:", err)
		}
	}()
	userProto.RegisterUserServer(server, userGRPC.NewUserGRPC(userUsecase, logger))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-stop
		logger.Info("Server user gracefully shutting down...")

		server.GracefulStop()
	}()

	logger.Info("Starting grpc server user")
	if err := server.Serve(listener); err != nil {
		logger.Errorf("User Server error: %v", err)
		os.Exit(1)
	}
	wg.Wait()
}

func init() {
	_ = godotenv.Load()
}
