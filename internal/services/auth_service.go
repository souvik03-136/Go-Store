// internal/services/auth_service.go

package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/souvik03-136/Go-Store/internal/auth"
	"github.com/souvik03-136/Go-Store/internal/models"
	"github.com/souvik03-136/Go-Store/internal/repository"
)

const tokenExpiry = 24 * time.Hour

// AuthService handles all authentication business logic.
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// TokenPair holds a signed JWT and the salt required to validate it.
type TokenPair struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}

// RegisterOAuth creates a new user record from the supplied credentials and
// returns a TokenPair. It rejects duplicate email addresses.
func (s *AuthService) RegisterOAuth(username, email, password string) (*models.User, *TokenPair, error) {
	// Check for duplicate email
	if existing, _ := s.userRepo.GetUserByEmail(email); existing != nil {
		return nil, nil, errors.New("a user with that email already exists")
	}

	user, err := models.NewUser(generateID(), username, email, password)
	if err != nil {
		return nil, nil, fmt.Errorf("building user: %w", err)
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, nil, fmt.Errorf("persisting user: %w", err)
	}

	pair, err := s.issueToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

// LoginOAuth authenticates a user by email + password and issues a TokenPair.
func (s *AuthService) LoginOAuth(email, password string) (*models.User, *TokenPair, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return nil, nil, errors.New("invalid credentials")
	}

	pair, err := s.issueToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

// RegisterAnonymous issues a TokenPair for an anonymous session without
// creating a persistent user record.
func (s *AuthService) RegisterAnonymous() (string, *TokenPair, error) {
	anonID := auth.GenerateAnonymousID()
	pair, err := s.issueToken(anonID)
	if err != nil {
		return "", nil, err
	}
	return anonID, pair, nil
}

// ValidateToken parses and validates a JWT, returning the subject claim.
func (s *AuthService) ValidateToken(tokenString, salt string) (string, error) {
	claims, err := auth.ValidateToken(s.jwtSecret, tokenString, salt)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	return claims.Subject, nil
}

func (s *AuthService) issueToken(subject string) (*TokenPair, error) {
	tokenString, salt, err := auth.GenerateToken(s.jwtSecret, subject, tokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}
	return &TokenPair{Token: tokenString, Salt: salt}, nil
}
