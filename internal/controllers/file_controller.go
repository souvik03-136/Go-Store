package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/merrors"
	"github.com/souvik03-136/Go-Store/internal/models"
	"github.com/souvik03-136/Go-Store/internal/repository"
)

type FileController struct {
	fileRepo *repository.FileRepository
}

// NewFileController creates a new instance of FileController.
func NewFileController(fileRepo *repository.FileRepository) *FileController {
	return &FileController{fileRepo: fileRepo}
}

// CreateFile handles the creation of a new file.
func (c *FileController) CreateFile(ctx *gin.Context) {
	var file models.File

	// Parse JSON request body into the file model
	if err := ctx.ShouldBindJSON(&file); err != nil {
		merrors.BadRequest(ctx, "Invalid request payload")
		return
	}

	// Save the file in the repository
	if err := c.fileRepo.CreateFile(&file); err != nil {
		merrors.InternalServer(ctx, "Error creating file")
		return
	}

	ctx.JSON(http.StatusCreated, file)
}

// GetFileByID handles fetching a file by ID.
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

// UpdateFile handles updating an existing file's information.
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
		merrors.InternalServer(ctx, "Error updating file")
		return
	}

	ctx.JSON(http.StatusOK, file)
}

// DeleteFile handles the deletion of a file by ID.
func (c *FileController) DeleteFile(ctx *gin.Context) {
	fileID := ctx.Query("id")

	if fileID == "" {
		merrors.BadRequest(ctx, "File ID is required")
		return
	}

	if err := c.fileRepo.DeleteFile(fileID); err != nil {
		merrors.InternalServer(ctx, "Error deleting file")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
