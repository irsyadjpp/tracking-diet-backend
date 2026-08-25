package usecase

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

// DashboardUsecase handles dashboard business logic
type DashboardUsecase struct {
	dashboardRepo domain.DashboardRepository
}

// NewDashboardUsecase creates a new dashboard use case
func NewDashboardUsecase(dashboardRepo domain.DashboardRepository) *DashboardUsecase {
	return &DashboardUsecase{
		dashboardRepo: dashboardRepo,
	}
}

// GetDashboardStats retrieves comprehensive dashboard statistics
func (uc *DashboardUsecase) GetDashboardStats(userID int) (*domain.DashboardStats, error) {
	stats, err := uc.dashboardRepo.GetDashboardStats(userID)
	if err != nil {
		return nil, err
	}

	// Add business logic calculations
	if stats.WeightChange != nil {
		// Determine if weight change is positive or negative
		// Negative = weight loss (usually good for weight loss goals)
		// Positive = weight gain (could be good for muscle gain goals)
	}

	return stats, nil
}

// GetProgressTrend retrieves trend data for a specific metric
func (uc *DashboardUsecase) GetProgressTrend(userID int, metricType string, days int) (*domain.ProgressTrend, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	trend, err := uc.dashboardRepo.GetProgressTrend(userID, metricType, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Add trend analysis
	if len(trend.DataPoints) > 1 {
		firstValue := trend.DataPoints[0].Value
		lastValue := trend.DataPoints[len(trend.DataPoints)-1].Value
		change := lastValue - firstValue
		changePercent := (change / firstValue) * 100

		// Add trend label
		if changePercent > 5 {
			trend.DataPoints[len(trend.DataPoints)-1].Label = "Significant increase"
		} else if changePercent < -5 {
			trend.DataPoints[len(trend.DataPoints)-1].Label = "Significant decrease"
		} else {
			trend.DataPoints[len(trend.DataPoints)-1].Label = "Stable"
		}
	}

	return trend, nil
}

// GetWeeklySummary retrieves weekly progress summary
func (uc *DashboardUsecase) GetWeeklySummary(userID int) (*domain.WeeklySummary, error) {
	weekStart := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
	
	summary, err := uc.dashboardRepo.GetWeeklySummary(userID, weekStart)
	if err != nil {
		return nil, err
	}

	// Add business logic for achievements and recommendations
	summary.Achievements = uc.generateAchievements(summary)
	summary.Recommendations = uc.generateRecommendations(summary)

	return summary, nil
}

// GetGoalProgress retrieves progress towards user's goals
func (uc *DashboardUsecase) GetGoalProgress(userID int) ([]domain.GoalProgress, error) {
	progress, err := uc.dashboardRepo.GetGoalProgress(userID)
	if err != nil {
		return nil, err
	}

	// Update status based on progress
	for i := range progress {
		if progress[i].ProgressPercent >= 100 {
			progress[i].Status = "completed"
		} else if progress[i].ProgressPercent >= 80 {
			progress[i].Status = "ahead"
		} else if progress[i].ProgressPercent >= 50 {
			progress[i].Status = "on_track"
		} else {
			progress[i].Status = "behind"
		}
	}

	return progress, nil
}

// GetComparisonReport retrieves comparison with previous periods
func (uc *DashboardUsecase) GetComparisonReport(userID int, period string) (*domain.ComparisonReport, error) {
	report, err := uc.dashboardRepo.GetComparisonReport(userID, period)
	if err != nil {
		return nil, err
	}

	// Add business logic for improvements and concerns
	report.Improvements = uc.analyzeImprovements(report)
	report.Concerns = uc.analyzeConcerns(report)

	return report, nil
}

// Helper functions
func (uc *DashboardUsecase) generateAchievements(summary *domain.WeeklySummary) []string {
	achievements := []string{}

	if summary.MealCompliance >= 80 {
		achievements = append(achievements, "Great meal logging consistency")
	}

	if summary.WorkoutCount >= 3 {
		achievements = append(achievements, "Active week with regular workouts")
	}

	if summary.SleepQualityAvg != nil && *summary.SleepQualityAvg >= 7 {
		achievements = append(achievements, "Good sleep quality maintained")
	}

	if summary.MoodAvg != nil && *summary.MoodAvg >= 7 {
		achievements = append(achievements, "Positive mood maintained")
	}

	if len(achievements) == 0 {
		achievements = append(achievements, "Keep tracking for achievements")
	}

	return achievements
}

func (uc *DashboardUsecase) generateRecommendations(summary *domain.WeeklySummary) []string {
	recommendations := []string{}

	if summary.MealCompliance < 60 {
		recommendations = append(recommendations, "Try to log meals more consistently")
	}

	if summary.WorkoutCount < 2 {
		recommendations = append(recommendations, "Consider adding more physical activity")
	}

	if summary.SleepQualityAvg != nil && *summary.SleepQualityAvg < 6 {
		recommendations = append(recommendations, "Focus on improving sleep quality")
	}

	if summary.AverageCalories != nil && *summary.AverageCalories > 2500 {
		recommendations = append(recommendations, "Review calorie intake if weight loss is a goal")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Continue with your current routine")
	}

	return recommendations
}

func (uc *DashboardUsecase) analyzeImprovements(report *domain.ComparisonReport) []string {
	improvements := []string{}

	if report.WeightChange != nil && *report.WeightChange < 0 {
		improvements = append(improvements, "Weight loss progress")
	}

	if report.ActivityChange != nil && *report.ActivityChange > 0 {
		improvements = append(improvements, "Increased activity level")
	}

	if report.HealthScoreChange != nil && *report.HealthScoreChange > 0 {
		improvements = append(improvements, "Health score improvement")
	}

	if len(improvements) == 0 {
		improvements = append(improvements, "Maintaining current progress")
	}

	return improvements
}

func (uc *DashboardUsecase) analyzeConcerns(report *domain.ComparisonReport) []string {
	concerns := []string{}

	if report.WeightChange != nil && *report.WeightChange > 1 {
		concerns = append(concerns, "Weight gain detected")
	}

	if report.CalorieChange != nil && *report.CalorieChange > 500 {
		concerns = append(concerns, "Significant calorie increase")
	}

	if report.HealthScoreChange != nil && *report.HealthScoreChange < -5 {
		concerns = append(concerns, "Health score decline")
	}

	if len(concerns) == 0 {
		concerns = append(concerns, "No major concerns detected")
	}

	return concerns
}