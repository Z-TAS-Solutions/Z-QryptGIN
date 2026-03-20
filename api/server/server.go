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

func NewServer(l zerolog.Logger, router *gin.Engine, config *configs.Config) *Server {
	return &Server{l: l, router: router, config: config}
}

// Serve creates a new http.Server with support for graceful shutdown
func (s *Server) Serve() {
	srv := &http.Server{
		Addr:    s.config.Server.Address,
		Handler: s.router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.l.Fatal().Err(err).Msg("listen")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 30 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.l.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.l.Fatal().Err(err).Msg("Server Shutdown")
	}
	// catching ctx.Done(). timeout of 30 seconds.

	<-ctx.Done()
	s.l.Info().Msg("Server shutdown timeout of 30 seconds")
	s.l.Info().Msg("Server exiting")
}
