// internal/storage/s3_storage.go

package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Storage implements Storage using Amazon S3.
type S3Storage struct {
	client *s3.S3
	bucket string
	region string
}

// NewS3Storage initialises a new S3 client and returns an S3Storage.
func NewS3Storage(accessKeyID, secretAccessKey, region, bucket string) (*S3Storage, error) {
	if accessKeyID == "" || secretAccessKey == "" || region == "" || bucket == "" {
		return nil, fmt.Errorf("accessKeyID, secretAccessKey, region, and bucket are all required")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID, secretAccessKey, "",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("creating AWS session: %w", err)
	}

	return &S3Storage{
		client: s3.New(sess),
		bucket: bucket,
		region: region,
	}, nil
}

// UploadFile uploads the given multipart file to S3 under the specified
// destination key and returns the public HTTPS URL of the object.
func (s *S3Storage) UploadFile(ctx context.Context, file *multipart.FileHeader, destination string) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("opening uploaded file: %w", err)
	}
	defer f.Close()

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(destination),
		Body:        f,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	}

	if _, err := s.client.PutObjectWithContext(ctx, input); err != nil {
		return "", fmt.Errorf("uploading to S3 (bucket=%s key=%s): %w", s.bucket, destination, err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, destination)
	return url, nil
}

// DeleteFile removes the object identified by objectKey from S3.
func (s *S3Storage) DeleteFile(ctx context.Context, objectKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	}
	if _, err := s.client.DeleteObjectWithContext(ctx, input); err != nil {
		return fmt.Errorf("deleting from S3 (bucket=%s key=%s): %w", s.bucket, objectKey, err)
	}
	return nil
}
