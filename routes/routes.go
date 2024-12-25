package routes

import (
	"github.com/arnab-afk/Zenv/database"
	"github.com/arnab-afk/Zenv/security"
	"github.com/gin-gonic/gin"
)

func SetupSecretRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/secrets", createSecret)
		api.GET("/secrets/:id", getSecret)
		api.PUT("/secrets/:id", updateSecret)
		api.DELETE("/secrets/:id", deleteSecret)
	}
}

func createSecret(c *gin.Context) {
	var input struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	encrypted := security.EncryptSecret([]byte(input.Value))

	secret := database.Secret{
		UserID:  c.GetUint("userID"),
		Name:    input.Name,
		Value:   encrypted,
		Version: 1,
	}

	database.DB.Create(&secret)
	c.JSON(201, gin.H{"id": secret.ID})
}
