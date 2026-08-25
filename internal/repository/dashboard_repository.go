package repository

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) domain.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetDashboardStats(userID int) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{
		UserID: userID,
	}

	// Get current weight (most recent body measurement)
	var bodyMeasurement domain.BodyMeasurement
	if err := r.db.Where("user_id = ?", userID).Order("measured_at DESC").First(&bodyMeasurement).Error; err == nil {
		stats.CurrentWeight = bodyMeasurement.WeightKg
		stats.CurrentBMI = bodyMeasurement.BMI
	}

	// Calculate weight change (compare with 30 days ago)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var oldWeight float64
	r.db.Model(&domain.BodyMeasurement{}).
		Where("user_id = ? AND measured_at >= ?", userID, thirtyDaysAgo).
		Order("measured_at ASC").
		Select("weight_kg").
		Scan(&oldWeight)

	if stats.CurrentWeight != nil && oldWeight > 0 {
		change := *stats.CurrentWeight - oldWeight
		stats.WeightChange = &change
	}

	// Calculate health score (simplified - would use actual metabolic data)
	stats.HealthScore = r.calculateHealthScore(userID)

	// Count activity streak (consecutive days with measurements)
	stats.ActivityStreak = r.calculateActivityStreak(userID)

	// Count meals logged today
	today := time.Now().Truncate(24 * time.Hour)
	r.db.Model(&domain.NutritionMeasurement{}).
		Where("user_id = ? AND measured_at >= ?", userID, today).
		Count(&stats.MealsLoggedToday)

	// Count workouts this week
	weekStart := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
	r.db.Model(&domain.FitnessMeasurement{}).
		Where("user_id = ? AND measured_at >= ?", userID, weekStart).
		Count(&stats.WorkoutsThisWeek)

	// Calculate average sleep quality
	var avgSleepQuality float64
	r.db.Model(&domain.WellbeingMeasurement{}).
		Where("user_id = ? AND measured_at >= ?", userID, time.Now().AddDate(0, 0, -7)).
		Select("AVG(sleep_quality)").
		Scan(&avgSleepQuality)
	stats.SleepQualityAvg = &avgSleepQuality

	stats.UpdatedAt = time.Now()
	return stats, nil
}

func (r *dashboardRepository) GetProgressTrend(userID int, metricType string, startDate, endDate time.Time) (*domain.ProgressTrend, error) {
	trend := &domain.ProgressTrend{
		UserID:     userID,
		MetricType: metricType,
		StartDate:  startDate,
		EndDate:    endDate,
		DataPoints: []domain.TrendPoint{},
	}

	switch metricType {
	case "weight":
		var measurements []domain.BodyMeasurement
		r.db.Where("user_id = ? AND measured_at >= ? AND measured_at <= ?", 
			userID, startDate, endDate).
			Order("measured_at ASC").
			Find(&measurements)

		for _, m := range measurements {
			if m.WeightKg != nil {
				trend.DataPoints = append(trend.DataPoints, domain.TrendPoint{
					Date:  m.MeasuredAt,
					Value: *m.WeightKg,
				})
			}
		}

	case "calories":
		var measurements []domain.NutritionMeasurement
		r.db.Where("user_id = ? AND measured_at >= ? AND measured_at <= ?", 
			userID, startDate, endDate).
			Order("measured_at ASC").
			Find(&measurements)

		for _, m := range measurements {
			if m.CaloriesInKcal != nil {
				trend.DataPoints = append(trend.DataPoints, domain.TrendPoint{
					Date:  m.MeasuredAt,
					Value: float64(*m.CaloriesInKcal),
				})
			}
		}

	case "sleep":
		var measurements []domain.WellbeingMeasurement
		r.db.Where("user_id = ? AND measured_at >= ? AND measured_at <= ?", 
			userID, startDate, endDate).
			Order("measured_at ASC").
			Find(&measurements)

		for _, m := range measurements {
			if m.SleepHours != nil {
				trend.DataPoints = append(trend.DataPoints, domain.TrendPoint{
					Date:  m.MeasuredAt,
					Value: *m.SleepHours,
				})
			}
		}
	}

	return trend, nil
}

