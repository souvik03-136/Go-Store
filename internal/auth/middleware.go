package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/merrors"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) settings.
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// RequestLogger logs the details of each incoming request.
func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Log the incoming request details
		ctx.Next()
		// Optionally log response details if needed
	}
}

// JWTAuthMiddleware verifies the JWT token and extracts claims.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			merrors.Unauthorized(ctx, "Authorization header is missing")
			return
		}

		// Extract the token from the header
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			merrors.Unauthorized(ctx, "Invalid authorization header format")
			return
		}
		tokenString := tokenParts[1]

		// Extract the salt (assuming it's sent as a query parameter, header, or some other way)
		salt := ctx.Query("salt")
		if salt == "" {
			merrors.Unauthorized(ctx, "Salt is missing")
			return
		}

		// Validate the token
		claims, err := ValidateToken(ctx, tokenString, salt)
		if err != nil {
			merrors.Unauthorized(ctx, "Invalid token")
			return
		}

		// Attach the claims to the context for use in the handlers
		ctx.Set("claims", claims)

		ctx.Next()
	}
}
