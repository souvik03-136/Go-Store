package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GC3Storage struct to hold the Google Cloud Storage client and bucket details
type GC3Storage struct {
	client *storage.Client
	bucket string
}

// NewGC3Storage initializes a new Google Cloud Storage client
func NewGC3Storage(ctx context.Context, credentialsFile, bucket string) (*GC3Storage, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %v", err)
	}
	return &GC3Storage{client: client, bucket: bucket}, nil
}

// UploadFile uploads a file to Google Cloud Storage
func (g *GC3Storage) UploadFile(ctx context.Context, file *multipart.FileHeader, destination string) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	wc := g.client.Bucket(g.bucket).Object(destination).NewWriter(ctx)
	if _, err = wc.Write([]byte{}); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucket, destination)
	return fileURL, nil
}

// DeleteFile deletes a file from Google Cloud Storage
func (g *GC3Storage) DeleteFile(ctx context.Context, filePath string) error {
	if err := g.client.Bucket(g.bucket).Object(filePath).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file from GCS: %v", err)
	}
	return nil
}
