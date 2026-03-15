// internal/controllers/file_controller.go

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/auth"
	"github.com/souvik03-136/Go-Store/internal/merrors"
	"github.com/souvik03-136/Go-Store/internal/services"
)

// FileController handles HTTP requests for file management endpoints.
type FileController struct {
	fileSvc *services.FileService
}

// NewFileController creates a new FileController.
func NewFileController(fileSvc *services.FileService) *FileController {
	return &FileController{fileSvc: fileSvc}
}

type updateFileRequest struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
}

// requesterID extracts the authenticated subject from the Gin context.
// It returns an empty string if the context has no authenticated user
// (e.g. on public endpoints), which downstream ownership checks will reject.
func requesterID(ctx *gin.Context) string {
	if v, ok := ctx.Get(auth.SubjectKey); ok {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// CreateFile handles POST /v1/files
// Expects a multipart/form-data request with the file under the "file" field.
func (c *FileController) CreateFile(ctx *gin.Context) {
	ownerID := requesterID(ctx)
	if ownerID == "" {
		merrors.Unauthorized(ctx, "authentication required")
		return
	}

	header, err := ctx.FormFile("file")
	if err != nil {
		merrors.BadRequest(ctx, "a file field is required in the multipart form")
		return
	}

	file, err := c.fileSvc.UploadFile(ctx.Request.Context(), header, ownerID)
	if err != nil {
		merrors.InternalServer(ctx, "file upload failed")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": file})
}

// GetFileByID handles GET /v1/files/:id
func (c *FileController) GetFileByID(ctx *gin.Context) {
	id := ctx.Param("id")

	file, err := c.fileSvc.GetFileByID(id)
	if err != nil {
		merrors.NotFound(ctx, "file not found")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": file})
}

// ListMyFiles handles GET /v1/files  (authenticated user's own files)
func (c *FileController) ListMyFiles(ctx *gin.Context) {
	ownerID := requesterID(ctx)
	if ownerID == "" {
		merrors.Unauthorized(ctx, "authentication required")
		return
	}

	files, err := c.fileSvc.ListFilesByOwner(ownerID)
	if err != nil {
		merrors.InternalServer(ctx, "failed to list files")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": files})
}

// UpdateFile handles PUT /v1/files/:id
func (c *FileController) UpdateFile(ctx *gin.Context) {
	id := ctx.Param("id")
	ownerID := requesterID(ctx)

	var req updateFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		merrors.Validation(ctx, err.Error())
		return
	}

	file, err := c.fileSvc.UpdateFile(id, ownerID, req.Name, req.ContentType, 0)
	if err != nil {
		switch err.Error() {
		case "permission denied":
			merrors.Forbidden(ctx, "you do not own this file")
		case "file not found":
			merrors.NotFound(ctx, "file not found")
		default:
			merrors.InternalServer(ctx, "failed to update file")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": file})
}

// DeleteFile handles DELETE /v1/files/:id
func (c *FileController) DeleteFile(ctx *gin.Context) {
	id := ctx.Param("id")
	ownerID := requesterID(ctx)

	if err := c.fileSvc.DeleteFile(ctx.Request.Context(), id, ownerID); err != nil {
		switch err.Error() {
		case "permission denied":
			merrors.Forbidden(ctx, "you do not own this file")
		case "file not found":
			merrors.NotFound(ctx, "file not found")
		default:
			merrors.InternalServer(ctx, "failed to delete file")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}
