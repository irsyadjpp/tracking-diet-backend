package repository

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type goalRepository struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) domain.GoalRepository {
	return &goalRepository{db: db}
}

func (r *goalRepository) Create(goal *domain.Goal) error {
	return r.db.Create(goal).Error
}

func (r *goalRepository) GetByID(id int) (*domain.Goal, error) {
	var goal domain.Goal
	if err := r.db.First(&goal, id).Error; err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *goalRepository) GetByUser(userID int) ([]domain.Goal, error) {
	var goals []domain.Goal
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *goalRepository) GetActiveGoals(userID int) ([]domain.Goal, error) {
	var goals []domain.Goal
	if err := r.db.Where("user_id = ? AND status = ?", userID, domain.GoalStatusActive).Order("target_date ASC").Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *goalRepository) Update(goal *domain.Goal) error {
	return r.db.Save(goal).Error
}

func (r *goalRepository) Delete(id int) error {
	return r.db.Delete(&domain.Goal{}, id).Error
}

func (r *goalRepository) UpdateProgress(goalID int, value float64) error {
	progress := &domain.GoalProgress{
		GoalID:          goalID,
		RecordedAt:      time.Now(),
		Value:           value,
		MeasurementType: "automatic",
	}
	return r.db.Create(progress).Error
}

func (r *goalRepository) GetProgressHistory(goalID int) ([]domain.GoalProgress, error) {
	var progress []domain.GoalProgress
	if err := r.db.Where("goal_id = ?", goalID).Order("recorded_at DESC").Find(&progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *goalRepository) CreateReminder(reminder *domain.GoalReminder) error {
	return r.db.Create(reminder).Error
}

func (r *goalRepository) GetReminders(goalID int) ([]domain.GoalReminder, error) {
	var reminders []domain.GoalReminder
	if err := r.db.Where("goal_id = ?", goalID).Find(&reminders).Error; err != nil {
		return nil, err
	}
	return reminders, nil
}