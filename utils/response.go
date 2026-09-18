package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}

func ErrorResponse(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}
