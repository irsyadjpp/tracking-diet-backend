package repository

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type exportRepository struct {
	db *gorm.DB
}

func NewExportRepository(db *gorm.DB) domain.ExportRepository {
	return &exportRepository{db: db}
}

func (r *exportRepository) Create(export *domain.Export) error {
	return r.db.Create(export).Error
}

func (r *exportRepository) GetByID(id int) (*domain.Export, error) {
	var export domain.Export
	if err := r.db.First(&export, id).Error; err != nil {
		return nil, err
	}
	return &export, nil
}

func (r *exportRepository) GetByUser(userID int) ([]domain.Export, error) {
	var exports []domain.Export
	if err := r.db.Where("user_id = ?", userID).Order("requested_at DESC").Find(&exports).Error; err != nil {
		return nil, err
	}
	return exports, nil
}

func (r *exportRepository) Update(export *domain.Export) error {
	return r.db.Save(export).Error
}

func (r *exportRepository) Delete(id int) error {
	return r.db.Delete(&domain.Export{}, id).Error
}

func (r *exportRepository) GetPendingExports() ([]domain.Export, error) {
	var exports []domain.Export
	if err := r.db.Where("status = ?", "pending").Find(&exports).Error; err != nil {
		return nil, err
	}
	return exports, nil
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) domain.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(report *domain.Report) error {
	return r.db.Create(report).Error
}

func (r *reportRepository) GetByID(id int) (*domain.Report, error) {
	var report domain.Report
	if err := r.db.First(&report, id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) GetByUser(userID int) ([]domain.Report, error) {
	var reports []domain.Report
	if err := r.db.Where("user_id = ?", userID).Order("generated_at DESC").Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *reportRepository) GetByUserAndType(userID int, reportType domain.ReportType) ([]domain.Report, error) {
	var reports []domain.Report
	if err := r.db.Where("user_id = ? AND report_type = ?", userID, reportType).Order("generated_at DESC").Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *reportRepository) Update(report *domain.Report) error {
	return r.db.Save(report).Error
}

func (r *reportRepository) Delete(id int) error {
	return r.db.Delete(&domain.Report{}, id).Error
}