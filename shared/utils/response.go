package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type PaginatedResponse[T any] struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Data       []T    `json:"data"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalCount int64  `json:"total_count"`
	TotalPages int    `json:"total_pages"`
}

func RespondSuccess[T any](ctx *gin.Context, statusCode int, message string, data T) {
	ctx.JSON(statusCode, SuccessResponse[T]{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func RespondError(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

func RespondErrorWithDetails(ctx *gin.Context, statusCode int, message string, err string) {
	ctx.JSON(statusCode, ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   err,
	})
}

func RespondPaginated[T any](ctx *gin.Context, message string, data []T, page, pageSize int, totalCount int64) {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, PaginatedResponse[T]{
		Status:     "success",
		Message:    message,
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	})
}
