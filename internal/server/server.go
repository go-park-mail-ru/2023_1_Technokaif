package server

import (
	"context"
	"net/http"
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

// Run launches http Server on chosen port with given handler
func (s *Server) Run(host, port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           host + ":" + port,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down Server without interrupting active connections
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
