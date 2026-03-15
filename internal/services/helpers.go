// internal/services/helpers.go

package services

import "github.com/google/uuid"

// generateID returns a new random UUID string.
func generateID() string {
	return uuid.New().String()
}
