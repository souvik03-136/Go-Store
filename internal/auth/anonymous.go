package auth

import "github.com/google/uuid"

func GenerateAnonymousID() string {
	return uuid.New().String()
}
