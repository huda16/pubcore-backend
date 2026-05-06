package helpers

import (
	"errors"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "pubcore-default-secret" // fallback for development only
	}
	return []byte(secret)
}

// GenerateToken creates a signed JWT token for the given user ID and email.
func GenerateToken(id uint, email string) string {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString(getSecretKey())

	return signedToken
}

// VerifyToken validates the Bearer token from the Authorization header.
func VerifyToken(c *gin.Context) (interface{}, error) {
	errResponse := errors.New("sign in to proceed")
	headerToken := c.Request.Header.Get("Authorization")

	if headerToken == "" {
		return nil, errResponse
	}

	parts := strings.Split(headerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errResponse
	}

	stringToken := parts[1]
	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}
		return getSecretKey(), nil
	})

	if err != nil || token == nil || !token.Valid {
		return nil, errResponse
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errResponse
	}

	return claims, nil
}