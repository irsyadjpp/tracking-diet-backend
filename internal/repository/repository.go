package repository

import (
	"time"

	"gorm.io/gorm"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
}

// PaginatedResult represents a paginated result
type PaginatedResult struct {
	Data       interface{}
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// Paginate applies pagination to a GORM query
func Paginate(db *gorm.DB, params PaginationParams) *gorm.DB {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	offset := (params.Page - 1) * params.PageSize
	return db.Offset(offset).Limit(params.PageSize)
}

// ApplyDateRangeFilter applies a date range filter to a GORM query
func ApplyDateRangeFilter(db *gorm.DB, dateField string, startDate, endDate *time.Time) *gorm.DB {
	if startDate != nil {
		db = db.Where(dateField+" >= ?", *startDate)
	}
	if endDate != nil {
		db = db.Where(dateField+" <= ?", *endDate)
	}
	return db
}

// BatchDelete performs a batch delete with the given conditions
func BatchDelete(db *gorm.DB, model interface{}, conditions map[string]interface{}) error {
	return db.Where(conditions).Delete(model).Error
}

// GetTotalCount returns the total count of records for a query
func GetTotalCount(db *gorm.DB, model interface{}) (int64, error) {
	var count int64
	if err := db.Model(model).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}