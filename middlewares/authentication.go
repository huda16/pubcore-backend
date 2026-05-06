package middlewares

import (
	"net/http"
	"pubcore/helpers"

	"github.com/gin-gonic/gin"
)

// Authentication is a middleware that validates the Bearer JWT token.
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for OPTIONS requests (preflight)
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		verifyToken, err := helpers.VerifyToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
			return
		}
		c.Set("userData", verifyToken)
		c.Next()
	}
}