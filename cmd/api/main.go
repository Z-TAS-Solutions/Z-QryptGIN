package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	
	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/handlers"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/server"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/configs"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

func main() {
	// 1. Load Config
	cfg := configs.NewConfig()

	// 2. Connect to DB
	db, err := database.NewDatabaseConnection(cfg.Database.DatabaseSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)

	// 4. Initialize Services
	userSvc := service.NewUserService(userRepo)

	// 5. Initialize Handlers
	userHandler := handlers.NewUserHandler(userSvc)

	// 6. Setup Gin Router
	router := gin.Default()
	router.Use(cfg.CorsNew())
	
	// API Routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", userHandler.Register)
	}

	// 7. Start Server using your existing Server struct
	srv := server.NewServer(log.Logger, router, cfg)
	srv.Serve()
}