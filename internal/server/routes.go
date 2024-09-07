package server

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"                              // or another appropriate driver for your DB
	"github.com/souvik03-136/Go-Store/internal/config" // Assuming you have a config package
	"github.com/souvik03-136/Go-Store/internal/controllers"
	"github.com/souvik03-136/Go-Store/internal/repository"
)

func InitRoutes(router *gin.Engine) {
	// Initialize database connection (PostgreSQL in this example)
	connStr := "user=your_user password=your_password dbname=your_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Don't forget to close the DB connection when the application stops
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)

	// Initialize configuration (assuming you have a config structure)
	cfg, err := config.LoadConfig() // You should implement this function to load your config
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize controllers
	userController := controllers.NewUserController(userRepo)

	// Handle the two return values from NewFileController
	fileController, err := controllers.NewFileController(fileRepo, cfg)
	if err != nil {
		log.Fatalf("Could not create file controller: %v", err)
	}

	// Auth routes
	router.POST("/v1/auth/oauth/register", controllers.RegisterOAuthUser)
	router.POST("/v1/auth/oauth/login", controllers.LoginOAuthUser)
	router.POST("/v1/auth/anonymous/register", controllers.RegisterAnonymousUser)
	router.POST("/v1/auth/logout", controllers.LogoutUser)
	router.GET("/v1/auth/validate", controllers.ValidateTokenHandler)

	// User routes
	router.POST("/v1/users", userController.CreateUser)       // Create a new user
	router.GET("/v1/users/:id", userController.GetUserByID)   // Get a user by ID
	router.PUT("/v1/users/:id", userController.UpdateUser)    // Update a user by ID
	router.DELETE("/v1/users/:id", userController.DeleteUser) // Delete a user by ID

	// File routes
	router.POST("/v1/files", fileController.CreateFile)       // Create a new file
	router.GET("/v1/files/:id", fileController.GetFileByID)   // Get a file by ID
	router.PUT("/v1/files/:id", fileController.UpdateFile)    // Update a file by ID
	router.DELETE("/v1/files/:id", fileController.DeleteFile) // Delete a file by ID
}
