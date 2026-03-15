// internal/utils/response.go

package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OK writes a successful JSON response with HTTP 200.
func OK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// Created writes a successful JSON response with HTTP 201.
func Created(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}
