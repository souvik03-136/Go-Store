// internal/services/file_service.go

package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/souvik03-136/Go-Store/internal/models"
	"github.com/souvik03-136/Go-Store/internal/repository"
	"github.com/souvik03-136/Go-Store/internal/storage"
)

// FileService handles all file-management business logic.
type FileService struct {
	fileRepo *repository.FileRepository
	storage  storage.Storage
}

// NewFileService creates a new FileService.
func NewFileService(fileRepo *repository.FileRepository, store storage.Storage) *FileService {
	return &FileService{fileRepo: fileRepo, storage: store}
}

// UploadFile uploads the multipart file to cloud storage, persists its
// metadata, and returns the created File record.
func (s *FileService) UploadFile(ctx context.Context, header *multipart.FileHeader, ownerID string) (*models.File, error) {
	if header == nil {
		return nil, errors.New("file is required")
	}
	if ownerID == "" {
		return nil, errors.New("owner id is required")
	}

	fileID := generateID()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	// Store as <ownerID>/<fileID><ext> so each owner's files are namespaced.
	objectKey := fmt.Sprintf("%s/%s%s", ownerID, fileID, ext)

	url, err := s.storage.UploadFile(ctx, header, objectKey)
	if err != nil {
		return nil, fmt.Errorf("uploading file to storage: %w", err)
	}

	file := models.NewFile(
		fileID,
		header.Filename,
		objectKey,
		url,
		header.Header.Get("Content-Type"),
		ownerID,
		header.Size,
	)

	if err := s.fileRepo.CreateFile(file); err != nil {
		// Attempt to clean up orphaned cloud object on DB failure.
		_ = s.storage.DeleteFile(ctx, objectKey)
		return nil, fmt.Errorf("persisting file metadata: %w", err)
	}

	return file, nil
}

// GetFileByID retrieves file metadata, enforcing that the requester owns
// the file or has been granted read access (permission check is optional
// and can be layered on top by the controller).
func (s *FileService) GetFileByID(id string) (*models.File, error) {
	if id == "" {
		return nil, errors.New("file id is required")
	}
	return s.fileRepo.GetFileByID(id)
}

// ListFilesByOwner returns all files belonging to ownerID.
func (s *FileService) ListFilesByOwner(ownerID string) ([]*models.File, error) {
	if ownerID == "" {
		return nil, errors.New("owner id is required")
	}
	return s.fileRepo.ListFilesByOwner(ownerID)
}

// UpdateFile applies non-zero field updates to the file identified by id.
// Only the owner may update their own files.
func (s *FileService) UpdateFile(id, requesterID, name, contentType string, size int64) (*models.File, error) {
	if id == "" {
		return nil, errors.New("file id is required")
	}

	file, err := s.fileRepo.GetFileByID(id)
	if err != nil {
		return nil, err
	}

	if !file.IsOwner(requesterID) {
		return nil, errors.New("permission denied")
	}

	if name != "" {
		if err := file.Rename(name); err != nil {
			return nil, err
		}
	}
	if contentType != "" {
		file.ContentType = contentType
	}
	if size > 0 {
		file.Size = size
	}

	if err := s.fileRepo.UpdateFile(file); err != nil {
		return nil, fmt.Errorf("persisting file update: %w", err)
	}

	return file, nil
}

// DeleteFile removes the cloud object and the metadata record.
// Only the file owner may delete their own files.
func (s *FileService) DeleteFile(ctx context.Context, id, requesterID string) error {
	if id == "" {
		return errors.New("file id is required")
	}

	file, err := s.fileRepo.GetFileByID(id)
	if err != nil {
		return err
	}

	if !file.IsOwner(requesterID) {
		return errors.New("permission denied")
	}

	// Delete from cloud storage first; if that fails do not remove the metadata
	// so the admin can retry the cloud deletion.
	if err := s.storage.DeleteFile(ctx, file.Path); err != nil {
		return fmt.Errorf("deleting object from storage: %w", err)
	}

	if err := s.fileRepo.DeleteFile(id); err != nil {
		return fmt.Errorf("deleting file metadata: %w", err)
	}

	return nil
}
