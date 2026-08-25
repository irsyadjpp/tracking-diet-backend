package repository

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type metabolicMeasurementRepository struct {
	db *gorm.DB
}

func NewMetabolicMeasurementRepository(db *gorm.DB) domain.MetabolicMeasurementRepository {
	return &metabolicMeasurementRepository{db: db}
}

func (r *metabolicMeasurementRepository) Create(mm *domain.MetabolicMeasurement) error {
	return r.db.Create(mm).Error
}

func (r *metabolicMeasurementRepository) GetByID(id int) (*domain.MetabolicMeasurement, error) {
	var mm domain.MetabolicMeasurement
	if err := r.db.First(&mm, id).Error; err != nil {
		return nil, err
	}
	return &mm, nil
}

func (r *metabolicMeasurementRepository) ListByUser(userID int) ([]domain.MetabolicMeasurement, error) {
	var list []domain.MetabolicMeasurement
	if err := r.db.Where("user_id = ?", userID).Order("measured_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *metabolicMeasurementRepository) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.MetabolicMeasurement, error) {
	var list []domain.MetabolicMeasurement
	if err := r.db.Where("user_id = ? AND measured_at >= ? AND measured_at <= ?", 
		userID, startDate, endDate).Order("measured_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *metabolicMeasurementRepository) Update(mm *domain.MetabolicMeasurement) error {
	return r.db.Save(mm).Error
}

func (r *metabolicMeasurementRepository) Delete(id int) error {
	return r.db.Delete(&domain.MetabolicMeasurement{}, id).Error
}
