package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/souvik03-136/Go-Store/internal/models"
)

type FileRepository struct {
	db *sql.DB
}

// NewFileRepository creates a new instance of FileRepository.
func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

// CreateFile inserts a new file record into the database.
func (r *FileRepository) CreateFile(file *models.File) error {
	query := `
		INSERT INTO files (id, name, path, size, content_type, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, file.ID, file.Name, file.Path, file.Size, file.ContentType, file.OwnerID, file.CreatedAt, file.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetFileByID retrieves a file from the database by its ID.
func (r *FileRepository) GetFileByID(id string) (*models.File, error) {
	var file models.File

	query := `
		SELECT id, name, path, size, content_type, owner_id, created_at, updated_at
		FROM files WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(&file.ID, &file.Name, &file.Path, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt, &file.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("file not found")
	} else if err != nil {
		return nil, err
	}

	return &file, nil
}

// UpdateFile updates a file's information in the database.
func (r *FileRepository) UpdateFile(file *models.File) error {
	query := `
		UPDATE files SET name = $1, path = $2, size = $3, content_type = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(query, file.Name, file.Path, file.Size, file.ContentType, time.Now(), file.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile removes a file from the database by its ID.
func (r *FileRepository) DeleteFile(id string) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
