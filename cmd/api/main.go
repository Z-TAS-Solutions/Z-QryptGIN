package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/handlers"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/server"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/configs"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/cache"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

func main() {
	// Load Config
	cfg := configs.NewConfig()

	// Connect to PostgreSQL
	db, err := database.NewDatabaseConnection(cfg.Database.DatabaseSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}

	redisClient, err := cache.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to Redis")
	}

	defer redisClient.Close()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(redisClient)

	// Initialize Services
	userSvc := service.NewUserService(userRepo, sessionRepo)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userSvc)

	// Setup Gin Router
	router := gin.Default()
	router.Use(cfg.CorsNew())

	// API Routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", userHandler.Register)
	}

	// Start Server using your existing Server struct
	srv := server.NewServer(log.Logger, router, cfg)
	srv.Serve()
}
