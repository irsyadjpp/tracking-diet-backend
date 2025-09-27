package repository

import (
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type LabMeasurementRepository struct {
	db *gorm.DB
}

func NewLabMeasurementRepository(db *gorm.DB) domain.LabMeasurementRepository {
	return &LabMeasurementRepository{db: db}
}

func (r *LabMeasurementRepository) CreateLabMeasurement(lab *domain.LabMeasurement) error {
	return r.db.Create(lab).Error
}

func (r *LabMeasurementRepository) GetLabMeasurementByID(id int64) (*domain.LabMeasurement, error) {
	var lab domain.LabMeasurement
	if err := r.db.First(&lab, id).Error; err != nil {
		return nil, err
	}
	return &lab, nil
}

func (r *LabMeasurementRepository) GetLabMeasurementsByUserID(userID int64) ([]*domain.LabMeasurement, error) {
	var labs []*domain.LabMeasurement
	if err := r.db.Where("user_id = ?", userID).Find(&labs).Error; err != nil {
		return nil, err
	}
	return labs, nil
}

func (r *LabMeasurementRepository) UpdateLabMeasurement(lab *domain.LabMeasurement) error {
	return r.db.Save(lab).Error
}

func (r *LabMeasurementRepository) DeleteLabMeasurement(id int64) error {
	return r.db.Delete(&domain.LabMeasurement{}, id).Error
}
