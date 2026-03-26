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
	// 1. Load the damn Logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()
	// 2. Loading the Config
	cfg := configs.NewConfig()

	// 3. Connect to PostgreSQL database
	db, err := database.NewDatabaseConnection(cfg.Database.DatabaseSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}

	// 4. Connect to Redis Chache 
	redisClient, err := cache.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to Redis")
	}

	defer redisClient.Close()

	// 5. Start the Email Service
	logger.Info().Msg("Initializing Email Service...")
	emailSvc, err := service.NewEmailService(
		cts,
		cfg.Gmail.ClientID,
		cfg.Gmail.ClientSecret,
		cfg.Gmail.TokenPath,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize Email Service")
	}

	// 6. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(redisClient)

	// 7. Initialize Services
	userSvc := service.NewUserService(userRepo, sessionRepo)

	// 8. Initialize Handlers
	userHandler := handlers.NewUserHandler(userSvc)

	// 9. Setup Gin Router
	router := gin.Default()
	router.Use(cfg.CorsNew())

	// API Routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", userHandler.Register)
	}

	// 10. Finally Starting the damn server 🥲
	srv := server.NewServer(log.Logger, router, cfg)
	srv.Serve()
}
