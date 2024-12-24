package auth

import (
	"net/http"

	"github.com/arnab-afk/Zenv/database"
	"github.com/arnab-afk/Zenv/security"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/auth/token", handleToken)
	router.GET("/secret", authMiddleware(), handleSecret)
}

func handleToken(c *gin.Context) {
	var request struct {
		AccessToken string `json:"access_token"`
		MFA         string `json:"mfa"`
		IP          string `json:"ip"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Verify MFA
	if !security.VerifyMFA(request.MFA) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid MFA"})
		return
	}

	// Verify IP
	if !security.VerifyIP(request.IP) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized IP"})
		return
	}

	// Store the access token in the database
	encryptedToken, err := security.Encrypt(request.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt access token"})
		return
	}
	_, err = database.DB.Exec("INSERT INTO access_tokens (token) VALUES ($1)", encryptedToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access token received"})
}

func handleSecret(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"secret": "This is a protected resource"})
}
