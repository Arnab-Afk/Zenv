// main.go
package main

import (
	"log"

	"github.com/arnab-afk/Zenv/config"
	"github.com/arnab-afk/Zenv/handlers"
	"github.com/arnab-afk/Zenv/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	cfg := config.Load()

	// Setup router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.RateLimit())
	r.Use(middleware.Authentication())

	// Routes
	api := r.Group("/api/v1")
	{
		secrets := api.Group("/secrets")
		{
			secrets.POST("/", handlers.CreateSecret)
			secrets.GET("/:id", handlers.GetSecret)
			secrets.PUT("/:id", handlers.UpdateSecret)
			secrets.DELETE("/:id", handlers.DeleteSecret)
		}

		keys := api.Group("/keys")
		{
			keys.POST("/rotate", handlers.RotateKeys)
		}
	}

	log.Fatal(r.Run(":8080"))
}
