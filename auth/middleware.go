package auth

import (
	"net/http"

	"github.com/arnab-afk/Zenv/database"
	"github.com/arnab-afk/Zenv/security"
	"github.com/gin-gonic/gin"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		decryptedToken, err := security.Decrypt(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		storedToken, err := getAccessToken(decryptedToken)
		if err != nil || storedToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getAccessToken(token string) (string, error) {
	var storedToken string
	err := database.DB.QueryRow("SELECT token FROM access_tokens WHERE token = $1", token).Scan(&storedToken)
	if err != nil {
		return "", err
	}
	return storedToken, nil
}
