package middleware

import (
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

	tokenstring = strings.TrimPrefix(tokenstring, "Bearer")

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte("comebywatch"), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization Token"})
		c.Abort()
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthoraised access"})
		c.Abort()
		return
	}
	fmt.Println("claims", claims)

	role, ok := claims["role"].(string)
	if !ok || role != "client" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised access"})
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
