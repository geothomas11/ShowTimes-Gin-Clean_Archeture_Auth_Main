package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AdminAuthMiddleware(c *gin.Context) {
	accessToken := c.Request.Header.Get("Authorization")

	accessToken = strings.TrimPrefix(accessToken, "Bearer")

	_, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("accessScret"), nil
	})
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	c.Next()
}
