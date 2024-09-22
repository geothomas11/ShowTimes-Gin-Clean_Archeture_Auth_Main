package middleware

import (
	"ShowTimes/pkg/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserAuthMiddleware(c *gin.Context) {
	tokenstring := c.GetHeader("Authorization")
	if tokenstring == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization"})
		c.Abort()
		return
	}

	accessToken := strings.TrimPrefix(tokenstring, "Bearer ")
	cfg, _ := config.LoadConfig()
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.User_AccessKey), nil
	})

	if err != nil {
		// The access token is invalid.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		c.Abort()
		return

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println("claims2", claims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthoraised access 1"})
		c.Abort()
		return
	}
	fmt.Println("claims", claims)

	role, ok := claims["role"].(string)
	if !ok || role != "Client" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised access 2"})
		c.Abort()
		return
	}
	id, ok := claims["id"].(float64)
	if !ok || id == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "error in retreving id"})
	}

	c.Set("role", role)
	c.Set("id", int(id))

	c.Next()
}
