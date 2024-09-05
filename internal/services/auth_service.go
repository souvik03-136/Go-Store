package services

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/auth"
	"github.com/souvik03-136/Go-Store/internal/merrors"
)

// AuthService handles the business logic for authentication.
type AuthService struct{}

// NewAuthService creates a new instance of AuthService.
func NewAuthService() *AuthService {
	return &AuthService{}
}

// GenerateToken generates a JWT token for a given username.
func (s *AuthService) GenerateToken(ctx *gin.Context, username string) (string, string, error) {
	salt, err := auth.GenerateDynamicSalt(ctx)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to generate dynamic salt")
		return "", "", err
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		Subject:   username,
	}

	signingSecret, err := auth.GetSigningSecret(ctx, salt)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to get signing secret")
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingSecret)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to sign the JWT token")
		return "", "", err
	}

	return tokenString, salt, nil
}

// ValidateToken validates a given JWT token using the provided salt.
func (s *AuthService) ValidateToken(ctx *gin.Context, tokenString string, salt string) (*jwt.StandardClaims, error) {
	signingSecret, err := auth.GetSigningSecret(ctx, salt)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to get signing secret")
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			merrors.Unauthorized(ctx, "Invalid JWT signature")
			return nil, err
		}
		merrors.Unauthorized(ctx, "Invalid JWT token")
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	merrors.Unauthorized(ctx, "Invalid JWT token")
	return nil, err
}

// HandleOAuthLogin processes OAuth login and returns the JWT token and any associated error.
func (s *AuthService) HandleOAuthLogin(ctx *gin.Context, oauthID string) (string, string, error) {
	// Perform OAuth login logic here (e.g., verify OAuth token, retrieve user details)
	// For now, we assume the OAuth login is successful and return a token.

	username := "user_from_oauth_" + oauthID // This should be the username or unique identifier from the OAuth provider.
	return s.GenerateToken(ctx, username)
}

// HandleAnonymousLogin generates an anonymous ID and returns a JWT token.
func (s *AuthService) HandleAnonymousLogin(ctx *gin.Context) (string, string, error) {
	anonymousID := auth.GenerateAnonymousID()
	return s.GenerateToken(ctx, anonymousID)
}

// CheckJWTSecret verifies the presence of the JWT secret key in environment variables.
func (s *AuthService) CheckJWTSecret(ctx *gin.Context) error {
	if os.Getenv("JWT_SECRET_KEY") == "" {
		merrors.InternalServer(ctx, "JWT secret key not set in environment variables")
		return nil
	}
	return nil
}
