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

// S3Storage struct to hold the S3 client and bucket details
type S3Storage struct {
	client *s3.S3
	bucket string
}

// NewS3Storage initializes a new S3 client
func NewS3Storage(accessKey, secretKey, region, bucket string) (*S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)
	return &S3Storage{client: svc, bucket: bucket}, nil
}

// UploadFile uploads a file to S3
func (s *S3Storage) UploadFile(ctx context.Context, file *multipart.FileHeader, destination string) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(destination),
		Body:   f,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, destination)
	return fileURL, nil
}

// DeleteFile deletes a file from S3
func (s *S3Storage) DeleteFile(ctx context.Context, filePath string) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}
	return nil
}
