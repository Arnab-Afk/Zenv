package main

import (
    "github.com/arnab-afk/Zenv/auth"
    "github.com/arnab-afk/Zenv/config"
    "github.com/arnab-afk/Zenv/database"
    "github.com/arnab-afk/Zenv/security"
    "github.com/arnab-afk/Zenv/monitoring"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize configuration
    config.LoadConfig()

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

    // Apply rate limiter middleware
    router.Use(rateLimiter)

    // Define routes
    auth.SetupRoutes(router)

    // Start the server
    router.Run(":8080")
}