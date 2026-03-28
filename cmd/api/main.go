package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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

	fmt.Println("Initializing Logger...")
	// 1. Initialize Logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	fmt.Println("Loading Config...")
	// 2. Load Configurations
	cfg := configs.NewConfig()

	fmt.Println("Connecting to PostgreSQL database...")
	// 3. Connect to PostgreSQL Database
	db, err := database.NewDatabaseConnection(cfg.Database.DatabaseSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}

	fmt.Println("Connecting to Redis cache...")
	// 4. Connect to Redis Chache
	redisClient, err := cache.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to Redis")
	}

	defer redisClient.Close()

	fmt.Println("Starting Email Service...")
	// 5. Start the Email Service
	logger.Info().Msg("Initializing Email Service...")
	emailSvc, err := service.NewEmailService(
		ctx,
		cfg.Gmail.ClientID,
		cfg.Gmail.ClientSecret,
		cfg.Gmail.TokenPath,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize Email Service")
	}

	fmt.Println("Initializing Repositories...")
	// 6. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	credentialRepo := repository.NewWebAuthnCredentialRepository(db)
	// sessionRepo := repository.NewSessionRepository(redisClient)

	fmt.Println("Initializing Services...")
	// 7. Initialize Services
	// userSvc := service.NewUserService(userRepo, sessionRepo, emailSvc)

	fmt.Println("Initializing WebAuthn...")
	// 7.5 Initialize WebAuthn for Passkey Registration/Authentication
	wa, err := service.InitWebAuthn()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize WebAuthn")
	}
	webauthnSvc := service.NewWebAuthnService(wa)
	webauthnSessionCache := cache.NewWebAuthnSessionCache(redisClient)

	// Create user registration service with WebAuthn support for complete registration flow
	userRegistrationSvc := service.NewUserRegistrationServiceWithWebAuthn(userRepo, credentialRepo, redisClient, emailSvc, webauthnSvc, db)

	fmt.Println("Initializing JWT Service...")
	// 7.6 Initialize JWT Service with EdDSA signing (hybrid stateless + stateful via Redis)
	// Generate EdDSA keys for token signing
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate EdDSA keys for JWT signing")
	}

	// Initialize JWT Service with Redis for session tracking
	jwtService := service.NewJWTService(privateKey, publicKey, "Z-QryptGIN", redisClient)
	logger.Info().Msg("JWT Service initialized with EdDSA algorithm and Redis session tracking (Hybrid Approach)")

	fmt.Println("Initializing Handlers...")
	// 8. Initialize Handlers
	// userHandler := handlers.NewUserHandler(userSvc)
	userRegistrationHandler := handlers.NewUserRegistrationHandler(userRegistrationSvc)
	webauthnHandler := handlers.NewWebAuthnHandlerWithRegistration(logger, webauthnSvc, userRepo, credentialRepo, webauthnSessionCache, userRegistrationSvc, redisClient, jwtService)

	fmt.Println("Setting up Gin Router...")
	// 9. Setup Gin Router
	router := gin.Default()
	router.Use(cfg.CorsNew())

	fmt.Println("Configuring API Routes...")
	// API Routes
	v1 := router.Group("/api/v1/admin")
	{
		// v1.POST("/users/RegisterUser", userHandler.Register)
		v1.POST("/users/register/new", userRegistrationHandler.Register)
		v1.POST("/users/register/verifyOTP", userRegistrationHandler.VerifyOTP)
		v1.POST("/users/register/resendOTP", userRegistrationHandler.ResendOTP)
	}

	// WebAuthn Routes (Passkey Registration & Authentication)
	webauthn := router.Group("/api/v1/webauthn")
	{
		// Passkey Registration Endpoints
		webauthn.POST("/register/begin", webauthnHandler.RegisterBegin)
		webauthn.POST("/register/finish", webauthnHandler.RegisterFinish)

		// Passkey Authentication Endpoints
		webauthn.POST("/login/begin", webauthnHandler.LoginBegin)
		webauthn.POST("/login/finish", webauthnHandler.LoginFinish)
	}

	fmt.Println("Starting the server...")
	// 10. Finally Starting the damn server 🥲
	srv := server.NewServer(log.Logger, router, cfg)
	srv.Serve()
}
