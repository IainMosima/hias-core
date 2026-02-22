package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Sort     string `json:"sort"`
	Order    string `json:"order"`
}

func GetPaginationParams(ctx *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	sort := ctx.DefaultQuery("sort", "created_at")
	order := ctx.DefaultQuery("order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	}
}

func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p PaginationParams) Limit() int {
	return p.PageSize
}
