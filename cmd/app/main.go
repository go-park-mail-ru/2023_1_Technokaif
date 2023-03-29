package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv" // load environment

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/app/internal/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/app/internal/init/app"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/app/internal/server"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// @title		Fluire API
// @version		1.0.1
// @description	Server API for Fluire Streaming Service Application

// @contact.name   Fluire API Support
// @contact.email  yarik1448kuzmin@gmail.com

// @host		localhost:4443
// @BasePath	/api/albums/feed

// @securityDefinitions	AuthKey
// @in					header
// @name				Authorization

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger can not be defined: %v\n", err)
		return
	}

	if err = godotenv.Load(); err != nil {
		logger.Errorf("error while loading environment: %v", err)
		return
	}

	db, tables, err := postgresql.InitPostgresDB()
	if err != nil {
		logger.Errorf("error while connecting to database: %v", err)
		return
	}

	router := app.Init(db, tables, logger)

	var srv server.Server
	go func() {
		if err := srv.Run(router, logger); err != nil {
			logger.Errorf("error while launching server: %v", err)
			os.Exit(1)
		}
	}()
	logger.Info("trying to launch server")
	
	timer := time.AfterFunc(2*time.Second, func() {
		logger.Infof("server launched at %s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
	})
	defer timer.Stop()

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
