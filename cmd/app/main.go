package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"strings"

	"github.com/joho/godotenv" // load environment
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/delivery"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/server"
)

//	@title			Fluire API
//	@version		0.1.0
//	@description	Server API for Fluire Streaming Service Application

//	@host		localhost:4443
//	@BasePath	/feed

//	@securityDefinitions	AuthKey
//	@in						header
//	@name					Authorization

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Logger can not be defined: %v\n", err)
		return
	}

	if err = godotenv.Load(); err != nil {
		logger.Errorf("Error while loading environment: %v", err)
		return
	}

	dbConfig, err := initDBConfig()
	if err != nil {
		logger.Errorf("Error while connecting to database: %v", err)
		return
	}

	db, err := repository.NewPostrgresDB(dbConfig)
	if err != nil {
		logger.Errorf("Error while connecting to database: %v", err)
		return
	}

	repository := repository.NewRepository(db, logger)
	services := usecase.NewUsecase(repository, logger)
	handler := delivery.NewHandler(services, logger)

	server := new(server.Server)
	serverRunConfig, err := initServerRunConfig()
	if err != nil {
		logger.Errorf("Error while launching server: %v", err)
		return
	}

	go func() {
		if err := server.Run(serverRunConfig, handler.InitRouter()); err != nil {
			logger.Errorf("Error while launching server: %v", err)
			log.Fatalf("Can't launch server: %v", err)
		}
	}()
	logger.Infof("Server launched at %s:%s", serverRunConfig.ServerHost, serverRunConfig.ServerPort)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Server gracefully shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Errorf("Error while shutting down server: %v", err)
	}

	if err := db.Close(); err != nil {
		logger.Errorf("Error while closing DB-connection: %v", err)
	}
}

// InitConfig inits DB configuration from environment variables
func initDBConfig() (repository.Config, error) {  // TODO CHECK FIELDS
	cfg := repository.Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if 	strings.TrimSpace(cfg.DBHost) == "" ||
		strings.TrimSpace(cfg.DBPort) == "" ||
		strings.TrimSpace(cfg.DBUser) == "" ||
		strings.TrimSpace(cfg.DBName) == "" ||
		strings.TrimSpace(cfg.DBPassword) == "" ||
		strings.TrimSpace(cfg.DBSSLMode) == "" {
		
		return repository.Config{}, errors.New("invalid db config")
	}


	return cfg, nil
}

func initServerRunConfig() (server.RunConfig, error) {
	cfg := server.RunConfig{
		ServerPort: os.Getenv("SERVER_PORT"),
		ServerHost: os.Getenv("SERVER_HOST"),
	}

	if 	strings.TrimSpace(cfg.ServerPort) == "" ||
		strings.TrimSpace(cfg.ServerHost) == "" {
		
		return server.RunConfig{}, errors.New("invalid server run config")
	}

	return cfg, nil
}