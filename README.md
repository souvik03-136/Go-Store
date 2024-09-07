# Cloud-Based File Storage System

## Overview

This project is a cloud-based file storage system built using GoLang, designed to allow users to upload, download, and share files securely. The system leverages cloud storage services like Amazon S3 or Google Cloud Storage for file storage and uses PostgreSQL or MySQL to manage file metadata. JWT (JSON Web Tokens) is used for user authentication and authorization to ensure secure access to the stored files.

## Features

- **User Authentication and Authorization:** Secure user authentication using JWT, ensuring that only authorized users can access the system.
- **File Upload and Download:** Users can upload files to the cloud storage and download them as needed.
- **File Sharing:** Users can share files with other users by managing permissions.
- **Metadata Management:** The system stores metadata about the files, such as file names, owners, and permissions, in a relational database.
- **Scalability and Security:** Designed to scale with cloud storage and implement security best practices for file management.

## Tech Stack

- **Backend:** GoLang with the Gin framework
- **File Storage:** Amazon S3 or Google Cloud Storage
- **Database:** PostgreSQL or MySQL
- **Authentication:** JWT for secure user authentication and authorization
- **Containerization:** Docker for containerizing the application
- **Task Automation:** Taskfile for automating common tasks (alternative to Makefile)

## Project Structure

```plaintext
cloud-file-storage/
│
├── cmd/
│   └── api/
│       └── main.go              # Entry point for the API server
│
├── internal/
│   ├── auth/
│   │   ├── jwt.go                # JWT generation and verification
│   │   └── middleware.go         # Authentication and authorization middleware
│   │
│   ├── config/
│   │   └── config.go             # Configuration loading (env variables, etc.)
│   │
│   ├── controllers/
│   │   ├── auth_controller.go    # Handlers for user registration and login
│   │   ├── file_controller.go    # Handlers for file upload, download, and sharing
│   │   └── user_controller.go    # Handlers for user management
│   │
│   ├── models/
│   │   ├── file.go               # File metadata model
│   │   ├── user.go               # User model
│   │   └── permission.go         # File permission model
│   │
│   ├── repository/
│   │   ├── file_repository.go    # Data access layer for file metadata
│   │   ├── user_repository.go    # Data access layer for users
│   │   └── permission_repository.go  # Data access layer for permissions
│   │
│   ├── server/
│   │   ├── routes.go             # API routes definition
│   │   └── server.go             # Server setup and initialization
│   │
│   ├── utils
│   │   └── base_response.go
│   │
│   ├── merrors
│   │   ├── conflict_409.go
│   │   ├── constants.go
│   │   ├── downstream_550.go
│   │   ├── forbidden_403.go
│   │   ├── handle_service_errors.go
│   │   ├── internal_server_500.go
│   │   ├── not_found_404.go
│   │   ├── service_unavailable_503.go
│   │   ├── unauthorized_401.go
│   │   └── validation_422.go
│   │
│   ├── services/
│   │   ├── auth_service.go       # Business logic for authentication
│   │   ├── file_service.go       # Business logic for file operations
│   │   └── user_service.go       # Business logic for user operations
│   │
│   └── storage/
│       ├── s3_storage.go         # Integration with Amazon S3
│       └── gcs_storage.go        # Integration with Google Cloud Storage
│
├── scripts/
│   ├── migrate.sh                # Script for running database migrations
│   └── deploy.sh                 # Script for deploying the application
│
├── database/
│   ├── migrations/
│   │   ├── 001_create_users_table.up.sql   # SQL migration file for creating users table
│   │   ├── 002_create_files_table.up.sql   # SQL migration file for creating files table
│   │   └── 003_create_permissions_table.up.sql  # SQL migration file for permissions
│   │
│   └── queries/
│       ├── file_queries.sql      # SQL queries for file operations
│       ├── user_queries.sql      # SQL queries for user operations
│       └── permission_queries.sql  # SQL queries for permissions
│
├── .env                          # Environment variables (database URL, API keys, etc.)
├── .gitignore                    # Ignored files for Git
├── Dockerfile                    # Docker configuration for containerizing the app
├── docker-compose.yml            # Docker Compose file for managing multi-container setups
├── go.mod                        # Go module file
├── go.sum                        # Go module dependencies
├── Taskfile.yml                  # Taskfile for automating tasks (alternative to Makefile)
└── README.md                     # Project documentation
```

## Setup and Installation

### Prerequisites

- GoLang (version 1.18+)
- Docker (for containerization)
- PostgreSQL or MySQL (for metadata storage)
- AWS account with S3 configured or Google Cloud account with Cloud Storage configured

### Steps

1. **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/cloud-file-storage.git
    cd cloud-file-storage
    ```

2. **Set up environment variables:**
    Create a `.env` file in the root directory with the following variables:
    ```env
    DB_URL="your-database-url"
    JWT_SECRET="your-jwt-secret"
    AWS_ACCESS_KEY_ID="your-aws-access-key-id"
    AWS_SECRET_ACCESS_KEY="your-aws-secret-access-key"
    AWS_S3_BUCKET="your-s3-bucket-name"
    GCP_CREDENTIALS_PATH="path-to-your-gcp-credentials-json"
    ```

3. **Run database migrations:**
    ```bash
    ./scripts/migrate.sh
    ```

4. **Build and run the application:**
    ```bash
    go run cmd/api/main.go
    ```

5. **Access the API:**
    The API will be available at `http://localhost:8080`.

## Usage

### Auth Routes

- **Register OAuth User:**
    ```http
    POST /v1/auth/oauth/register
    ```
    Send a POST request with OAuth user details to register a new user.

- **Login OAuth User:**
    ```http
    POST /v1/auth/oauth/login
    ```
    Send a POST request with OAuth credentials to log in and receive a JWT token.

- **Register Anonymous User:**
    ```http
    POST /v1/auth/anonymous/register
    ```
    Send a POST request to register an anonymous user.

- **Logout User:**
    ```http
    POST /v1/auth/logout
    ```
    Send a POST request to log out the user and invalidate the JWT token.

- **Validate Token:**
    ```http
    GET /v1/auth/validate
    ```
    Send a GET request to validate the JWT token.

### User Routes

- **Create a New User:**
    ```http
    POST /v1/users
    ```
    Send a POST request with user details to create a new user.

- **Get a User by ID:**
    ```http
    GET /v1/users/:id
    ```
    Send a GET request with the user ID to retrieve user details.

- **Update a User by ID:**
    ```http
    PUT /v1/users/:id
    ```
    Send a PUT request with user details to update the user.

- **Delete a User by ID:**
    ```http
    DELETE /v1/users/:id
    ```
    Send a DELETE request with the user ID to remove the user.

### File Routes

- **Create a New File:**
    ```http
    POST /v1/files
    ```
    Send a POST request with the file and metadata to create a new file.

- **Get a File by ID:**
    ```http
    GET /v1/files/:id
    ```
    Send a GET request with the file ID to retrieve file details.

- **Update a File by ID:**
    ```http
    PUT /v1/files/:id
    ```
    Send a PUT request with updated file details to modify the file.

- **Delete a File by ID:**
    ```http
    DELETE /v1/files/:id
    ```
    Send a DELETE request with the file ID to remove the file.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.

## Acknowledgments

- Inspired by modern cloud storage systems and best practices in cloud architecture.
