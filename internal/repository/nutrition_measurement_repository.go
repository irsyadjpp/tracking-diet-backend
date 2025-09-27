package repository

import (
    "github.com/irsyadjpp/tracking-diet-backend/internal/domain"
    "gorm.io/gorm"
)

type nutritionMeasurementRepository struct {
    db *gorm.DB
}

func NewNutritionMeasurementRepository(db *gorm.DB) domain.NutritionMeasurementRepository {
    return &nutritionMeasurementRepository{db: db}
}

func (r *nutritionMeasurementRepository) Create(nm *domain.NutritionMeasurement) error {
    return r.db.Create(nm).Error
}

func (r *nutritionMeasurementRepository) GetByID(id int) (*domain.NutritionMeasurement, error) {
    var nm domain.NutritionMeasurement
    if err := r.db.First(&nm, id).Error; err != nil {
        return nil, err
    }
    return &nm, nil
}

func (r *nutritionMeasurementRepository) ListByUser(userID int) ([]domain.NutritionMeasurement, error) {
    var list []domain.NutritionMeasurement
    if err := r.db.Where("user_id = ?", userID).Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func (r *nutritionMeasurementRepository) Update(nm *domain.NutritionMeasurement) error {
    return r.db.Save(nm).Error
}

func (r *nutritionMeasurementRepository) Delete(id int) error {
    return r.db.Delete(&domain.NutritionMeasurement{}, id).Error
}
