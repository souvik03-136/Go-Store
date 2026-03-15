// internal/models/user.go

package models

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents an authenticated user in the system.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never exposed in JSON responses
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser creates and validates a new User, hashing the plain-text password.
func NewUser(id, username, email, plainPassword string) (*User, error) {
	if id == "" || username == "" || email == "" || plainPassword == "" {
		return nil, errors.New("id, username, email, and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	now := time.Now()
	return &User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// CheckPassword returns true if the plain-text password matches the stored hash.
func (u *User) CheckPassword(plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword)) == nil
}

// SetPassword hashes a new plain-text password and stamps UpdatedAt.
func (u *User) SetPassword(plainPassword string) error {
	if plainPassword == "" {
		return errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}
	u.PasswordHash = string(hash)
	u.UpdatedAt = time.Now()
	return nil
}
