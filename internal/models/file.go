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
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	OwnerID     string    `json:"owner_id"` // References the user who uploaded the file
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewFile creates a new File instance.
func NewFile(id, name, path, contentType, ownerID string, size int64) *File {
	return &File{
		ID:          id,
		Name:        name,
		Path:        path,
		Size:        size,
		ContentType: contentType,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

var db *sql.DB

// UpdateFile updates the file information.
func (f *File) UpdateFile(name, path, contentType string, size int64) {
	f.Name = name
	f.Path = path
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

// DeleteFile deletes the file by setting its attributes to empty or nil.
// In practice, you would probably delete the record from your database and the actual file from storage.
func (f *File) DeleteFile() error {
	f.Name = ""
	f.Path = ""
	f.Size = 0
	f.ContentType = ""
	f.UpdatedAt = time.Now()
	return nil
}

// GetAllFiles returns a list of all files in the system.
// In practice, this would involve querying your database to retrieve the list of files.
func GetAllFiles() ([]*File, error) {
	// This would typically be a database query, returning a list of files.
	// For demonstration, we'll return a dummy list of files.
	files := []*File{
		{
			ID:          "file1",
			Name:        "example1.txt",
			Path:        "/files/example1.txt",
			Size:        1024,
			ContentType: "text/plain",
			OwnerID:     "user1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "file2",
			Name:        "example2.jpg",
			Path:        "/files/example2.jpg",
			Size:        2048,
			ContentType: "image/jpeg",
			OwnerID:     "user2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	return files, nil
}

// GetFileByID retrieves a file by its ID.
func GetFileByID(id string) (*File, error) {
	var file File

	query := `SELECT id, name, path, size, content_type, owner_id, created_at, updated_at FROM files WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&file.ID, &file.Name, &file.Path, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt, &file.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("file not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to query file: %v", err)
	}

	return &file, nil
}
