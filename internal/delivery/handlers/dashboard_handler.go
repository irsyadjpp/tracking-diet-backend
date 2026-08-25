package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/dto"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/middleware"
	"github.com/irsyadjpp/tracking-diet-backend/internal/usecase"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

// DashboardHandler handles dashboard HTTP requests
type DashboardHandler struct {
	dashboardUC *usecase.DashboardUsecase
	logger      *logger.Logger
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardUC *usecase.DashboardUsecase, log *logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		dashboardUC: dashboardUC,
		logger:      log,
	}
}

// GetStats handles getting dashboard statistics
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	stats, err := h.dashboardUC.GetDashboardStats(userID)
	if err != nil {
		h.logger.Error("Failed to get dashboard stats", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to get dashboard stats", "INTERNAL_ERROR")
		return
	}

	response := h.toStatsResponse(stats)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// GetProgressTrend handles getting progress trend data
func (h *DashboardHandler) GetProgressTrend(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	metricType := r.URL.Query().Get("metric_type")
	if metricType == "" {
		metricType = "weight" // default
	}

	daysStr := r.URL.Query().Get("days")
	days := 30 // default
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	trend, err := h.dashboardUC.GetProgressTrend(userID, metricType, days)
	if err != nil {
		h.logger.Error("Failed to get progress trend", err, map[string]interface{}{"user_id": userID, "metric_type": metricType})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to get progress trend", "INTERNAL_ERROR")
		return
	}

	response := h.toTrendResponse(trend)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// GetWeeklySummary handles getting weekly summary
func (h *DashboardHandler) GetWeeklySummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	summary, err := h.dashboardUC.GetWeeklySummary(userID)
	if err != nil {
		h.logger.Error("Failed to get weekly summary", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to get weekly summary", "INTERNAL_ERROR")
		return
	}

	response := h.toWeeklySummaryResponse(summary)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// GetGoalProgress handles getting goal progress
func (h *DashboardHandler) GetGoalProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	progress, err := h.dashboardUC.GetGoalProgress(userID)
	if err != nil {
		h.logger.Error("Failed to get goal progress", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to get goal progress", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.GoalProgressResponse, len(progress))
	for i, p := range progress {
		response[i] = h.toGoalProgressResponse(&p)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// GetComparisonReport handles getting comparison report
func (h *DashboardHandler) GetComparisonReport(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "week_over_week" // default
	}

	report, err := h.dashboardUC.GetComparisonReport(userID, period)
	if err != nil {
		h.logger.Error("Failed to get comparison report", err, map[string]interface{}{"user_id": userID, "period": period})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to get comparison report", "INTERNAL_ERROR")
		return
	}

	response := h.toComparisonReportResponse(report)
	middleware.JSONResponder(w, http.StatusOK, response)
}

func (h *DashboardHandler) toStatsResponse(stats *usecase.DashboardStats) dto.DashboardStatsResponse {
	return dto.DashboardStatsResponse{
		UserID:            stats.UserID,
		CurrentWeight:     stats.CurrentWeight,
		WeightChange:      stats.WeightChange,
		CurrentBMI:        stats.CurrentBMI,
		HealthScore:       stats.HealthScore,
		ActivityStreak:     stats.ActivityStreak,
		MealsLoggedToday:   stats.MealsLoggedToday,
		WorkoutsThisWeek:   stats.WorkoutsThisWeek,
		SleepQualityAvg:    stats.SleepQualityAvg,
		UpdatedAt:          stats.UpdatedAt,
	}
}

func (h *DashboardHandler) toTrendResponse(trend *usecase.ProgressTrend) dto.ProgressTrendResponse {
	dataPoints := make([]dto.TrendPointResponse, len(trend.DataPoints))
	for i, dp := range trend.DataPoints {
		dataPoints[i] = dto.TrendPointResponse{
			Date:  dp.Date,
			Value: dp.Value,
			Label: dp.Label,
		}
	}

	return dto.ProgressTrendResponse{
		UserID:     trend.UserID,
		MetricType: trend.MetricType,
		StartDate:  trend.StartDate,
		EndDate:    trend.EndDate,
		DataPoints: dataPoints,
	}
}

func (h *DashboardHandler) toWeeklySummaryResponse(summary *usecase.WeeklySummary) dto.WeeklySummaryResponse {
	return dto.WeeklySummaryResponse{
		UserID:          summary.UserID,
		WeekStart:       summary.WeekStart,
		WeekEnd:         summary.WeekEnd,
		AverageCalories:  summary.AverageCalories,
		AverageSteps:     summary.AverageSteps,
		WeightChange:     summary.WeightChange,
		WorkoutCount:     summary.WorkoutCount,
		MealCompliance:   summary.MealCompliance,
		SleepQualityAvg:  summary.SleepQualityAvg,
		MoodAvg:          summary.MoodAvg,
		Achievements:     summary.Achievements,
		Recommendations:  summary.Recommendations,
	}
}

func (h *DashboardHandler) toGoalProgressResponse(progress *usecase.GoalProgress) dto.GoalProgressResponse {
	return dto.GoalProgressResponse{
		UserID:               progress.UserID,
		GoalType:             progress.GoalType,
		TargetValue:          progress.TargetValue,
		CurrentValue:         progress.CurrentValue,
		ProgressPercent:      progress.ProgressPercent,
		StartDate:            progress.StartDate,
		TargetDate:           progress.TargetDate,
		Status:               progress.Status,
		EstimatedCompletion: progress.EstimatedCompletion,
	}
}

func (h *DashboardHandler) toComparisonReportResponse(report *usecase.ComparisonReport) dto.ComparisonReportResponse {
	return dto.ComparisonReportResponse{
		UserID:            report.UserID,
		ComparisonPeriod:  report.ComparisonPeriod,
		WeightChange:      report.WeightChange,
		BMIChange:         report.BMIChange,
		CalorieChange:     report.CalorieChange,
		ActivityChange:    report.ActivityChange,
		HealthScoreChange: report.HealthScoreChange,
		Improvements:      report.Improvements,
		Concerns:          report.Concerns,
	}
}