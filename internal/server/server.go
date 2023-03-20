package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
)

// Server is the wrapper for http.Server
type Server struct {
	httpServer *http.Server
}

// RunConfig is server's config for host and port
type RunConfig struct {
	ServerHost string
	ServerPort string
}

func InitServerRunConfig() (RunConfig, error) {
	cfg := RunConfig{
		ServerPort: os.Getenv("SERVER_PORT"),
		ServerHost: os.Getenv("SERVER_HOST"),
	}

	if strings.TrimSpace(cfg.ServerPort) == "" ||
		strings.TrimSpace(cfg.ServerHost) == "" {

		return RunConfig{}, errors.New("invalid server run config")
	}

	return cfg, nil
}

// Run launches http Server on chosen port with given handler
func (s *Server) Run(handler http.Handler, logger logger.Logger) error {
	cfg, err := InitServerRunConfig()
	if err != nil {
		logger.Errorf("error while launching server: %v", err)
		return fmt.Errorf("can't init server config: %w", err)
	}

	s.httpServer = &http.Server{
		Addr:           cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		logger.Errorf("error while launching server: %v", err)
		return err
	}
	return nil
}

// Shutdown gracefully shuts down Server without interrupting active connections
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
