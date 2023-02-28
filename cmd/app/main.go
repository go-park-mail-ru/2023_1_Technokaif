package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv" // load environment

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/delivery"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/server"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Logger can not be defined: "+err.Error())
	}

	if err = godotenv.Load(); err != nil {
		logger.Error("Error while loading environment: " + err.Error())
	}

	db, err := repository.NewPostrgresDB(InitDBConfig())
	if err != nil {
		logger.Error("Error while connecting to database: " + err.Error())
		return
	}

	repository := repository.NewRepository(db, logger)
	services := usecase.NewUsecase(repository, logger)
	handler := delivery.NewHandler(services, logger)

	server := new(server.Server)
	host := os.Getenv("SERVERHOST")
	port := os.Getenv("SERVERPORT")

	go func() {
		if err := server.Run(host, port, handler.InitRouter()); err != nil {
			logger.Error("Error while launching server: " + err.Error())
		}
	}()
	logger.Info("Server launched at " + host + ":" + port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Server gracefully shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error("Error while shutting down server: " + err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Error("Error while closing DB-connection: " + err.Error())
	}
}

// InitConfig inits DB configuration from environment variables
func InitDBConfig() repository.Config {
	return repository.Config{
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
		User:     os.Getenv("DBUSER"),
		DBName:   os.Getenv("DBNAME"),
		Password: os.Getenv("DBPASSWORD"),
		SSLMode:  os.Getenv("DBPASSWORD"),
	}
}
