package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	StorageProvider string // Determines the storage provider (e.g., "s3" or "gcs")
	GoogleCloud     GoogleCloudConfig
	AWS             AWSConfig
}

type GoogleCloudConfig struct {
	ProjectID      string
	BucketName     string
	CredentialsKey string
}

type AWSConfig struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
}

func LoadConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system env variables")
	}

	// Populate Google Cloud config
	googleCloudConfig := GoogleCloudConfig{
		ProjectID:      os.Getenv("GOOGLE_CLOUD_PROJECT_ID"),
		BucketName:     os.Getenv("GOOGLE_CLOUD_BUCKET_NAME"),
		CredentialsKey: os.Getenv("GOOGLE_CLOUD_CREDENTIALS_KEY"),
	}

	// Populate AWS config
	awsConfig := AWSConfig{
		Region:          os.Getenv("AWS_REGION"),
		BucketName:      os.Getenv("AWS_BUCKET_NAME"),
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	// Read the storage provider (e.g., "s3" or "gcs")
	storageProvider := os.Getenv("STORAGE_PROVIDER")

	// Combine into main config
	config := &Config{
		StorageProvider: storageProvider,
		GoogleCloud:     googleCloudConfig,
		AWS:             awsConfig,
	}

	return config, nil
}
