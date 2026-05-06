package helpers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	id := uint(1)
	email := "test@example.com"

	token := GenerateToken(id, email)
	assert.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	id := uint(1)
	email := "test@example.com"
	token := GenerateToken(id, email)

	// Setup Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	claims, err := VerifyToken(c)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	mapClaims := claims.(jwt.MapClaims)
	assert.Equal(t, float64(id), mapClaims["id"])
	assert.Equal(t, email, mapClaims["email"])
}
