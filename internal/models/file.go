package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// File represents a file stored in the system.
type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Url         string    `json:"url"` // URL to access the file
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	OwnerID     string    `json:"owner_id"` // References the user who uploaded the file
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewFile creates a new File instance.
func NewFile(id, name, path, url, contentType, ownerID string, size int64) *File {
	return &File{
		ID:          id,
		Name:        name,
		Path:        path,
		Url:         url,
		Size:        size,
		ContentType: contentType,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

var db *sql.DB // Assuming this will be initialized externally

// UpdateFile updates the file information.
func (f *File) UpdateFile(name, path, url, contentType string, size int64) {
	f.Name = name
	f.Path = path
	f.Url = url
	f.ContentType = contentType
	f.Size = size
	f.UpdatedAt = time.Now()
}

// IsOwner checks if the provided userID matches the file's owner ID.
func (f *File) IsOwner(userID string) bool {
	return f.OwnerID == userID
}

// RenameFile allows renaming the file.
func (f *File) RenameFile(newName string) error {
	if newName == "" {
		return errors.New("file name cannot be empty")
	}
	f.Name = newName
	f.UpdatedAt = time.Now()
	return nil
}

// DeleteFile clears the file metadata, mimicking deletion.
// In practice, this would involve deleting the actual file and its metadata from the database and storage.
func (f *File) DeleteFile() error {
	f.Name = ""
	f.Path = ""
	f.Url = ""
	f.Size = 0
	f.ContentType = ""
	f.UpdatedAt = time.Now()
	return nil
}

// GetAllFiles returns a list of all files in the system.
func GetAllFiles() ([]*File, error) {
	query := `SELECT id, name, path, url, size, content_type, owner_id, created_at, updated_at FROM files`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %v", err)
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		var file File
		if err := rows.Scan(&file.ID, &file.Name, &file.Path, &file.Url, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file: %v", err)
		}
		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return files, nil
}

// GetFileByID retrieves a file by its ID.
func GetFileByID(id string) (*File, error) {
	var file File

	query := `SELECT id, name, path, url, size, content_type, owner_id, created_at, updated_at FROM files WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&file.ID, &file.Name, &file.Path, &file.Url, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt, &file.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("file not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to query file: %v", err)
	}

	return &file, nil
}

// CreateFile saves a new file's metadata in the database.
func CreateFile(file *File) error {
	query := `INSERT INTO files (id, name, path, url, size, content_type, owner_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := db.Exec(query, file.ID, file.Name, file.Path, file.Url, file.Size, file.ContentType, file.OwnerID, file.CreatedAt, file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert file: %v", err)
	}
	return nil
}

// UpdateFileInDB updates a file's metadata in the database.
func UpdateFileInDB(file *File) error {
	query := `UPDATE files SET name = $1, path = $2, url = $3, size = $4, content_type = $5, updated_at = $6 WHERE id = $7`
	_, err := db.Exec(query, file.Name, file.Path, file.Url, file.Size, file.ContentType, file.UpdatedAt, file.ID)
	if err != nil {
		return fmt.Errorf("failed to update file: %v", err)
	}
	return nil
}

// DeleteFileByID deletes a file by ID from the database.
func DeleteFileByID(id string) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}
