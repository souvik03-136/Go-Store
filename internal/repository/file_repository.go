// internal/repository/file_repository.go

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/souvik03-136/Go-Store/internal/models"
)

// FileRepository handles all database operations for file metadata.
type FileRepository struct {
	db *sql.DB
}

// NewFileRepository creates a new FileRepository backed by the given *sql.DB.
func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

// CreateFile inserts a new file metadata record.
func (r *FileRepository) CreateFile(file *models.File) error {
	query := `
		INSERT INTO files (id, name, path, url, size, content_type, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query,
		file.ID, file.Name, file.Path, file.Url,
		file.Size, file.ContentType, file.OwnerID,
		file.CreatedAt, file.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	return nil
}

// GetFileByID retrieves file metadata by UUID.
func (r *FileRepository) GetFileByID(id string) (*models.File, error) {
	query := `
		SELECT id, name, path, url, size, content_type, owner_id, created_at, updated_at
		FROM files WHERE id = $1
	`
	return r.scanFile(r.db.QueryRow(query, id))
}

// ListFilesByOwner retrieves all files owned by the specified user.
func (r *FileRepository) ListFilesByOwner(ownerID string) ([]*models.File, error) {
	query := `
		SELECT id, name, path, url, size, content_type, owner_id, created_at, updated_at
		FROM files WHERE owner_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("listing files for owner %s: %w", ownerID, err)
	}
	defer rows.Close()

	var files []*models.File
	for rows.Next() {
		f := &models.File{}
		if err := rows.Scan(
			&f.ID, &f.Name, &f.Path, &f.Url,
			&f.Size, &f.ContentType, &f.OwnerID,
			&f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning file row: %w", err)
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating file rows: %w", err)
	}
	return files, nil
}

// UpdateFile updates mutable file metadata fields.
func (r *FileRepository) UpdateFile(file *models.File) error {
	query := `
		UPDATE files
		SET name = $1, path = $2, url = $3, size = $4, content_type = $5, updated_at = $6
		WHERE id = $7
	`
	result, err := r.db.Exec(query,
		file.Name, file.Path, file.Url,
		file.Size, file.ContentType,
		time.Now(), file.ID,
	)
	if err != nil {
		return fmt.Errorf("updating file %s: %w", file.ID, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("file not found")
	}
	return nil
}

// DeleteFile removes a file metadata record by UUID.
func (r *FileRepository) DeleteFile(id string) error {
	result, err := r.db.Exec(`DELETE FROM files WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting file %s: %w", id, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("file not found")
	}
	return nil
}

func (r *FileRepository) scanFile(row *sql.Row) (*models.File, error) {
	f := &models.File{}
	err := row.Scan(
		&f.ID, &f.Name, &f.Path, &f.Url,
		&f.Size, &f.ContentType, &f.OwnerID,
		&f.CreatedAt, &f.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("file not found")
	}
	if err != nil {
		return nil, fmt.Errorf("scanning file: %w", err)
	}
	return f, nil
}
