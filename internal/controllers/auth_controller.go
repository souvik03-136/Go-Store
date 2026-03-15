// internal/controllers/auth_controller.go

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/merrors"
	"github.com/souvik03-136/Go-Store/internal/services"
)

// AuthController handles HTTP requests for authentication endpoints.
type AuthController struct {
	authSvc *services.AuthService
}

// NewAuthController creates a new AuthController.
func NewAuthController(authSvc *services.AuthService) *AuthController {
	return &AuthController{authSvc: authSvc}
}

// registerOAuthRequest is the expected JSON body for OAuth registration.
type registerOAuthRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// loginOAuthRequest is the expected JSON body for OAuth login.
type loginOAuthRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// validateTokenRequest holds the query parameters for token validation.
type validateTokenRequest struct {
	Token string `form:"token" binding:"required"`
	Salt  string `form:"salt"  binding:"required"`
}

// RegisterOAuthUser handles POST /v1/auth/oauth/register
func (c *AuthController) RegisterOAuthUser(ctx *gin.Context) {
	var req registerOAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		merrors.Validation(ctx, err.Error())
		return
	}

	user, pair, err := c.authSvc.RegisterOAuth(req.Username, req.Email, req.Password)
	if err != nil {
		merrors.Conflict(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": pair.Token,
		"salt":  pair.Salt,
	})
}

// LoginOAuthUser handles POST /v1/auth/oauth/login
func (c *AuthController) LoginOAuthUser(ctx *gin.Context) {
	var req loginOAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		merrors.Validation(ctx, err.Error())
		return
	}

	user, pair, err := c.authSvc.LoginOAuth(req.Email, req.Password)
	if err != nil {
		merrors.Unauthorized(ctx, "invalid email or password")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": pair.Token,
		"salt":  pair.Salt,
	})
}

// RegisterAnonymousUser handles POST /v1/auth/anonymous/register
func (c *AuthController) RegisterAnonymousUser(ctx *gin.Context) {
	anonID, pair, err := c.authSvc.RegisterAnonymous()
	if err != nil {
		merrors.InternalServer(ctx, "failed to create anonymous session")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"anonymous_id": anonID,
		"token":        pair.Token,
		"salt":         pair.Salt,
	})
}

// LogoutUser handles POST /v1/auth/logout
// Token invalidation is stateless here; integrate a token blacklist (e.g. Redis)
// if you need server-side revocation.
func (c *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// ValidateTokenHandler handles GET /v1/auth/validate
func (c *AuthController) ValidateTokenHandler(ctx *gin.Context) {
	var req validateTokenRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		merrors.BadRequest(ctx, "token and salt query parameters are required")
		return
	}

	subject, err := c.authSvc.ValidateToken(req.Token, req.Salt)
	if err != nil {
		merrors.Unauthorized(ctx, "invalid or expired token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"subject": subject, "valid": true})
}
