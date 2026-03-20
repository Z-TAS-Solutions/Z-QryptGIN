package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/configs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	l      zerolog.Logger
	router *gin.Engine
	config *configs.Config
}

// NewServer initializes the server and applies global middleware
func NewServer(l zerolog.Logger, router *gin.Engine, config *configs.Config) *Server {
	// Apply the ErrorHandler middleware globally to the router instance
	// This ensures every request passing through this router is governed
	// by our centralized error-to-JSON logic.
	router.Use(ErrorHandler())

	// If you have a Recovery middleware (highly recommended for production), 
	// it should generally be applied here as well to prevent server crashes.
	router.Use(gin.Recovery())

	return &Server{
		l:      l,
		router: router,
		config: config,
	}
}

// Serve creates a new http.Server with support for graceful shutdown
func (s *Server) Serve() {
	srv := &http.Server{
		Addr:         s.config.Server.Address,
		Handler:      s.router.Handler(),
		ReadTimeout:  10 * time.Second, // Production-ready timeouts
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Run server in a goroutine so that it doesn't block the graceful shutdown cleanup
	go func() {
		s.l.Info().Str("address", s.config.Server.Address).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.l.Fatal().Err(err).Msg("Failed to listen and serve")
		}
	}()

	// Signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C), SIGTERM (Kubernetes/Docker shutdown)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// Block until a signal is received
	sig := <-quit
	s.l.Info().Str("signal", sig.String()).Msg("Shutdown signal received")

	// Context with timeout for the shutdown process
	// This gives active connections 30 seconds to finish before being forced closed
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.l.Info().Msg("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		s.l.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	// Check if the context timed out or was completed successfully
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			s.l.Warn().Msg("Server shutdown timed out (30s limit exceeded)")
		} else {
			s.l.Info().Msg("Server exiting cleanly")
		}
	}
}