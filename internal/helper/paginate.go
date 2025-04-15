package helper

import (
	"math"

	"github.com/HasanNugroho/golang-starter/internal/model"
)

func BuildPagination(filter *model.PaginationFilter, totalItems int64) model.Pagination {
	// Hitung total halaman
	totalPages := int(math.Ceil(float64(totalItems) / float64(filter.Limit)))

	// Buat response dengan pagination
	return model.Pagination{
		Limit:      filter.Limit,
		Page:       filter.Page,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
