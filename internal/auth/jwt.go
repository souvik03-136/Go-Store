// internal/auth/jwt.go

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims holds the JWT standard claims plus our custom fields.
type Claims struct {
	Subject string `json:"sub"`
	jwt.RegisteredClaims
}

// GenerateDynamicSalt creates a cryptographically-random 16-byte salt
// encoded as base64, used to make each token's signing key unique.
func GenerateDynamicSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generating salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// signingKey derives a per-token HMAC key by SHA-256 hashing the base secret
// concatenated with the per-token salt. This means a compromised token cannot
// be used to sign other tokens even if the attacker knows the salt.
func signingKey(baseSecret, salt string) []byte {
	combined := baseSecret + salt
	hash := sha256.Sum256([]byte(combined))
	return hash[:]
}

// GenerateToken produces a signed JWT and the salt used to derive its signing
// key. Both values must be returned to the client; the salt is required to
// validate the token later.
func GenerateToken(baseSecret, subject string, expiry time.Duration) (tokenString, salt string, err error) {
	salt, err = GenerateDynamicSalt()
	if err != nil {
		return "", "", err
	}

	now := time.Now()
	claims := Claims{
		Subject: subject,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(signingKey(baseSecret, salt))
	if err != nil {
		return "", "", fmt.Errorf("signing token: %w", err)
	}

	return tokenString, salt, nil
}

// ValidateToken parses and validates a JWT, returning the claims on success.
// The same salt that was returned by GenerateToken must be provided.
func ValidateToken(baseSecret, tokenString, salt string) (*Claims, error) {
	key := signingKey(baseSecret, salt)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
