package routes

import (
	"net/http"

	"github.com/arnab-afk/Zenv/database"
	"github.com/arnab-afk/Zenv/monitoring"
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
	// ...existing code...
}

func getSecret(c *gin.Context) {
	userID := c.GetUint("userID")
	secretID := c.Param("id")

	var secret database.Secret
	if err := database.DB.Where("id = ? AND user_id = ?", secretID, userID).First(&secret).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
		return
	}

	decrypted, err := security.DecryptSecret(secret.Value)
	if err != nil {
		monitoring.LogError(err, map[string]interface{}{"secret_id": secretID})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      secret.ID,
		"name":    secret.Name,
		"value":   string(decrypted),
		"version": secret.Version,
	})
}

func updateSecret(c *gin.Context) {
	userID := c.GetUint("userID")
	secretID := c.Param("id")

	var input struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var secret database.Secret
	if err := database.DB.Where("id = ? AND user_id = ?", secretID, userID).First(&secret).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
		return
	}

	encrypted := security.EncryptSecret([]byte(input.Value))

	secret.Name = input.Name
	secret.Value = encrypted
	secret.Version++

	database.DB.Save(&secret)

	monitoring.AuditLog(string(userID), "update_secret", secretID, true)
	c.JSON(http.StatusOK, gin.H{"id": secret.ID, "version": secret.Version})
}

func deleteSecret(c *gin.Context) {
	userID := c.GetUint("userID")
	secretID := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", secretID, userID).Delete(&database.Secret{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
		return
	}

	monitoring.AuditLog(string(userID), "delete_secret", secretID, true)
	c.JSON(http.StatusOK, gin.H{"message": "Secret deleted"})
}
