// internal/services/user_service.go

package services

import (
	"errors"
	"fmt"

	"github.com/souvik03-136/Go-Store/internal/models"
	"github.com/souvik03-136/Go-Store/internal/repository"
)

// UserService handles all user management business logic.
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser validates input, hashes the password, persists the user, and
// returns the created User (PasswordHash is never exposed via JSON).
func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password are all required")
	}

	// Reject duplicate emails up front with a friendly error.
	if existing, _ := s.userRepo.GetUserByEmail(email); existing != nil {
		return nil, errors.New("a user with that email already exists")
	}

	user, err := models.NewUser(generateID(), username, email, password)
	if err != nil {
		return nil, fmt.Errorf("building user model: %w", err)
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("persisting user: %w", err)
	}

	return user, nil
}

// GetUserByID fetches a single user by their UUID.
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}
	return s.userRepo.GetUserByID(id)
}

// UpdateUser applies non-empty field updates to the user identified by id.
// If password is non-empty it is re-hashed before saving.
func (s *UserService) UpdateUser(id, username, email, password string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		if err := user.SetPassword(password); err != nil {
			return nil, fmt.Errorf("updating password: %w", err)
		}
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("persisting user update: %w", err)
	}

	return user, nil
}

// DeleteUser removes the user identified by id.
func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("user id is required")
	}
	return s.userRepo.DeleteUser(id)
}
