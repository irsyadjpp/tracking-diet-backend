package repository

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type wellbeingMeasurementRepository struct {
	db *gorm.DB
}

func NewWellbeingMeasurementRepository(db *gorm.DB) domain.WellbeingMeasurementRepository {
	return &wellbeingMeasurementRepository{db: db}
}

func (r *wellbeingMeasurementRepository) Create(wm *domain.WellbeingMeasurement) error {
	return r.db.Create(wm).Error
}

func (r *wellbeingMeasurementRepository) GetByID(id int) (*domain.WellbeingMeasurement, error) {
	var wm domain.WellbeingMeasurement
	if err := r.db.First(&wm, id).Error; err != nil {
		return nil, err
	}
	return &wm, nil
}

func (r *wellbeingMeasurementRepository) ListByUser(userID int) ([]domain.WellbeingMeasurement, error) {
	var list []domain.WellbeingMeasurement
	if err := r.db.Where("user_id = ?", userID).Order("measured_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *wellbeingMeasurementRepository) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.WellbeingMeasurement, error) {
	var list []domain.WellbeingMeasurement
	if err := r.db.Where("user_id = ? AND measured_at >= ? AND measured_at <= ?", 
		userID, startDate, endDate).Order("measured_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *wellbeingMeasurementRepository) Update(wm *domain.WellbeingMeasurement) error {
	return r.db.Save(wm).Error
}

func (r *wellbeingMeasurementRepository) Delete(id int) error {
	return r.db.Delete(&domain.WellbeingMeasurement{}, id).Error
}
