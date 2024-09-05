package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/souvik03-136/Go-Store/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user record into the database.
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user from the database by their ID.
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates a user's information in the database.
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users SET username = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, time.Now(), user.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser removes a user from the database by their ID.
func (r *UserRepository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
