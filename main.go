package main

import (
	"github.com/arnab-afk/Zenv/auth"
	"github.com/arnab-afk/Zenv/config"
	"github.com/arnab-afk/Zenv/database"
	"github.com/arnab-afk/Zenv/monitoring"
	"github.com/arnab-afk/Zenv/routes"
	"github.com/arnab-afk/Zenv/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize configuration
	config.LoadConfig()

	//Initialize logging
	if err := monitoring.SetupLogging(); err != nil {
		panic("Failed to setup logging" + err.Error())
	}
	defer monitoring.Logger.Sync()

	// Initialize database
	database.InitDB()

	// Initialize logging
	monitoring.SetupLogging()

	// Initialize key rotation scheduler
	security.StartKeyRotationScheduler()

	// Initialize rate limiter
	rateLimiter := security.SetupRateLimiter()

	// Set up Gin router
	router := gin.Default()

	// Apply middlewares
	router.Use(rateLimiter)
	router.Use(auth.AuthMiddleware())
	router.Use(monitoring.RequestLogger())
	router.Use(monitoring.MetricsMiddleware())

	// Setup routes
	routes.SetupSecretRoutes(router)

	// Start the server
	monitoring.Logger.Info("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		monitoring.Logger.Fatal("Server failed to start: ", zap.String("error", err.Error()))
	}
}
