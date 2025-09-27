package repository

import (
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type labTestRepository struct {
	db *gorm.DB
}

func NewLabTestRepository(db *gorm.DB) domain.LabTestRepository {
	return &labTestRepository{db: db}
}

func (r *labTestRepository) Create(lt *domain.LabTest) error {
	return r.db.Create(lt).Error
}

func (r *labTestRepository) GetByID(id int) (*domain.LabTest, error) {
	var lt domain.LabTest
	if err := r.db.First(&lt, id).Error; err != nil {
		return nil, err
	}
	return &lt, nil
}

func (r *labTestRepository) ListByUser(userID int) ([]domain.LabTest, error) {
	var list []domain.LabTest
	if err := r.db.Where("user_id = ?", userID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *labTestRepository) Update(lt *domain.LabTest) error {
	return r.db.Save(lt).Error
}

func (r *labTestRepository) Delete(id int) error {
	return r.db.Delete(&domain.LabTest{}, id).Error
}
