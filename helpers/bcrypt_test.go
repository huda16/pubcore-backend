package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPass(t *testing.T) {
	password := "secret123"
	hashed := HashPass(password)

	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestComparePass(t *testing.T) {
	password := "secret123"
	hashed := HashPass(password)

	// Test correct password
	isValid := CompatePass([]byte(hashed), []byte(password))
	assert.True(t, isValid)

	// Test incorrect password
	isInvalid := CompatePass([]byte(hashed), []byte("wrongpass"))
	assert.False(t, isInvalid)
}
