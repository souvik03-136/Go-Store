package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Function to generate dynamic salt
func generateDynamicSalt() string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(salt)
}

// Combine base secret with dynamic salt
func getSigningSecret(salt string) []byte {
	baseSecret := os.Getenv("JWT_SECRET_KEY")
	combined := baseSecret + salt
	hash := sha256.Sum256([]byte(combined))
	return hash[:]
}

// Generate a JWT with a dynamic component
func GenerateToken(username string) (string, string, error) {
	salt := generateDynamicSalt()

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getSigningSecret(salt))

	return tokenString, salt, err
}

// Validate JWT
func ValidateToken(tokenString string, salt string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getSigningSecret(salt), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
