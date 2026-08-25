package usecase

import (
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

// GamificationUsecase handles gamification business logic
type GamificationUsecase struct {
	achievementRepo domain.AchievementRepository
	userRepo        domain.UserRepository
	bodyRepo        domain.BodyMeasurementRepository
	nutritionRepo   domain.NutritionMeasurementRepository
	fitnessRepo     domain.FitnessMeasurementRepository
}

// NewGamificationUsecase creates a new gamification use case
func NewGamificationUsecase(achievementRepo domain.AchievementRepository, userRepo domain.UserRepository, bodyRepo domain.BodyMeasurementRepository, nutritionRepo domain.NutritionMeasurementRepository, fitnessRepo domain.FitnessMeasurementRepository) *GamificationUsecase {
	return &GamificationUsecase{
		achievementRepo: achievementRepo,
		userRepo:        userRepo,
		bodyRepo:        bodyRepo,
		nutritionRepo:   nutritionRepo,
		fitnessRepo:     fitnessRepo,
	}
}

// CheckAndGrantAchievements checks user's progress and grants achievements
func (uc *GamificationUsecase) CheckAndGrantAchievements(userID int) error {
	// Check for streak achievement
	if err := uc.checkStreakAchievement(userID); err != nil {
		return fmt.Errorf("failed to check streak achievement: %w", err)
	}

	// Check for weight loss achievement
	if err := uc.checkWeightLossAchievement(userID); err != nil {
		return fmt.Errorf("failed to check weight loss achievement: %w", err)
	}

	// Check for workout achievement
	if err := uc.checkWorkoutAchievement(userID); err != nil {
		return fmt.Errorf("failed to check workout achievement: %w", err)
	}

	// Check for nutrition achievement
	if err := uc.checkNutritionAchievement(userID); err != nil {
		return fmt.Errorf("failed to check nutrition achievement: %w", err)
	}

	return nil
}

// GetUserAchievements retrieves all achievements for a user
func (uc *GamificationUsecase) GetUserAchievements(userID int) ([]domain.UserAchievement, error) {
	achievements, err := uc.achievementRepo.GetUserAchievements(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user achievements: %w", err)
	}
	return achievements, nil
}

// GetAvailableAchievements retrieves all available achievements
func (uc *GamificationUsecase) GetAvailableAchievements() ([]domain.Achievement, error) {
	achievements, err := uc.achievementRepo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get available achievements: %w", err)
	}
	return achievements, nil
}

