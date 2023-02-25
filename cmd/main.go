package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv" // load environment
	"go.uber.org/zap"          // logger

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/delivery"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/repository"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/usecase"
	"github.com/go-park-mail-ru/2023_1_Technokaif/server"
)

func main() {
	logger, err := zap.NewProduction()
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

	repository := repository.NewRepository(db)
	services := usecase.NewUsecase(repository)
	handler := delivery.NewHandler(services)

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
