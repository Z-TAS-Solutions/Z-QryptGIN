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

	fmt.Println("Initializing Services...")
	// 7. Initialize Services (continued)
	profileSvc := service.NewProfileService(userRepo)
	notificationSvc := service.NewNotificationService(redisClient)
	sessionSvc := service.NewSessionService(redisClient)
	adminDashboardSvc := service.NewAdminDashboardService()
	adminUserSvc := service.NewAdminUserService(userRepo)
	adminDeviceSvc := service.NewAdminDeviceService()
	adminAuthSvc := service.NewAdminAuthService()
	adminSettingsSvc := service.NewAdminSettingsService()

	fmt.Println("Initializing Handlers...")
	// 8. Initialize Handlers
	// userHandler := handlers.NewUserHandler(userSvc)
	userRegistrationHandler := handlers.NewUserRegistrationHandler(userRegistrationSvc)
	webauthnHandler := handlers.NewWebAuthnHandlerWithRegistration(logger, webauthnSvc, userRepo, credentialRepo, webauthnSessionCache, userRegistrationSvc, redisClient, jwtService)
	userAccountHandler := handlers.NewUserAccountHandler(logger)
	profileHandler := handlers.NewProfileHandler(profileSvc)
	notificationHandler := handlers.NewNotificationHandler(notificationSvc)
	sessionHandler := handlers.NewSessionHandler(sessionSvc)
	adminDashboardHandler := handlers.NewAdminDashboardHandler(adminDashboardSvc)
	adminUserHandler := handlers.NewAdminUserHandler(adminUserSvc)
	adminDeviceHandler := handlers.NewAdminDeviceHandler(adminDeviceSvc)
	adminAuthHandler := handlers.NewAdminAuthHandler(adminAuthSvc)
	adminSettingsHandler := handlers.NewAdminSettingsHandler(adminSettingsSvc)
	adminSecurityHandler := handlers.NewAdminSecurityHandler(adminSettingsSvc)

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

	// Admin Routes
	admin := router.Group("/api/v1/admin")
	{
		dash := admin.Group("/dashboard")
		{
			dash.GET("/analytics", adminDashboardHandler.GetAnalytics)
			dash.GET("/auth-trends", adminDashboardHandler.GetAuthTrends)
			dash.GET("/recent-auth-activity", adminDashboardHandler.GetRecentAuthActivity)
		}

		users := admin.Group("/users")
		{
			users.GET("", adminUserHandler.ListUsers)
			users.GET("/:userId", adminUserHandler.GetUser)
			users.PATCH("/:userId/lock-status", adminUserHandler.UpdateLockStatus)
			users.DELETE("/:userId", adminUserHandler.DeleteUser)
		}

		devices := admin.Group("/devices")
		{
			devices.GET("", adminDeviceHandler.ListDevices)
			devices.POST("/:deviceId/force-logout", adminDeviceHandler.ForceLogout)
		}

		security := admin.Group("/security")
		{
			security.PATCH("/mfa-enforcement", adminSecurityHandler.EnforceMFA)
		}

		authLogs := admin.Group("/auth")
		{
			authLogs.GET("/logs", adminAuthHandler.GetAuthLogs)
			authLogs.GET("/analytics", adminAuthHandler.GetAuthAnalytics)
			// Admin Public Auth (login/refresh) usually doesn't require Bearer token
			authLogs.POST("/login", adminAuthHandler.Login)
			authLogs.POST("/refresh", adminAuthHandler.Refresh)
		}

		settings := admin.Group("/settings")
		{
			settings.GET("", adminSettingsHandler.GetSettings)
		}
	}

	// User Routes
	user := router.Group("/api/v1/user")
	{
		auth := user.Group("/auth")
		{
			register := auth.Group("/register")
			{
				register.POST("/options", nil) // userAuthHandler.RegisterOptions
				register.POST("/verify", nil) // userAuthHandler.RegisterVerify
			}

			login := auth.Group("/login")
			{
				login.POST("/options", nil) // userAuthHandler.LoginOptions
				login.POST("/verify", nil) // userAuthHandler.LoginVerify
			}

			mfa := auth.Group("/mfa")
			{
				mfa.POST("/send", nil) // userMfaHandler.Send
				mfa.POST("/respond", nil) // userMfaHandler.Respond
			}
		}

		// Profile
		profile := user.Group("/profile", server.RequireAuth(jwtService, sessionRepo))
		{
			profile.GET("", profileHandler.GetProfile)
			profile.PATCH("", profileHandler.UpdateProfile)
		}

		// Notifications
		notifications := user.Group("/notifications", server.RequireAuth(jwtService, sessionRepo))
		{
			notifications.GET("", notificationHandler.FetchNotifications)
			notifications.PATCH("/read-all", notificationHandler.MarkAllRead)
			notifications.PATCH("/:notificationId/status", notificationHandler.UpdateStatus)
		}

		// Sessions
		sessions := user.Group("/sessions", server.RequireAuth(jwtService, sessionRepo))
		{
			sessions.GET("", sessionHandler.FetchActiveSessions)
			sessions.POST("/logout-others", sessionHandler.SignOutOtherDevices)
		}

		// External User Account Functions (Protected - requires valid JWT)
		user.POST("/force-logout-devices", server.RequireAuth(jwtService, sessionRepo), userAccountHandler.ForceLogoutAllDevices)
		user.DELETE("/account/delete", server.RequireAuth(jwtService, sessionRepo), userAccountHandler.DeleteAccount)
	}

	fmt.Println("Starting the server...")
	// 10. Finally Starting the damn server 🥲
	srv := server.NewServer(log.Logger, router, cfg)
	srv.Serve()
}
