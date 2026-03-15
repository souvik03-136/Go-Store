// internal/repository/user_repository.go

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/souvik03-136/Go-Store/internal/models"
)

// UserRepository handles all database operations for users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository backed by the given *sql.DB.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user record. The PasswordHash field must already be
// populated (bcrypt hashing is the caller's responsibility).
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query,
		user.ID, user.Username, user.Email,
		user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

// GetUserByID fetches a single user by their UUID.
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users WHERE id = $1
	`
	return r.scanUser(r.db.QueryRow(query, id))
}

// GetUserByEmail fetches a single user by email address.
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users WHERE email = $1
	`
	return r.scanUser(r.db.QueryRow(query, email))
}

// GetUserByUsername fetches a single user by their username.
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users WHERE username = $1
	`
	return r.scanUser(r.db.QueryRow(query, username))
}

// UpdateUser updates a user's mutable fields and stamps updated_at.
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, updated_at = $4
		WHERE id = $5
	`
	result, err := r.db.Exec(query,
		user.Username, user.Email, user.PasswordHash, time.Now(), user.ID,
	)
	if err != nil {
		return fmt.Errorf("updating user %s: %w", user.ID, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

// DeleteUser removes a user by their UUID.
func (r *UserRepository) DeleteUser(id string) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting user %s: %w", id, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) scanUser(row *sql.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email,
		&u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("scanning user: %w", err)
	}
	return &u, nil
}
