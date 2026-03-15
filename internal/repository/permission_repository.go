// internal/repository/permission_repository.go

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/souvik03-136/Go-Store/internal/models"
)

// PermissionRepository handles all database operations for file permissions.
type PermissionRepository struct {
	db *sql.DB
}

// NewPermissionRepository creates a new PermissionRepository.
func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// CreatePermission inserts a new permission record.
func (r *PermissionRepository) CreatePermission(p *models.Permission) error {
	query := `
		INSERT INTO permissions (id, file_id, user_id, can_read, can_write, can_delete, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query,
		p.ID, p.FileID, p.UserID,
		p.CanRead, p.CanWrite, p.CanDelete,
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("creating permission: %w", err)
	}
	return nil
}

// GetPermission retrieves the permission record for a specific user/file pair.
func (r *PermissionRepository) GetPermission(userID, fileID string) (*models.Permission, error) {
	query := `
		SELECT id, file_id, user_id, can_read, can_write, can_delete, created_at, updated_at
		FROM permissions WHERE user_id = $1 AND file_id = $2
	`
	var p models.Permission
	err := r.db.QueryRow(query, userID, fileID).Scan(
		&p.ID, &p.FileID, &p.UserID,
		&p.CanRead, &p.CanWrite, &p.CanDelete,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("permission not found")
	}
	if err != nil {
		return nil, fmt.Errorf("getting permission: %w", err)
	}
	return &p, nil
}

// ListPermissionsByFile retrieves all permissions for a given file.
func (r *PermissionRepository) ListPermissionsByFile(fileID string) ([]*models.Permission, error) {
	query := `
		SELECT id, file_id, user_id, can_read, can_write, can_delete, created_at, updated_at
		FROM permissions WHERE file_id = $1
	`
	rows, err := r.db.Query(query, fileID)
	if err != nil {
		return nil, fmt.Errorf("listing permissions for file %s: %w", fileID, err)
	}
	defer rows.Close()

	var perms []*models.Permission
	for rows.Next() {
		p := &models.Permission{}
		if err := rows.Scan(
			&p.ID, &p.FileID, &p.UserID,
			&p.CanRead, &p.CanWrite, &p.CanDelete,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning permission row: %w", err)
		}
		perms = append(perms, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating permission rows: %w", err)
	}
	return perms, nil
}

// UpdatePermission updates the access flags for an existing user/file permission.
func (r *PermissionRepository) UpdatePermission(p *models.Permission) error {
	query := `
		UPDATE permissions
		SET can_read = $1, can_write = $2, can_delete = $3, updated_at = $4
		WHERE file_id = $5 AND user_id = $6
	`
	result, err := r.db.Exec(query,
		p.CanRead, p.CanWrite, p.CanDelete, time.Now(),
		p.FileID, p.UserID,
	)
	if err != nil {
		return fmt.Errorf("updating permission: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

// DeletePermission removes the permission for a specific user/file pair.
func (r *PermissionRepository) DeletePermission(userID, fileID string) error {
	result, err := r.db.Exec(
		`DELETE FROM permissions WHERE user_id = $1 AND file_id = $2`,
		userID, fileID,
	)
	if err != nil {
		return fmt.Errorf("deleting permission: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("permission not found")
	}
	return nil
}
