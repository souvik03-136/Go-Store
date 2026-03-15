// internal/models/file.go

package models

import (
	"errors"
	"time"
)

// File represents the metadata of a stored file.
type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Url         string    `json:"url"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewFile constructs a File with the current timestamp.
func NewFile(id, name, path, url, contentType, ownerID string, size int64) *File {
	now := time.Now()
	return &File{
		ID:          id,
		Name:        name,
		Path:        path,
		Url:         url,
		Size:        size,
		ContentType: contentType,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Rename changes the file name and stamps UpdatedAt.
func (f *File) Rename(newName string) error {
	if newName == "" {
		return errors.New("file name cannot be empty")
	}
	f.Name = newName
	f.UpdatedAt = time.Now()
	return nil
}

// IsOwner returns true when the given userID matches the file's owner.
func (f *File) IsOwner(userID string) bool {
	return f.OwnerID == userID
}
