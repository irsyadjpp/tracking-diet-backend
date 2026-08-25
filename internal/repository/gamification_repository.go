package repository

import (
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type achievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) domain.AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) Create(achievement *domain.Achievement) error {
	return r.db.Create(achievement).Error
}

func (r *achievementRepository) GetByID(id int) (*domain.Achievement, error) {
	var achievement domain.Achievement
	if err := r.db.First(&achievement, id).Error; err != nil {
		return nil, err
	}
	return &achievement, nil
}

func (r *achievementRepository) ListAll() ([]domain.Achievement, error) {
	var achievements []domain.Achievement
	if err := r.db.Where("is_public = ?", true).Order("tier DESC, points DESC").Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

func (r *achievementRepository) ListByType(achievementType domain.AchievementType) ([]domain.Achievement, error) {
	var achievements []domain.Achievement
	if err := r.db.Where("type = ? AND is_public = ?", achievementType, true).Order("tier DESC, points DESC").Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

func (r *achievementRepository) Update(achievement *domain.Achievement) error {
	return r.db.Save(achievement).Error
}

func (r *achievementRepository) Delete(id int) error {
	return r.db.Delete(&domain.Achievement{}, id).Error
}

func (r *achievementRepository) GrantToUser(userID, achievementID int) error {
	userAchievement := &domain.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
		IsCompleted:   true,
	}
	return r.db.Create(userAchievement).Error
}

func (r *achievementRepository) GetUserAchievements(userID int) ([]domain.UserAchievement, error) {
	var userAchievements []domain.UserAchievement
	if err := r.db.Where("user_id = ?", userID).Preload("Achievement").Order("earned_at DESC").Find(&userAchievements).Error; err != nil {
		return nil, err
	}
	return userAchievements, nil
}

func (r *achievementRepository) GetLeaderboard(limit int) ([]domain.Leaderboard, error) {
	var leaderboard []domain.Leaderboard
	if err := r.db.Order("total_points DESC").Limit(limit).Find(&leaderboard).Error; err != nil {
		return nil, err
	}
	return leaderboard, nil
}

func (r *achievementRepository) UpdateLeaderboard(userID int, points int) error {
	var leaderboard domain.Leaderboard
	if err := r.db.Where("user_id = ?", userID).First(&leaderboard).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new leaderboard entry
			leaderboard = domain.Leaderboard{
				UserID:      userID,
				TotalPoints: points,
				Rank:        0,
			}
			return r.db.Create(&leaderboard).Error
		}
		return err
	}

	// Update existing entry
	leaderboard.TotalPoints += points
	return r.db.Save(&leaderboard).Error
}