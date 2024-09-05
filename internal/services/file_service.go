package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/souvik03-136/Go-Store/internal/merrors"
	"github.com/souvik03-136/Go-Store/internal/models"
)

// FileService handles business logic for file operations.
type FileService struct{}

// NewFileService creates a new instance of FileService.
func NewFileService() *FileService {
	return &FileService{}
}

// CreateFile creates a new file record.
func (s *FileService) CreateFile(ctx *gin.Context, name, path, contentType, ownerID string, size int64) (*models.File, error) {
	if name == "" || path == "" || contentType == "" || ownerID == "" {
		merrors.BadRequest(ctx, "Name, path, content type, and owner ID are required")
		return nil, errors.New("name, path, content type, and owner ID are required")
	}

	// Generate unique file ID
	fileID := uuid.New().String()

	// Create a new file model
	file := models.NewFile(fileID, name, path, contentType, ownerID, size)

	// In a real-world scenario, you would store this in a database.
	return file, nil
}

// GetFileByID retrieves a file by its ID.
func (s *FileService) GetFileByID(ctx *gin.Context, fileID string) (*models.File, error) {
	// This would normally query the database for the file.
	file, err := models.GetFileByID(fileID) // You need to implement this function in models/files.go.
	if err != nil {
		merrors.InternalServer(ctx, "Failed to retrieve file")
		return nil, err
	}

	if file == nil {
		merrors.NotFound(ctx, "File not found")
		return nil, errors.New("file not found")
	}

	return file, nil
}

// UpdateFile updates the file's information.
func (s *FileService) UpdateFile(ctx *gin.Context, file *models.File, name, path, contentType string, size int64) (*models.File, error) {
	if name == "" && path == "" && contentType == "" && size == 0 {
		merrors.BadRequest(ctx, "No updates provided")
		return nil, errors.New("no updates provided")
	}

	// Update file details
	file.UpdateFile(name, path, contentType, size)

	// In a real-world scenario, this would involve saving the updated record to the database.
	return file, nil
}

// DeleteFile removes the file by its ID.
func (s *FileService) DeleteFile(ctx *gin.Context, file *models.File) error {
	if file == nil {
		merrors.NotFound(ctx, "File not found")
		return errors.New("file not found")
	}

	// Perform the deletion
	err := file.DeleteFile()
	if err != nil {
		merrors.InternalServer(ctx, "Failed to delete file")
		return err
	}

	// In a real-world scenario, you'd also delete the file from the database and storage.
	return nil
}

// RenameFile renames the file with a new name.
func (s *FileService) RenameFile(ctx *gin.Context, file *models.File, newName string) (*models.File, error) {
	if newName == "" {
		merrors.BadRequest(ctx, "New file name is required")
		return nil, errors.New("new file name is required")
	}

	err := file.RenameFile(newName)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to rename file")
		return nil, err
	}

	// In a real-world scenario, this would involve saving the updated record to the database.
	return file, nil
}

// ListAllFiles retrieves all the files in the system.
func (s *FileService) ListAllFiles(ctx *gin.Context) ([]*models.File, error) {
	files, err := models.GetAllFiles()
	if err != nil {
		merrors.InternalServer(ctx, "Failed to retrieve files")
		return nil, err
	}

	return files, nil
}

// CheckFileOwnership checks if a user is the owner of a file.
func (s *FileService) CheckFileOwnership(ctx *gin.Context, file *models.File, userID string) bool {
	return file.IsOwner(userID)
}
