package repository

import (
	"database/sql"
	"errors"
)

type Permission struct {
	ID     string
	UserID string
	FileID string
	Level  string // e.g., "read", "write", "admin"
}

type PermissionRepository struct {
	db *sql.DB
}

// NewPermissionRepository creates a new instance of PermissionRepository.
func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// GrantPermission grants a user a specific level of permission for a file.
func (r *PermissionRepository) GrantPermission(permission *Permission) error {
	query := `
		INSERT INTO permissions (id, user_id, file_id, level)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, permission.ID, permission.UserID, permission.FileID, permission.Level)
	if err != nil {
		return err
	}
	return nil
}

// GetPermission checks if a user has permission to access a file.
func (r *PermissionRepository) GetPermission(userID, fileID string) (*Permission, error) {
	var permission Permission

	query := `
		SELECT id, user_id, file_id, level
		FROM permissions WHERE user_id = $1 AND file_id = $2
	`
	err := r.db.QueryRow(query, userID, fileID).Scan(&permission.ID, &permission.UserID, &permission.FileID, &permission.Level)
	if err == sql.ErrNoRows {
		return nil, errors.New("no permission found")
	} else if err != nil {
		return nil, err
	}

	return &permission, nil
}

// RevokePermission revokes a user's permission for a file.
func (r *PermissionRepository) RevokePermission(userID, fileID string) error {
	query := `DELETE FROM permissions WHERE user_id = $1 AND file_id = $2`
	_, err := r.db.Exec(query, userID, fileID)
	if err != nil {
		return err
	}
	return nil
}
