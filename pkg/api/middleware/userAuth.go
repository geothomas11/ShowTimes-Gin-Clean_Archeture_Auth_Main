package middleware

import (
	"ShowTimes/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserAuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization"})
		c.Abort()
		return
	}

	accessToken := strings.TrimPrefix(tokenString, "Bearer ")
	cfg, _ := config.LoadConfig()

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.User_AccessKey), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised access *"})
		c.Abort()
		return
	}

	role, ok := claims["role"].(string)
	if !ok || (role != "Client" && role != "admin") { // Allow both "Client" and "admin" roles
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorised access #"})
		c.Abort()
		return
	}

	id, ok := claims["id"].(float64)
	if !ok || id == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "error in retrieving id"})
		c.Abort()
		return
	}

	c.Set("role", role)
	c.Set("id", int(id))

	c.Next()
}
