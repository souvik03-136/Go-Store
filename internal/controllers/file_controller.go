package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/config"
	"github.com/souvik03-136/Go-Store/internal/merrors"
	"github.com/souvik03-136/Go-Store/internal/models"
	"github.com/souvik03-136/Go-Store/internal/repository"
	"github.com/souvik03-136/Go-Store/internal/storage"
)

type FileController struct {
	fileRepo *repository.FileRepository
	storage  storage.Storage // This will be either S3 or Google Cloud Storage
}

// NewFileController creates a new FileController with the specified repository and configuration
func NewFileController(fileRepo *repository.FileRepository, cfg *config.Config) (*FileController, error) {
	var store storage.Storage
	var err error

	if cfg.StorageProvider == "s3" {
		// Initialize AWS S3 Storage
		store, err = storage.NewS3Storage(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, cfg.AWS.Region, cfg.AWS.BucketName)
	} else if cfg.StorageProvider == "gcs" {
		// Initialize Google Cloud Storage
		ctx := context.Background()
		store, err = storage.NewGC3Storage(ctx, cfg.GoogleCloud.CredentialsKey, cfg.GoogleCloud.BucketName)
	} else {
		// Return a custom error message
		return nil, errors.New("unsupported storage provider")
	}

	if err != nil {
		return nil, err
	}

	return &FileController{fileRepo: fileRepo, storage: store}, nil
}

// CreateFile handles the creation of a new file, uploads to storage, and saves metadata in the repository.
func (c *FileController) CreateFile(ctx *gin.Context) {
	// Get file from form-data
	file, err := ctx.FormFile("file")
	if err != nil {
		merrors.BadRequest(ctx, "File upload failed")
		return
	}

	// Upload file to cloud storage
	fileURL, err := c.storage.UploadFile(ctx, file, file.Filename)
	if err != nil {
		merrors.InternalServer(ctx, "Error uploading file to storage")
		return
	}

	// Create file metadata
	var fileModel models.File
	if err := ctx.ShouldBindJSON(&fileModel); err != nil {
		merrors.BadRequest(ctx, "Invalid request payload")
		return
	}

	fileModel.Url = fileURL // Save the URL returned by cloud storage

	// Save the file metadata in the repository
	if err := c.fileRepo.CreateFile(&fileModel); err != nil {
		merrors.InternalServer(ctx, "Error saving file metadata")
		return
	}

	ctx.JSON(http.StatusCreated, fileModel)
}

// GetFileByID handles fetching a file's metadata by ID from the repository.
func (c *FileController) GetFileByID(ctx *gin.Context) {
	fileID := ctx.Query("id")

	if fileID == "" {
		merrors.BadRequest(ctx, "File ID is required")
		return
	}

	file, err := c.fileRepo.GetFileByID(fileID)
	if err != nil {
		merrors.NotFound(ctx, "File not found")
		return
	}

	ctx.JSON(http.StatusOK, file)
}

// UpdateFile handles updating an existing file's metadata in the repository.
func (c *FileController) UpdateFile(ctx *gin.Context) {
	fileID := ctx.Query("id")

	if fileID == "" {
		merrors.BadRequest(ctx, "File ID is required")
		return
	}

	var file models.File
	if err := ctx.ShouldBindJSON(&file); err != nil {
		merrors.BadRequest(ctx, "Invalid request payload")
		return
	}

	file.ID = fileID
	if err := c.fileRepo.UpdateFile(&file); err != nil {
		merrors.InternalServer(ctx, "Error updating file metadata")
		return
	}

	ctx.JSON(http.StatusOK, file)
}

// DeleteFile handles the deletion of a file by ID, deletes the file from cloud storage, and removes metadata from the repository.
func (c *FileController) DeleteFile(ctx *gin.Context) {
	fileID := ctx.Query("id")

	if fileID == "" {
		merrors.BadRequest(ctx, "File ID is required")
		return
	}

	// Fetch the file metadata from the repository
	file, err := c.fileRepo.GetFileByID(fileID)
	if err != nil {
		merrors.NotFound(ctx, "File not found")
		return
	}

	// Delete the file from cloud storage
	err = c.storage.DeleteFile(ctx, file.Url)
	if err != nil {
		merrors.InternalServer(ctx, "Error deleting file from storage")
		return
	}

	// Delete the file metadata from the repository
	if err := c.fileRepo.DeleteFile(fileID); err != nil {
		merrors.InternalServer(ctx, "Error deleting file metadata")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
