package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv" // load environment

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/delivery"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/server"
)

func main() {
	logger, err := logger.NewFLogger()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Logger can not be defined: "+err.Error())
	}

	if err = godotenv.Load(); err != nil {
		logger.Error("Error while loading environment: " + err.Error())
	}

	db, err := repository.NewPostrgresDB(InitConfig())
	if err != nil {
		logger.Error("Error while connecting to database: " + err.Error())
	}

	repository := repository.NewRepository(db, logger)
	services := usecase.NewUsecase(repository, logger)
	handler := delivery.NewHandler(services, logger)

	server := new(server.Server)
	if err := server.Run(os.Getenv("SERVERHOST"), os.Getenv("SERVERPORT"), handler.InitRouter()); err != nil {
		logger.Error("Error while launching server: " + err.Error())
	}
}

// InitConfig inits DB configuration from environment variables
func InitConfig() repository.Config {
	return repository.Config{
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
		User:     os.Getenv("DBUSER"),
		DBName:   os.Getenv("DBNAME"),
		Password: os.Getenv("DBPASSWORD"),
		SSLMode:  os.Getenv("DBSSLMODE"),
	}
}
