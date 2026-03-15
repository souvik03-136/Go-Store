// internal/merrors/merrors.go

// Package merrors provides helpers for returning consistent JSON error responses
// from Gin handlers. Every helper writes the response and calls ctx.Abort() so
// the handler chain stops immediately.
package merrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// errorBody is the JSON shape returned for all error responses.
type errorBody struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func respond(ctx *gin.Context, code int, errType, message string) {
	ctx.AbortWithStatusJSON(code, errorBody{
		Error: apiError{Code: code, Type: errType, Message: message},
	})
}

// BadRequest responds with HTTP 400.
func BadRequest(ctx *gin.Context, message string) {
	respond(ctx, http.StatusBadRequest, "bad_request", message)
}

// Unauthorized responds with HTTP 401.
func Unauthorized(ctx *gin.Context, message string) {
	respond(ctx, http.StatusUnauthorized, "unauthorized", message)
}

// Forbidden responds with HTTP 403.
func Forbidden(ctx *gin.Context, message string) {
	respond(ctx, http.StatusForbidden, "forbidden", message)
}

// NotFound responds with HTTP 404.
func NotFound(ctx *gin.Context, message string) {
	respond(ctx, http.StatusNotFound, "not_found", message)
}

// Conflict responds with HTTP 409.
func Conflict(ctx *gin.Context, message string) {
	respond(ctx, http.StatusConflict, "conflict", message)
}

// Validation responds with HTTP 422.
func Validation(ctx *gin.Context, message string) {
	respond(ctx, http.StatusUnprocessableEntity, "validation_error", message)
}

// InternalServer responds with HTTP 500.
func InternalServer(ctx *gin.Context, message string) {
	respond(ctx, http.StatusInternalServerError, "internal_server_error", message)
}

// ServiceUnavailable responds with HTTP 503.
func ServiceUnavailable(ctx *gin.Context, message string) {
	respond(ctx, http.StatusServiceUnavailable, "service_unavailable", message)
}

// Downstream responds with a non-standard HTTP 550 for upstream dependency failures.
func Downstream(ctx *gin.Context, message string) {
	respond(ctx, 550, "downstream_error", message)
}
