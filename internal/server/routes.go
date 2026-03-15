// internal/server/routes.go

package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/auth"
	"github.com/souvik03-136/Go-Store/internal/controllers"
	"github.com/souvik03-136/Go-Store/internal/repository"
	"github.com/souvik03-136/Go-Store/internal/services"
	"github.com/souvik03-136/Go-Store/internal/storage"
)

// initRoutes wires all repositories, services, controllers, and routes together.
// The db connection is owned by Server and must not be closed here.
func initRoutes(router *gin.Engine, db *sql.DB, store storage.Storage, jwtSecret string) {
	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	authSvc := services.NewAuthService(userRepo, jwtSecret)
	userSvc := services.NewUserService(userRepo)
	fileSvc := services.NewFileService(fileRepo, store)

	// ── Controllers ───────────────────────────────────────────────────────────
	authCtrl := controllers.NewAuthController(authSvc)
	userCtrl := controllers.NewUserController(userSvc)
	fileCtrl := controllers.NewFileController(fileSvc)

	// ── Global middleware ─────────────────────────────────────────────────────
	router.Use(auth.CORSMiddleware())
	router.Use(auth.RequestLogger())

	// ── Health check (unauthenticated) ────────────────────────────────────────
	router.GET("/healthz", func(ctx *gin.Context) {
		if err := db.Ping(); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "db": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ── Auth routes (public) ──────────────────────────────────────────────────
	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/oauth/register", authCtrl.RegisterOAuthUser)
			authGroup.POST("/oauth/login", authCtrl.LoginOAuthUser)
			authGroup.POST("/anonymous/register", authCtrl.RegisterAnonymousUser)
			authGroup.POST("/logout", authCtrl.LogoutUser)
			authGroup.GET("/validate", authCtrl.ValidateTokenHandler)
		}

		// ── Protected routes ──────────────────────────────────────────────────
		protected := v1.Group("")
		protected.Use(auth.JWTAuthMiddleware(jwtSecret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.POST("", userCtrl.CreateUser)
				users.GET("/:id", userCtrl.GetUserByID)
				users.PUT("/:id", userCtrl.UpdateUser)
				users.DELETE("/:id", userCtrl.DeleteUser)
			}

			// File routes
			files := protected.Group("/files")
			{
				files.POST("", fileCtrl.CreateFile)
				files.GET("", fileCtrl.ListMyFiles)
				files.GET("/:id", fileCtrl.GetFileByID)
				files.PUT("/:id", fileCtrl.UpdateFile)
				files.DELETE("/:id", fileCtrl.DeleteFile)
			}
		}
	}
}
