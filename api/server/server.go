package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/handlers"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/configs"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
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

	// Mock DB instance for DI setup (assume initialized somewhere else actually)
	var db *gorm.DB

	// Initialization of layers based on api_creation_guidelines.md step 6

	// Utilities
	utilityHandler := handlers.NewUtilityHandler()

	// Profile
	userRepo := repository.NewUserRepository(db)
	profileService := service.NewProfileService(userRepo)
	profileHandler := handlers.NewProfileHandler(profileService)

	// User Notifications
	userNotificationsRepo := repository.NewUserNotificationsRepository(db)
	userNotificationsSvc := service.NewUserNotificationsService(userNotificationsRepo)
	userNotificationsHandler := handlers.NewUserNotificationsHandler(userNotificationsSvc)

	// User Sessions
	userSessionsRepo := repository.NewUserSessionsRepository(db)
	userSessionsSvc := service.NewUserSessionsService(userSessionsRepo)
	userSessionsHandler := handlers.NewUserSessionsHandler(userSessionsSvc)

	// User Account
	userAccountRepo := repository.NewUserAccountRepository(db)
	userAccountSvc := service.NewUserAccountService(userAccountRepo)
	userAccountHandler := handlers.NewUserAccountHandler(userAccountSvc)

	// Dashboard
	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardSvc := service.NewDashboardService(dashboardRepo)
	dashboardHandler := handlers.NewDashboardHandler(dashboardSvc)

	// Admin Users
	adminUserRepo := repository.NewAdminUserRepository(db)
	adminUserSvc := service.NewAdminUserService(adminUserRepo)
	adminUserHandler := handlers.NewAdminUserHandler(adminUserSvc)

	// Admin Devices
	adminDeviceRepo := repository.NewAdminDeviceRepository(db)
	adminDeviceSvc := service.NewAdminDeviceService(adminDeviceRepo)
	adminDeviceHandler := handlers.NewAdminDeviceHandler(adminDeviceSvc)

	// Admin Security
	adminSecurityRepo := repository.NewAdminSecurityRepository(db)
	adminSecuritySvc := service.NewAdminSecurityService(adminSecurityRepo)
	adminSecurityHandler := handlers.NewAdminSecurityHandler(adminSecuritySvc)

	// Admin Auth Logs
	adminAuthLogsRepo := repository.NewAdminAuthLogsRepository(db)
	adminAuthLogsSvc := service.NewAdminAuthLogsService(adminAuthLogsRepo)
	adminAuthLogsHandler := handlers.NewAdminAuthLogsHandler(adminAuthLogsSvc)

	// Admin Auth
	adminAuthRepo := repository.NewAdminAuthRepository(db)
	adminAuthSvc := service.NewAdminAuthService(adminAuthRepo)
	adminAuthHandler := handlers.NewAdminAuthHandler(adminAuthSvc)

	// Admin Settings
	adminSettingsRepo := repository.NewAdminSettingsRepository(db)
	adminSettingsSvc := service.NewAdminSettingsService(adminSettingsRepo)
	adminSettingsHandler := handlers.NewAdminSettingsHandler(adminSettingsSvc)

	// User Auth (Passkey/WebAuthn)
	userAuthRepo := repository.NewUserAuthRepository(db)
	userAuthSvc := service.NewUserAuthService(userAuthRepo)
	userAuthHandler := handlers.NewUserAuthHandler(userAuthSvc)

	// User MFA
	userMfaRepo := repository.NewUserMfaRepository(db)
	userMfaSvc := service.NewUserMfaService(userMfaRepo)
	userMfaHandler := handlers.NewUserMfaHandler(userMfaSvc)

	// Register Routes (V1)
	v1 := router.Group("/api/v1")
	{
		// Utility
		v1.GET("/ping", utilityHandler.Ping)

		// Admin Dashboard
		admin := v1.Group("/admin")
		// In reality, add role middleware here: admin.Use(AdminAuthMiddleware())
		{
			dash := admin.Group("/dashboard")
			dash.GET("/analytics", dashboardHandler.GetAnalytics)
			dash.GET("/auth-trends", dashboardHandler.GetAuthTrends)
			dash.GET("/recent-auth-activity", dashboardHandler.GetRecentAuthActivity)

			users := admin.Group("/users")
			users.GET("", adminUserHandler.ListUsers)
			users.GET("/:userId", adminUserHandler.GetUser)
			users.PATCH("/:userId/lock-status", adminUserHandler.UpdateLockStatus)
			users.DELETE("/:userId", adminUserHandler.DeleteUser)

			devices := admin.Group("/devices")
			devices.GET("", adminDeviceHandler.ListDevices)
			devices.POST("/:deviceId/force-logout", adminDeviceHandler.ForceLogout)

			security := admin.Group("/security")
			security.PATCH("/mfa-enforcement", adminSecurityHandler.EnforceMfa)

			authLogs := admin.Group("/auth")
			authLogs.GET("/logs", adminAuthLogsHandler.GetAuthLogs)
			authLogs.GET("/analytics", adminAuthLogsHandler.GetAuthAnalytics)

			// Admin Public Auth (login/refresh) usually doesn't require Bearer token
			authLogs.POST("/login", adminAuthHandler.Login)
			authLogs.POST("/refresh", adminAuthHandler.Refresh)

			settings := admin.Group("/settings")
			settings.GET("", adminSettingsHandler.GetSettings)
		}

		user := v1.Group("/user")
		{
			auth := user.Group("/auth")
			{
				register := auth.Group("/register")
				register.POST("/options", userAuthHandler.RegisterOptions)
				register.POST("/verify", userAuthHandler.RegisterVerify)

				login := auth.Group("/login")
				login.POST("/options", userAuthHandler.LoginOptions)
				login.POST("/verify", userAuthHandler.LoginVerify)

				mfa := auth.Group("/mfa")
				mfa.POST("/send", userMfaHandler.Send)
				mfa.POST("/respond", userMfaHandler.Respond)
			}

			// Profile
			profile := user.Group("/profile")
			{
				profile.GET("", profileHandler.GetProfile)
				profile.PATCH("", profileHandler.UpdateProfile)
			}

			// Notifications
			notifications := user.Group("/notifications")
			{
				notifications.GET("", userNotificationsHandler.FetchNotifications)
				notifications.PATCH("/read-all", userNotificationsHandler.MarkAllRead) // Registered before parameterized routes
				notifications.PATCH("/:notificationId/status", userNotificationsHandler.UpdateStatus)
			}

			// Sessions
			sessions := user.Group("/sessions")
			{
				sessions.GET("", userSessionsHandler.FetchActiveSessions)
				sessions.POST("/logout-others", userSessionsHandler.SignOutOtherDevices)
			}

			// External User Account Functions
			user.POST("/force-logout-devices", userAccountHandler.ForceLogoutAllDevices)
			user.DELETE("/account/delete", userAccountHandler.DeleteAccount)
		}
	}

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
