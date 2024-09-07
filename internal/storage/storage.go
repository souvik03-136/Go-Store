package storage

import (
	"context"
	"mime/multipart"
)

// Storage interface that both S3 and GC3 will implement
type Storage interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, destination string) (string, error)
	DeleteFile(ctx context.Context, filePath string) error
}
