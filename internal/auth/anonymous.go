// internal/auth/anonymous.go

package auth

import "github.com/google/uuid"

// GenerateAnonymousID returns a new random UUID string to identify
// an anonymous (unauthenticated) session.
func GenerateAnonymousID() string {
	return "anon_" + uuid.New().String()
}
