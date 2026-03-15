// internal/auth/middleware.go

package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ClaimsKey  = "claims"
	SubjectKey = "subject"
)

// CORSMiddleware sets permissive CORS headers suitable for development.
// In production, replace the wildcard origin with your actual domain.
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Salt")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// RequestLogger logs the method, path, status code, and latency of each request.
func RequestLogger() gin.HandlerFunc {
	return gin.Logger()
}

// JWTAuthMiddleware verifies the Bearer token in the Authorization header.
// The token salt must be supplied via the X-Salt request header.
// On success, the parsed *Claims are stored in the Gin context under ClaimsKey.
func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is missing",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header format must be: Bearer <token>",
			})
			return
		}
		tokenString := parts[1]

		salt := ctx.GetHeader("X-Salt")
		if salt == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "X-Salt header is required",
			})
			return
		}

		claims, err := ValidateToken(jwtSecret, tokenString, salt)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		ctx.Set(ClaimsKey, claims)
		ctx.Set(SubjectKey, claims.Subject)
		ctx.Next()
	}
}
