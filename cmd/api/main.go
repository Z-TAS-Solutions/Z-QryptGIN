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
	sessionRepo := repository.NewSessionRepository(redisClient)
	notificationRepo := repository.NewNotificationRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	fmt.Println("Initializing Services...")
	// 7. Initialize Services
	sessionSvc := service.NewSessionService(sessionRepo)
	notificationSvc := service.NewNotificationService(notificationRepo)
	dashboardSvc := service.NewDashboardService(dashboardRepo)

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
	userHandler := handlers.NewUserHandler(sessionSvc, notificationSvc)
	sessionHandler := handlers.NewSessionHandler(sessionSvc)
	dashboardHandler := handlers.NewDashboardHandler(dashboardSvc)
	userRegistrationHandler := handlers.NewUserRegistrationHandler(userRegistrationSvc)
	webauthnHandler := handlers.NewWebAuthnHandlerWithRegistration(logger, webauthnSvc, userRepo, credentialRepo, webauthnSessionCache, userRegistrationSvc, redisClient, jwtService)

	fmt.Println("Setting up Gin Router...")
	// 9. Setup Gin Router
	router := gin.Default()
	router.Use(cfg.CorsNew())

	fmt.Println("Configuring API Routes...")
	// API Routes
	admin := router.Group("/api/v1/admin")
	{
		// User registration endpoints (no auth required during registration)
		admin.POST("/users/register/new", userRegistrationHandler.Register)
		admin.POST("/users/register/verifyOTP", userRegistrationHandler.VerifyOTP)
		admin.POST("/users/register/resendOTP", userRegistrationHandler.ResendOTP)

		// Admin dashboard endpoints (require authentication and admin role)
		protected := admin.Group("")
		protected.Use(server.RequireAuth(jwtService, sessionRepo))
		protected.Use(server.RequireRole("Admin"))
		protected.Use(server.ValidateRoleConsistency(jwtService, sessionRepo))
		{
			// Authentication trends endpoint
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/auth-trends", dashboardHandler.GetAuthenticationTrends)
				dashboard.GET("/metrics", dashboardHandler.GetDashboardMetrics)
			}
			mfa_admin := protected.Group("/mfa")
			{
				mfa_admin.GET("/enforce", nil) // adminMfaHandler.GetPendingMfaRequests
			}
		}
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

	// User Routes
	user := router.Group("/api/v1/user")
	{
		// Protected routes (require authentication)
		protected := user.Group("")
		protected.Use(server.RequireAuth(jwtService, sessionRepo))
		{
			// Notifications endpoint
			protected.GET("/notifications", userHandler.GetNotifications)
			protected.PATCH("/notifications/:notificationId/status", userHandler.UpdateNotificationStatus)
			protected.PATCH("/notifications/read-all", userHandler.MarkAllAsRead)

			// Session management routes
			protected.GET("/sessions", sessionHandler.GetActiveSessions)
			protected.POST("/sessions/logout-others", sessionHandler.LogoutOthers)

			// Dashboard routes (legacy - keeping for backward compatibility)
			dashboard := protected.Group("/dashboard")
			{
				session := dashboard.Group("/session")
				{
					session.GET("/activeSessions", userHandler.GetActiveSessions)
				}
			}
		}

		// Auth routes (non-protected)
		auth := user.Group("/auth")
		{
			mfa := auth.Group("/mfa")
			{
				mfa.POST("/send", nil)    // userMfaHandler.Send
				mfa.POST("/respond", nil) // userMfaHandler.Respond
			}
		}
	}

	fmt.Println("Starting the server...")
	// 10. Finally Starting the damn server 🥲
	srv := server.NewServer(log.Logger, router, cfg)
	srv.Serve()
}
