package usecase

import (
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

// ExportUsecase handles data export business logic
type ExportUsecase struct {
	exportRepo domain.ExportRepository
	bodyRepo   domain.BodyMeasurementRepository
	nutritionRepo domain.NutritionMeasurementRepository
	metabolicRepo domain.MetabolicMeasurementRepository
	fitnessRepo domain.FitnessMeasurementRepository
	wellbeingRepo domain.WellbeingMeasurementRepository
	labTestRepo domain.LabTestRepository
	goalRepo domain.GoalRepository
}

// NewExportUsecase creates a new export use case
func NewExportUsecase(exportRepo domain.ExportRepository, bodyRepo domain.BodyMeasurementRepository, nutritionRepo domain.NutritionMeasurementRepository, metabolicRepo domain.MetabolicMeasurementRepository, fitnessRepo domain.FitnessMeasurementRepository, wellbeingRepo domain.WellbeingMeasurementRepository, labTestRepo domain.LabTestRepository, goalRepo domain.GoalRepository) *ExportUsecase {
	return &ExportUsecase{
		exportRepo: exportRepo,
		bodyRepo: bodyRepo,
		nutritionRepo: nutritionRepo,
		metabolicRepo: metabolicRepo,
		fitnessRepo: fitnessRepo,
		wellbeingRepo: wellbeingRepo,
		labTestRepo: labTestRepo,
		goalRepo: goalRepo,
	}
}

// RequestExport creates a new export request
func (uc *ExportUsecase) RequestExport(userID int, exportType domain.ExportType, format domain.ExportFormat, startDate, endDate *time.Time) (*domain.Export, error) {
	export := &domain.Export{
		UserID:     userID,
		ExportType: exportType,
		Format:     format,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     "pending",
	}

	if err := uc.exportRepo.Create(export); err != nil {
		return nil, fmt.Errorf("failed to create export request: %w", err)
	}

	// In a real implementation, this would trigger a background job to process the export
	// For now, we'll mark it as completed immediately
	export.Status = "completed"
	now := time.Now()
	export.CompletedAt = &now
	
	// Set expiration to 7 days from now
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	export.ExpiresAt = &expiresAt

	if err := uc.exportRepo.Update(export); err != nil {
		return nil, fmt.Errorf("failed to update export status: %w", err)
	}

	return export, nil
}

// GetExport retrieves an export by ID
func (uc *ExportUsecase) GetExport(id int) (*domain.Export, error) {
	export, err := uc.exportRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get export: %w", err)
	}
	return export, nil
}

// GetUserExports retrieves all exports for a user
func (uc *ExportUsecase) GetUserExports(userID int) ([]domain.Export, error) {
	exports, err := uc.exportRepo.GetByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user exports: %w", err)
	}
	return exports, nil
}

// DeleteExport deletes an export
func (uc *ExportUsecase) DeleteExport(id int) error {
	if err := uc.exportRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete export: %w", err)
	}
	return nil
}

// ReportUsecase handles report generation business logic
type ReportUsecase struct {
	reportRepo domain.ReportRepository
	bodyRepo   domain.BodyMeasurementRepository
	nutritionRepo domain.NutritionMeasurementRepository
	metabolicRepo domain.MetabolicMeasurementRepository
}

// NewReportUsecase creates a new report use case
func NewReportUsecase(reportRepo domain.ReportRepository, bodyRepo domain.BodyMeasurementRepository, nutritionRepo domain.NutritionMeasurementRepository, metabolicRepo domain.MetabolicMeasurementRepository) *ReportUsecase {
	return &ReportUsecase{
		reportRepo: reportRepo,
		bodyRepo: bodyRepo,
		nutritionRepo: nutritionRepo,
		metabolicRepo: metabolicRepo,
	}
}

// GenerateReport generates a new report
func (uc *ReportUsecase) GenerateReport(userID int, reportType domain.ReportType, title, description string, startDate, endDate time.Time) (*domain.Report, error) {
	// Gather data based on report type
	var data interface{}
	
	switch reportType {
	case domain.ReportTypeProgress:
		data = uc.generateProgressData(userID, startDate, endDate)
	case domain.ReportTypeHealthScore:
		data = uc.generateHealthScoreData(userID, startDate, endDate)
	case domain.ReportTypeNutrition:
		data = uc.generateNutritionData(userID, startDate, endDate)
	case domain.ReportTypeComprehensive:
		data = uc.generateComprehensiveData(userID, startDate, endDate)
	default:
		data = map[string]interface{}{"message": "Report type not implemented"}
	}

	report := &domain.Report{
		UserID:      userID,
		ReportType:  reportType,
		Title:       title,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Data:        data,
	}

	if err := uc.reportRepo.Create(report); err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}

// GetReport retrieves a report by ID
func (uc *ReportUsecase) GetReport(id int) (*domain.Report, error) {
	report, err := uc.reportRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	return report, nil
}

// GetUserReports retrieves all reports for a user
func (uc *ReportUsecase) GetUserReports(userID int) ([]domain.Report, error) {
	reports, err := uc.reportRepo.GetByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reports: %w", err)
	}
	return reports, nil
}

// GetReportsByType retrieves reports of a specific type for a user
func (uc *ReportUsecase) GetReportsByType(userID int, reportType domain.ReportType) ([]domain.Report, error) {
	reports, err := uc.reportRepo.GetByUserAndType(userID, reportType)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by type: %w", err)
	}
	return reports, nil
}

// DeleteReport deletes a report
func (uc *ReportUsecase) DeleteReport(id int) error {
	if err := uc.reportRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}
	return nil
}

// Helper functions for data generation
func (uc *ReportUsecase) generateProgressData(userID int, startDate, endDate time.Time) interface{} {
	bodyMeasurements, _ := uc.bodyRepo.ListByUserWithDateRange(userID, startDate, endDate)
	
	return map[string]interface{}{
		"body_measurements": bodyMeasurements,
		"total_records":     len(bodyMeasurements),
		"period": map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}
}

func (uc *ReportUsecase) generateHealthScoreData(userID int, startDate, endDate time.Time) interface{} {
	metabolicMeasurements, _ := uc.metabolicRepo.ListByUserWithDateRange(userID, startDate, endDate)
	
	return map[string]interface{}{
		"metabolic_measurements": metabolicMeasurements,
		"total_records":           len(metabolicMeasurements),
		"period": map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}
}

func (uc *ReportUsecase) generateNutritionData(userID int, startDate, endDate time.Time) interface{} {
	nutritionMeasurements, _ := uc.nutritionRepo.ListByUserWithDateRange(userID, startDate, endDate)
	
	return map[string]interface{}{
		"nutrition_measurements": nutritionMeasurements,
		"total_records":           len(nutritionMeasurements),
		"period": map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}
}

func (uc *ReportUsecase) generateComprehensiveData(userID int, startDate, endDate time.Time) interface{} {
	bodyMeasurements, _ := uc.bodyRepo.ListByUserWithDateRange(userID, startDate, endDate)
	nutritionMeasurements, _ := uc.nutritionRepo.ListByUserWithDateRange(userID, startDate, endDate)
	metabolicMeasurements, _ := uc.metabolicRepo.ListByUserWithDateRange(userID, startDate, endDate)
	
	return map[string]interface{}{
		"body_measurements":     bodyMeasurements,
		"nutrition_measurements": nutritionMeasurements,
		"metabolic_measurements": metabolicMeasurements,
		"total_records": map[string]int{
			"body":      len(bodyMeasurements),
			"nutrition": len(nutritionMeasurements),
			"metabolic": len(metabolicMeasurements),
		},
		"period": map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}
}