// internal/storage/storage.go

package storage

import (
	"context"
	"mime/multipart"
)

// Storage is the common interface implemented by every cloud storage backend.
// UploadFile uploads the multipart file to the given destination path and
// returns the public URL of the stored object.
// DeleteFile removes the object at the given object key / URL.
type Storage interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, destination string) (string, error)
	DeleteFile(ctx context.Context, objectKey string) error
}
