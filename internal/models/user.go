package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User instance.
func NewUser(id, username, email, password string) (*User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// UpdateUser updates the user information.
func (u *User) UpdateUser(username, email, password string) error {
	u.Username = username
	u.Email = email
	if password != "" {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	u.UpdatedAt = time.Now()
	return nil
}

// CheckPassword checks if the provided password matches the user's hashed password.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// hashPassword hashes a plain text password.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ResetPassword allows the user to reset their password.
func (u *User) ResetPassword(newPassword string) error {
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// GetUserByID retrieves a user by their ID (simulated).
func GetUserByID(userID string) (*User, error) {
	// Simulate a DB lookup (you can replace this with real DB code)
	// For now, we'll just return a dummy user for illustration purposes
	if userID == "valid-id" {
		return &User{
			ID:        userID,
			Username:  "testuser",
			Email:     "testuser@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	return nil, errors.New("user not found")
}
