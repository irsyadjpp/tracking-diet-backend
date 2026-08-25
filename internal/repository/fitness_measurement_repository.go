package repository

import (
    "time"

    "github.com/irsyadjpp/tracking-diet-backend/internal/domain"
    "gorm.io/gorm"
)

type fitnessMeasurementRepository struct {
    db *gorm.DB
}

func NewFitnessMeasurementRepository(db *gorm.DB) domain.FitnessMeasurementRepository {
    return &fitnessMeasurementRepository{db: db}
}

func (r *fitnessMeasurementRepository) Create(fm *domain.FitnessMeasurement) error {
    return r.db.Create(fm).Error
}

func (r *fitnessMeasurementRepository) GetByID(id int) (*domain.FitnessMeasurement, error) {
    var fm domain.FitnessMeasurement
    if err := r.db.First(&fm, id).Error; err != nil {
        return nil, err
    }
    return &fm, nil
}

func (r *fitnessMeasurementRepository) ListByUser(userID int) ([]domain.FitnessMeasurement, error) {
    var list []domain.FitnessMeasurement
    if err := r.db.Where("user_id = ?", userID).Order("measured_at DESC").Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func (r *fitnessMeasurementRepository) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.FitnessMeasurement, error) {
    var list []domain.FitnessMeasurement
    if err := r.db.Where("user_id = ? AND measured_at >= ? AND measured_at <= ?", 
        userID, startDate, endDate).Order("measured_at DESC").Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func (r *fitnessMeasurementRepository) Update(fm *domain.FitnessMeasurement) error {
    return r.db.Save(fm).Error
}

func (r *fitnessMeasurementRepository) Delete(id int) error {
    return r.db.Delete(&domain.FitnessMeasurement{}, id).Error
}
