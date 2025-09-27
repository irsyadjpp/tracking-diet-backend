package repository

import (
    "github.com/irsyadjpp/tracking-diet-backend/internal/domain"
    "gorm.io/gorm"
)

type bodyMeasurementRepository struct {
    db *gorm.DB
}

func NewBodyMeasurementRepository(db *gorm.DB) domain.BodyMeasurementRepository {
    return &bodyMeasurementRepository{db: db}
}

func (r *bodyMeasurementRepository) Create(bm *domain.BodyMeasurement) error {
    return r.db.Create(bm).Error
}

func (r *bodyMeasurementRepository) GetByID(id int) (*domain.BodyMeasurement, error) {
    var bm domain.BodyMeasurement
    if err := r.db.First(&bm, id).Error; err != nil {
        return nil, err
    }
    return &bm, nil
}

func (r *bodyMeasurementRepository) ListByUser(userID int) ([]domain.BodyMeasurement, error) {
    var list []domain.BodyMeasurement
    if err := r.db.Where("user_id = ?", userID).Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func (r *bodyMeasurementRepository) Update(bm *domain.BodyMeasurement) error {
    return r.db.Save(bm).Error
}

func (r *bodyMeasurementRepository) Delete(id int) error {
    return r.db.Delete(&domain.BodyMeasurement{}, id).Error
}
