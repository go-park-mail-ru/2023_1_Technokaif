package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
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

// runConfig is server's config for host and port
type runConfig struct {
	ServerHost string
	ServerPort string
}

func initServerRunConfig() (runConfig, error) {
	cfg := runConfig{
		ServerPort: os.Getenv("SERVER_PORT"),
		ServerHost: os.Getenv("SERVER_HOST"),
	}

	if strings.TrimSpace(cfg.ServerPort) == "" ||
		strings.TrimSpace(cfg.ServerHost) == "" {

		return runConfig{}, errors.New("invalid server run config")
	}

	return cfg, nil
}

func (s *Server) Init(handler http.Handler) error {
	cfg, err := initServerRunConfig()
	if err != nil {
		return fmt.Errorf("can't init server config: %w", err)
	}

	s.httpServer = &http.Server{
		Addr:           cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}
	return nil
}

// Run launches http Server on chosen port with given handler
func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully shuts down Server without interrupting active connections
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