// GetLeaderboard retrieves the leaderboard
func (uc *GamificationUsecase) GetLeaderboard(limit int) ([]domain.Leaderboard, error) {
	leaderboard, err := uc.achievementRepo.GetLeaderboard(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	return leaderboard, nil
}

// AddPoints adds points to a user's leaderboard score
func (uc *GamificationUsecase) AddPoints(userID int, points int) error {
	if err := uc.achievementRepo.UpdateLeaderboard(userID, points); err != nil {
		return fmt.Errorf("failed to add points: %w", err)
	}
	return nil
}

// CreateAchievement creates a new achievement (admin function)
func (uc *GamificationUsecase) CreateAchievement(name, description string, achievementType domain.AchievementType, icon string, points int, tier string, isPublic bool) (*domain.Achievement, error) {
	achievement := &domain.Achievement{
		Name:        name,
		Description: description,
		Type:        achievementType,
		Icon:        icon,
		Points:      points,
		Tier:        tier,
		IsPublic:    isPublic,
	}

	if err := uc.achievementRepo.Create(achievement); err != nil {
		return nil, fmt.Errorf("failed to create achievement: %w", err)
	}

	return achievement, nil
}

// Helper functions for achievement checking
func (uc *GamificationUsecase) checkStreakAchievement(userID int) error {
	// Calculate activity streak
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	measurements, err := uc.bodyRepo.ListByUserWithDateRange(userID, thirtyDaysAgo, time.Now())
	if err != nil {
		return err
	}

	// Count unique days with measurements
	uniqueDays := make(map[string]bool)
	for _, m := range measurements {
		dateKey := m.MeasuredAt.Format("2006-01-02")
		uniqueDays[dateKey] = true
	}

	streakDays := len(uniqueDays)

	// Grant streak achievements based on streak length
	if streakDays >= 7 {
		// Grant 7-day streak achievement
		uc.grantAchievementIfNotEarned(userID, "7 Day Streak")
	}
	if streakDays >= 30 {
		// Grant 30-day streak achievement
		uc.grantAchievementIfNotEarned(userID, "30 Day Streak")
	}

	return nil
}

func (uc *GamificationUsecase) checkWeightLossAchievement(userID int) error {
	// Get body measurements to check weight loss
	measurements, err := uc.bodyRepo.ListByUser(userID)
	if err != nil || len(measurements) < 2 {
		return nil
	}

	// Calculate weight loss from first to last measurement
	firstWeight := measurements[len(measurements)-1].WeightKg
	lastWeight := measurements[0].WeightKg

	if firstWeight != nil && lastWeight != nil {
		weightLoss := *firstWeight - *lastWeight

		// Grant weight loss achievements
		if weightLoss >= 1.0 {
			uc.grantAchievementIfNotEarned(userID, "First Kilogram Lost")
		}
		if weightLoss >= 5.0 {
			uc.grantAchievementIfNotEarned(userID, "5 Kilograms Lost")
		}
		if weightLoss >= 10.0 {
			uc.grantAchievementIfNotEarned(userID, "10 Kilograms Lost")
		}
	}

	return nil
}

func (uc *GamificationUsecase) checkWorkoutAchievement(userID int) error {
	// Get fitness measurements to check workout count
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	workouts, err := uc.fitnessRepo.ListByUserWithDateRange(userID, oneWeekAgo, time.Now())
	if err != nil {
		return nil
	}

	workoutCount := len(workouts)

	// Grant workout achievements
	if workoutCount >= 3 {
		uc.grantAchievementIfNotEarned(userID, "Weekly Warrior")
	}
	if workoutCount >= 5 {
		uc.grantAchievementIfNotEarned(userID, "Gym Rat")
	}

	return nil
}

func (uc *GamificationUsecase) checkNutritionAchievement(userID int) error {
	// Get nutrition measurements to check logging consistency
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	entries, err := uc.nutritionRepo.ListByUserWithDateRange(userID, oneWeekAgo, time.Now())
	if err != nil {
		return nil
	}

	// Count unique days with nutrition entries
	uniqueDays := make(map[string]bool)
	for _, n := range entries {
		dateKey := n.MeasuredAt.Format("2006-01-02")
		uniqueDays[dateKey] = true
	}

	loggingDays := len(uniqueDays)

	// Grant nutrition achievements
	if loggingDays >= 7 {
		uc.grantAchievementIfNotEarned(userID, "Perfect Week Logger")
	}
	if loggingDays >= 30 {
		uc.grantAchievementIfNotEarned(userID, "Consistency Master")
	}

	return nil
}

func (uc *GamificationUsecase) grantAchievementIfNotEarned(userID int, achievementName string) error {
	// Find achievement by name
	achievements, err := uc.achievementRepo.ListAll()
	if err != nil {
		return err
	}

	var achievementID int
	for _, a := range achievements {
		if a.Name == achievementName {
			achievementID = a.ID
			break
		}
	}

	if achievementID == 0 {
		return fmt.Errorf("achievement not found: %s", achievementName)
	}

	// Check if user already has this achievement
	userAchievements, err := uc.achievementRepo.GetUserAchievements(userID)
	if err != nil {
		return err
	}

	for _, ua := range userAchievements {
		if ua.AchievementID == achievementID {
			return nil // Already earned
		}
	}

	// Grant the achievement
	if err := uc.achievementRepo.GrantToUser(userID, achievementID); err != nil {
		return err
	}

	// Add points to leaderboard
	achievement, err := uc.achievementRepo.GetByID(achievementID)
	if err != nil {
		return err
	}

	return uc.AddPoints(userID, achievement.Points)
}