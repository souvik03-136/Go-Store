package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	JWT             JWTConfig
	StorageProvider string
	GoogleCloud     GoogleCloudConfig
	AWS             AWSConfig
}

type ServerConfig struct {
	Port string
	Env  string // "development" | "production"
}

type DatabaseConfig struct {
	Driver   string // "postgres" | "mysql"
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// DSN returns the data source name string for the configured driver.
func (d DatabaseConfig) DSN() string {
	switch d.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
			d.User, d.Password, d.Host, d.Port, d.Name)
	default: // postgres
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
	}
}

type JWTConfig struct {
	SecretKey string
}

type GoogleCloudConfig struct {
	ProjectID       string
	BucketName      string
	CredentialsFile string
}

type AWSConfig struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
}

// LoadConfig reads the .env file (if present) and then environment variables,
// returning a validated Config or an error.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading configuration from environment")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
			Env:  getEnvOrDefault("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Driver:   getEnvOrDefault("DB_DRIVER", "postgres"),
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			Name:     getEnvOrDefault("DB_DATABASE", "gostore"),
			User:     getEnvOrDefault("DB_USERNAME", "postgres"),
			Password: os.Getenv("DB_PASSWORD"),
			SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
		},
		StorageProvider: getEnvOrDefault("STORAGE_PROVIDER", "s3"),
		GoogleCloud: GoogleCloudConfig{
			ProjectID:       os.Getenv("GOOGLE_CLOUD_PROJECT_ID"),
			BucketName:      os.Getenv("GOOGLE_CLOUD_BUCKET_NAME"),
			CredentialsFile: os.Getenv("GOOGLE_CLOUD_CREDENTIALS_FILE"),
		},
		AWS: AWSConfig{
			Region:          getEnvOrDefault("AWS_REGION", "us-east-1"),
			BucketName:      os.Getenv("AWS_BUCKET_NAME"),
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWT.SecretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY is required")
	}
	if c.Database.Password == "" {
		log.Println("Warning: DB_PASSWORD is not set")
	}
	return nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
