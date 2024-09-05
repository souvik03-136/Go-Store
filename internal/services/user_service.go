package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/souvik03-136/Go-Store/internal/merrors"
	"github.com/souvik03-136/Go-Store/internal/models"
)

// UserService handles the business logic for user operations.
type UserService struct{}

// NewUserService creates a new instance of UserService.
func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser creates a new user record.
func (s *UserService) CreateUser(ctx *gin.Context, username, email, password string) (*models.User, error) {
	// Validate input
	if username == "" || email == "" || password == "" {
		merrors.BadRequest(ctx, "Username, email, and password are required")
		return nil, errors.New("username, email, and password are required")
	}

	// Generate user ID
	userID := uuid.New().String()

	// Create the user model
	user, err := models.NewUser(userID, username, email, password)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to create user")
		return nil, err
	}

	// In a real application, you'd save the user in a database here
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(ctx *gin.Context, userID string) (*models.User, error) {
	// Simulate a function to retrieve user by ID (you'll implement the actual DB query)
	user, err := models.GetUserByID(userID) // You need to implement this in models/user.go
	if err != nil {
		merrors.InternalServer(ctx, "Failed to retrieve user by ID")
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser updates the user's information.
func (s *UserService) UpdateUser(ctx *gin.Context, user *models.User, username, email, password string) (*models.User, error) {
	if username == "" && email == "" && password == "" {
		merrors.BadRequest(ctx, "No updates provided")
		return nil, errors.New("no updates provided")
	}

	// Update user model
	err := user.UpdateUser(username, email, password)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to update user")
		return nil, err
	}

	// In a real application, you'd save the updated user in the database here
	return user, nil
}

// DeleteUser removes a user by their ID.
func (s *UserService) DeleteUser(ctx *gin.Context, user *models.User) error {
	// In a real application, you'd delete the user from the database here
	// Here, we're just simulating the deletion by clearing user data
	user.Username = ""
	user.Email = ""
	user.Password = ""
	user.UpdatedAt = time.Now()

	// In a real application, ensure database deletion
	return nil
}
