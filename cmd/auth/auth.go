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
	authGRPC "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/delivery/grpc"
	authProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/proto/generated"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"

	authRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/repository/postgresql"
	authUsecase "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/usecase"
	userRepository "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/repository/postgresql"
)

var (
	reg = prometheus.NewRegistry()

	grpcMetrics = grpcPrometheus.NewServerMetrics()

	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "fluire",
		Subsystem: "auth",
		Name:      "req_counter",
	}, []string{"name"})
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
	authRepo := authRepository.NewPostgreSQL(db, tables)

	authUsecase := authUsecase.NewUsecase(authRepo, userRepo)

	listener, err := net.Listen("tcp", os.Getenv(config.AuthListenParam))
	defer listener.Close()

	if err != nil {
		logger.Errorf("Cant listen port: %v", err)
		return
	}
	reg.MustRegister(grpcMetrics, customizedCounterMetric.WithLabelValues("auth_test"))

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
	)

	grpcMetrics.InitializeMetrics(server)

	httpMetricsServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: os.Getenv("AUTH_EXPORTER_ENDPOINT")}
	go func() {
		if err := httpMetricsServer.ListenAndServe(); err != nil {
			logger.Errorf("Unable to start a http album metrics server:", err)
		}
	}()

	authProto.RegisterAuthorizationServer(server, authGRPC.NewAuthGRPC(authUsecase, logger))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-stop
		logger.Info("Server auth gracefully shutting down...")

		server.GracefulStop()
	}()

	logger.Info("Starting grpc server auth")
	if err := server.Serve(listener); err != nil {
		logger.Errorf("Auth Server error: %v", err)
		os.Exit(1)
	}
	wg.Wait()
}

func init() {
	_ = godotenv.Load()
}
