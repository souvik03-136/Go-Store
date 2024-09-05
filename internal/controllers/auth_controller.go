package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/auth"
	"github.com/souvik03-136/Go-Store/internal/merrors"
)

// RegisterOAuthUser handles user registration via OAuth.
func RegisterOAuthUser(ctx *gin.Context) {
	username := ctx.PostForm("username")
	if username == "" {
		merrors.BadRequest(ctx, "Username is required")
		return
	}

	token, salt, err := auth.GenerateToken(ctx, username)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to generate token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"salt":  salt,
	})
}

// LoginOAuthUser handles user login via OAuth.
func LoginOAuthUser(ctx *gin.Context) {
	username := ctx.PostForm("username")
	if username == "" {
		merrors.BadRequest(ctx, "Username is required")
		return
	}

	token, salt, err := auth.GenerateToken(ctx, username)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to generate token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"salt":  salt,
	})
}

// ValidateTokenHandler validates the provided JWT token.
func ValidateTokenHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	salt := ctx.Query("salt")

	if token == "" || salt == "" {
		merrors.BadRequest(ctx, "Token and salt are required")
		return
	}

	claims, err := auth.ValidateToken(ctx, token, salt)
	if err != nil {
		merrors.Unauthorized(ctx, "Invalid token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"claims": claims,
	})
}

// RegisterAnonymousUser handles the registration of an anonymous user.
func RegisterAnonymousUser(ctx *gin.Context) {
	anonymousID := auth.GenerateAnonymousID() // Assuming this function generates a unique ID

	token, salt, err := auth.GenerateToken(ctx, anonymousID)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to generate token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"anonymous_id": anonymousID,
		"token":        token,
		"salt":         salt,
	})
}

// LogoutUser handles user logout and token invalidation (if applicable).
func LogoutUser(ctx *gin.Context) {
	// Invalidate the user's token if necessary (e.g., by adding it to a blacklist)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
	})
}
