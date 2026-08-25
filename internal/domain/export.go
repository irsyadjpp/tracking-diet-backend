package domain

import "time"

// ExportFormat represents different export formats
type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatPDF  ExportFormat = "pdf"
	ExportFormatXLSX ExportFormat = "xlsx"
)

// ExportType represents different types of data to export
type ExportType string

const (
	ExportTypeBodyMeasurements    ExportType = "body_measurements"
	ExportTypeNutrition           ExportType = "nutrition"
	ExportTypeMetabolic           ExportType = "metabolic"
	ExportTypeFitness             ExportType = "fitness"
	ExportTypeWellbeing           ExportType = "wellbeing"
	ExportTypeLabTests            ExportType = "lab_tests"
	ExportTypeGoals               ExportType = "goals"
	ExportTypeAllData             ExportType = "all_data"
)

// Export represents a data export request
type Export struct {
	ID          int         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int         `gorm:"not null;index" json:"user_id"`
	ExportType  ExportType  `gorm:"size:50;not null" json:"export_type"`
	Format      ExportFormat `gorm:"size:10;not null" json:"format"`
	StartDate   *time.Time  `json:"start_date"`
	EndDate     *time.Time  `json:"end_date"`
	Status      string      `gorm:"size:20;default:pending" json:"status"` // pending, processing, completed, failed
	FileURL     *string     `gorm:"size:500" json:"file_url"`
	FileSize    *int64      `json:"file_size"`
	ExpiresAt   *time.Time  `json:"expires_at"`
	RequestedAt time.Time   `gorm:"autoCreateTime" json:"requested_at"`
	CompletedAt *time.Time  `json:"completed_at"`
}

// TableName specifies the table name for Export model
func (Export) TableName() string {
	return "exports"
}

// ReportType represents different types of reports
type ReportType string

const (
	ReportTypeProgress      ReportType = "progress"
	ReportTypeHealthScore   ReportType = "health_score"
	ReportTypeNutrition     ReportType = "nutrition"
	ReportTypeGoals         ReportType = "goals"
	ReportTypeComprehensive ReportType = "comprehensive"
)

// Report represents a generated report
type Report struct {
	ID          int         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int         `gorm:"not null;index" json:"user_id"`
	ReportType  ReportType  `gorm:"size:50;not null" json:"report_type"`
	Title       string      `gorm:"size:200;not null" json:"title"`
	Description string      `gorm:"type:text" json:"description"`
	StartDate   time.Time   `gorm:"not null" json:"start_date"`
	EndDate     time.Time   `gorm:"not null" json:"end_date"`
	Data        interface{} `gorm:"type:jsonb" json:"data"`
	GeneratedAt time.Time   `gorm:"autoCreateTime" json:"generated_at"`
}

// TableName specifies the table name for Report model
func (Report) TableName() string {
	return "reports"
}

// ExportRepository defines operations for managing data exports
type ExportRepository interface {
	Create(export *Export) error
	GetByID(id int) (*Export, error)
	GetByUser(userID int) ([]Export, error)
	Update(export *Export) error
	Delete(id int) error
	GetPendingExports() ([]Export, error)
}

// ReportRepository defines operations for managing reports
type ReportRepository interface {
	Create(report *Report) error
	GetByID(id int) (*Report, error)
	GetByUser(userID int) ([]Report, error)
	GetByUserAndType(userID int, reportType ReportType) ([]Report, error)
	Update(report *Report) error
	Delete(id int) error
}