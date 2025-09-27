package repository

import (
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type aiRecommendationRepository struct {
	db *gorm.DB
}

func NewAIRecommendationRepository(db *gorm.DB) domain.AIRecommendationRepository {
	return &aiRecommendationRepository{db: db}
}

func (r *aiRecommendationRepository) Create(ar *domain.AIRecommendation) error {
	return r.db.Create(ar).Error
}

func (r *aiRecommendationRepository) GetByID(id int) (*domain.AIRecommendation, error) {
	var ar domain.AIRecommendation
	if err := r.db.First(&ar, id).Error; err != nil {
		return nil, err
	}
	return &ar, nil
}

func (r *aiRecommendationRepository) ListByUser(userID int) ([]domain.AIRecommendation, error) {
	var list []domain.AIRecommendation
	if err := r.db.Where("user_id = ?", userID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *aiRecommendationRepository) Delete(id int) error {
	return r.db.Delete(&domain.AIRecommendation{}, id).Error
}
