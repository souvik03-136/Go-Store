// internal/controllers/user_controller.go

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/merrors"
	"github.com/souvik03-136/Go-Store/internal/services"
)

// UserController handles HTTP requests for user management endpoints.
type UserController struct {
	userSvc *services.UserService
}

// NewUserController creates a new UserController.
func NewUserController(userSvc *services.UserService) *UserController {
	return &UserController{userSvc: userSvc}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type updateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"    binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
}

// CreateUser handles POST /v1/users
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		merrors.Validation(ctx, err.Error())
		return
	}

	user, err := c.userSvc.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		merrors.Conflict(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

// GetUserByID handles GET /v1/users/:id
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.userSvc.GetUserByID(id)
	if err != nil {
		merrors.NotFound(ctx, "user not found")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser handles PUT /v1/users/:id
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		merrors.Validation(ctx, err.Error())
		return
	}

	user, err := c.userSvc.UpdateUser(id, req.Username, req.Email, req.Password)
	if err != nil {
		merrors.InternalServer(ctx, "failed to update user")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser handles DELETE /v1/users/:id
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.userSvc.DeleteUser(id); err != nil {
		merrors.NotFound(ctx, "user not found")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
