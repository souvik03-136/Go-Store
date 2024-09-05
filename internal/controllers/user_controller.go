package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/merrors"
	"github.com/souvik03-136/Go-Store/internal/models"
	"github.com/souvik03-136/Go-Store/internal/repository"
)

type UserController struct {
	userRepo *repository.UserRepository
}

// NewUserController creates a new instance of UserController.
func NewUserController(userRepo *repository.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

// CreateUser handles user registration.
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		merrors.BadRequest(ctx, "Invalid request payload")
		return
	}

	if err := c.userRepo.CreateUser(&user); err != nil {
		merrors.InternalServer(ctx, "Error creating user")
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// GetUserByID handles fetching a user by their ID.
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Query("id")

	if userID == "" {
		merrors.BadRequest(ctx, "User ID is required")
		return
	}

	user, err := c.userRepo.GetUserByID(userID)
	if err != nil {
		merrors.NotFound(ctx, "User not found")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateUser handles updating a user's information.
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Query("id")

	if userID == "" {
		merrors.BadRequest(ctx, "User ID is required")
		return
	}

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		merrors.BadRequest(ctx, "Invalid request payload")
		return
	}

	user.ID = userID
	if err := c.userRepo.UpdateUser(&user); err != nil {
		merrors.InternalServer(ctx, "Error updating user")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser handles the deletion of a user by their ID.
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Query("id")

	if userID == "" {
		merrors.BadRequest(ctx, "User ID is required")
		return
	}

	if err := c.userRepo.DeleteUser(userID); err != nil {
		merrors.InternalServer(ctx, "Error deleting user")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
