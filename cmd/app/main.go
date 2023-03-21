package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv" // load environment

	initApp "github.com/go-park-mail-ru/2023_1_Technokaif/init/app"
	initDB "github.com/go-park-mail-ru/2023_1_Technokaif/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/server"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
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
	ctx := context.Background()

	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger can not be defined: %v\n", err)
		return
	}

	if err = godotenv.Load(); err != nil {
		logger.Errorf("error while loading environment: %v", err)
		return
	}

	db, tables, err := initDB.InitPostgresDB()
	if err != nil {
		logger.Errorf("error while connecting to database: %v", err)
		return
	}

	router := initApp.Init(db, tables, logger)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(router, logger); err != nil {
			log.Fatalf("can't launch server: %v", err)
		}
	}()
	logger.Infof("server launched at %s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("server gracefully shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("error while shutting down server: %v", err)
	}

	if err := db.Close(); err != nil {
		logger.Errorf("error while closing DB connection: %v", err)
	}
}