func (r *dashboardRepository) GetWeeklySummary(userID int, weekStart time.Time) (*domain.WeeklySummary, error) {
	weekEnd := weekStart.AddDate(0, 0, 7)
	summary := &domain.WeeklySummary{
		UserID:    userID,
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
	}

	// Calculate average calories
	var avgCalories float64
	r.db.Model(&domain.NutritionMeasurement{}).
		Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, weekStart, weekEnd).
		Select("AVG(calories_in_kcal)").
		Scan(&avgCalories)
	summary.AverageCalories = &avgCalories

	// Calculate average steps
	var avgSteps int
	r.db.Model(&domain.NutritionMeasurement{}).
		Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, weekStart, weekEnd).
		Select("AVG(step_count)").
		Scan(&avgSteps)
	summary.AverageSteps = &avgSteps

	// Calculate weight change
	var firstWeight, lastWeight float64
	r.db.Model(&domain.BodyMeasurement{}).
		Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, weekStart, weekEnd).
		Order("measured_at ASC").
		Select("weight_kg").
		Scan(&firstWeight)
	
	r.db.Model(&domain.BodyMeasurement{}).
		Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, weekStart, weekEnd).
		Order("measured_at DESC").
		Select("weight_kg").
		Scan(&lastWeight)

	if firstWeight > 0 && lastWeight > 0 {
		summary.WeightChange = &lastWeight
	}

	// Count workouts
	r.db.Model(&domain.FitnessMeasurement{}).
		Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, weekStart, weekEnd).
		Count(&summary.WorkoutCount)

	// Calculate meal compliance (simplified - count days with nutrition entries)
	var daysWithEntries int64
	r.db.Model(&domain.NutritionMeasurement{}).
		Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, weekStart, weekEnd).
		Select("COUNT(DISTINCT measured_at)").
		Scan(&daysWithEntries)
	
	if daysWithEntries > 0 {
		summary.MealCompliance = float64(daysWithEntries) / 7.0 * 100
	}

	// Calculate average sleep quality
	var avgSleepQuality float64
	r.db.Model(&domain.WellbeingMeasurement{}).
		Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, weekStart, weekEnd).
		Select("AVG(sleep_quality)").
		Scan(&avgSleepQuality)
	summary.SleepQualityAvg = &avgSleepQuality

	// Calculate average mood
	var avgMood float64
	r.db.Model(&domain.WellbeingMeasurement{}).
		Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, weekStart, weekEnd).
		Select("AVG(mood_level)").
		Scan(&avgMood)
	summary.MoodAvg = &avgMood

	// Generate achievements and recommendations (simplified)
	summary.Achievements = []string{"Logged meals 5+ days", "Workout at least once"}
	summary.Recommendations = []string{"Increase protein intake", "Track sleep more consistently"}

	return summary, nil
}

func (r *dashboardRepository) GetGoalProgress(userID int) ([]domain.GoalProgress, error) {
	// This is a simplified implementation
	// In a real implementation, you would have a goals table and track progress
	progress := []domain.GoalProgress{
		{
			UserID:          userID,
			GoalType:        "weight_loss",
			TargetValue:     70.0,
			CurrentValue:    75.0,
			ProgressPercent: 50.0,
			StartDate:       time.Now().AddDate(0, 0, -30),
			TargetDate:      time.Now().AddDate(0, 0, 30),
			Status:          "on_track",
		},
	}

	return progress, nil
}

func (r *dashboardRepository) GetComparisonReport(userID int, period string) (*domain.ComparisonReport, error) {
	report := &domain.ComparisonReport{
		UserID:           userID,
		ComparisonPeriod: period,
	}

	var startDate time.Time
	switch period {
	case "week_over_week":
		startDate = time.Now().AddDate(0, 0, -7)
	case "month_over_month":
		startDate = time.Now().AddDate(0, -1, 0)
	default:
		startDate = time.Now().AddDate(0, 0, -7)
	}

	// Calculate weight change
	var currentWeight, oldWeight float64
	r.db.Model(&domain.BodyMeasurement{}).
		Where("user_id = ?", userID).
		Order("measured_at DESC").
		Select("weight_kg").
		Scan(&currentWeight)

	r.db.Model(&domain.BodyMeasurement{}).
		Where("user_id = ? AND measured_at >= ?", userID, startDate).
		Order("measured_at ASC").
		Select("weight_kg").
		Scan(&oldWeight)

	if currentWeight > 0 && oldWeight > 0 {
		change := currentWeight - oldWeight
		report.WeightChange = &change
	}

	// Generate improvements and concerns (simplified)
	report.Improvements = []string{"Activity level increased", "Sleep quality improved"}
	report.Concerns = []string{"Calorie intake slightly high", "Stress levels elevated"}

	return report, nil
}

// Helper functions
func (r *dashboardRepository) calculateHealthScore(userID int) *float64 {
	// Simplified health score calculation
	// In a real implementation, this would use metabolic measurements
	score := 75.0 // Default score
	return &score
}

func (r *dashboardRepository) calculateActivityStreak(userID int) int {
	// Calculate consecutive days with any measurements
	streak := 0
	currentDate := time.Now()

	for i := 0; i < 365; i++ { // Check up to a year back
		checkDate := currentDate.AddDate(0, 0, -i)
		dateStart := checkDate.Truncate(24 * time.Hour)
		dateEnd := dateStart.Add(24 * time.Hour)

		var count int64
		r.db.Model(&domain.BodyMeasurement{}).
			Where("user_id = ? AND measured_at >= ? AND measured_at < ?", userID, dateStart, dateEnd).
			Count(&count)

		if count > 0 {
			streak++
		} else if i > 0 {
			break
		}
	}

	return streak
}