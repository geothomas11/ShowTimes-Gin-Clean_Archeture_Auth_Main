package middleware

import (
	"ShowTimes/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AdminAuthMiddleware(c *gin.Context) {

	accessToken := c.Request.Header.Get("Authorization")

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	cfg, _ := config.LoadConfig()
	_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Admin_AccessKey), nil
	})

	if err != nil {
		// The access token is invalid.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		c.Abort()
		return

	}

	c.Next()
}
