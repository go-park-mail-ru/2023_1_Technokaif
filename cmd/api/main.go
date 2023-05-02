package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv" // load environment

	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/config"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/api/init/app"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/api/init/server"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/internal/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/file"
	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"

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

	logger, err := logger.NewLogger(commonHttp.GetReqIDFromRequest)
	if err != nil {
		log.Fatalf("logger can not be defined: %v\n", err)
	}

	if err := file.InitPaths(); err != nil {
		logger.Errorf("can't init paths: %v", err)
		return
	}

	db, tables, err := postgresql.InitPostgresDB()
	if err != nil {
		logger.Errorf("error while connecting to database: %v", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("error while closing DB connection: %v", err)
		}
	}()

	router, err := app.Init(db, tables, logger)
	if err != nil {
		logger.Errorf("error while initialization api app: %v", err)
		return
	}

	var srv server.Server
	if err := srv.Init(os.Getenv(config.ApiListenParam), router); err != nil {
		logger.Errorf("error while launching server: %v", err)
	}

	go func() {
		if err := srv.Run(); err != nil {
			logger.Errorf("server error: %v", err)
			os.Exit(1)
		}
	}()
	logger.Info("trying to launch server")

	timer := time.AfterFunc(2*time.Second, func() {
		logger.Infof("server launched at :%s", os.Getenv(config.ApiListenParam))
	})
	defer timer.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("server gracefully shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("error while shutting down server: %v", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error while loading environment: %v", err)
	}
}
