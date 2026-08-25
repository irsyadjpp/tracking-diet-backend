package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

// NotificationUsecase handles notification business logic
type NotificationUsecase struct {
	notificationRepo domain.NotificationRepository
}

// NewNotificationUsecase creates a new notification use case
func NewNotificationUsecase(notificationRepo domain.NotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{
		notificationRepo: notificationRepo,
	}
}

// SendNotification sends a notification to a user
func (uc *NotificationUsecase) SendNotification(userID int, notificationType domain.NotificationType, title, message string, data interface{}, channels []domain.NotificationChannel, priority domain.NotificationPriority, actionURL, actionLabel *string) (*domain.Notification, error) {
	// Get user preferences to determine which channels to use
	pref, err := uc.notificationRepo.GetPreferences(userID)
	if err != nil {
		// Create default preferences if not found
		if err := uc.notificationRepo.CreateDefaultPreferences(userID); err != nil {
			return nil, fmt.Errorf("failed to create default preferences: %w", err)
		}
		pref, err = uc.notificationRepo.GetPreferences(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get preferences: %w", err)
		}
	}

	// Filter channels based on user preferences
	activeChannels := uc.filterChannels(channels, pref, notificationType)
	if len(activeChannels) == 0 {
		// Always send in-app notification if enabled
		if pref.InAppEnabled {
			activeChannels = []domain.NotificationChannel{domain.NotificationChannelInApp}
		} else {
			return nil, errors.New("no enabled notification channels")
		}
	}

	notification := &domain.Notification{
		UserID:      userID,
		Type:        notificationType,
		Title:       title,
		Message:     message,
		Data:        data,
		Channels:    activeChannels,
		Priority:    priority,
		SentAt:      time.Now(),
		ActionURL:   actionURL,
		ActionLabel: actionLabel,
	}

	if err := uc.notificationRepo.Create(notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Send notifications through each channel (implementation would go here)
	// For now, we just create the notification record
	for _, channel := range activeChannels {
		uc.sendThroughChannel(notification, channel)
	}

	return notification, nil
}

// GetNotifications retrieves all notifications for a user
func (uc *NotificationUsecase) GetNotifications(userID int) ([]domain.Notification, error) {
	notifications, err := uc.notificationRepo.GetByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	return notifications, nil
}

// GetUnreadNotifications retrieves unread notifications for a user
func (uc *NotificationUsecase) GetUnreadNotifications(userID int) ([]domain.Notification, error) {
	notifications, err := uc.notificationRepo.GetUnread(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %w", err)
	}
	return notifications, nil
}

// MarkAsRead marks a notification as read
func (uc *NotificationUsecase) MarkAsRead(id int) error {
	if err := uc.notificationRepo.MarkAsRead(id); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

// DeleteNotification deletes a notification
func (uc *NotificationUsecase) DeleteNotification(id int) error {
	if err := uc.notificationRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// GetPreferences retrieves notification preferences for a user
func (uc *NotificationUsecase) GetPreferences(userID int) (*domain.NotificationPreference, error) {
	pref, err := uc.notificationRepo.GetPreferences(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}
	return pref, nil
}

// UpdatePreferences updates notification preferences for a user
func (uc *NotificationUsecase) UpdatePreferences(userID int, emailEnabled, pushEnabled, smsEnabled, inAppEnabled *bool, goalRemindersEnabled, mealRemindersEnabled, weeklyReportsEnabled *bool, quietHoursStart, quietHoursEnd *int) (*domain.NotificationPreference, error) {
	pref, err := uc.notificationRepo.GetPreferences(userID)
	if err != nil {
		// Create default preferences if not found
		if err := uc.notificationRepo.CreateDefaultPreferences(userID); err != nil {
			return nil, fmt.Errorf("failed to create default preferences: %w", err)
		}
		pref, err = uc.notificationRepo.GetPreferences(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get preferences: %w", err)
		}
	}

	if emailEnabled != nil {
		pref.EmailEnabled = *emailEnabled
	}
	if pushEnabled != nil {
		pref.PushEnabled = *pushEnabled
	}
	if smsEnabled != nil {
		pref.SMSEnabled = *smsEnabled
	}
	if inAppEnabled != nil {
		pref.InAppEnabled = *inAppEnabled
	}
	if goalRemindersEnabled != nil {
		pref.GoalRemindersEnabled = *goalRemindersEnabled
	}
	if mealRemindersEnabled != nil {
		pref.MealRemindersEnabled = *mealRemindersEnabled
	}
	if weeklyReportsEnabled != nil {
		pref.WeeklyReportsEnabled = *weeklyReportsEnabled
	}
	if quietHoursStart != nil {
		pref.QuietHoursStart = quietHoursStart
	}
	if quietHoursEnd != nil {
		pref.QuietHoursEnd = quietHoursEnd
	}

	if err := uc.notificationRepo.UpdatePreferences(pref); err != nil {
		return nil, fmt.Errorf("failed to update preferences: %w", err)
	}

	return pref, nil
}

// SendGoalReminder sends a goal reminder notification
func (uc *NotificationUsecase) SendGoalReminder(userID int, goalTitle string, progressPercent float64) (*domain.Notification, error) {
	title := "Goal Progress Reminder"
	message := fmt.Sprintf("Your goal '%s' is %.1f%% complete. Keep up the good work!", goalTitle, progressPercent)
	data := map[string]interface{}{
		"goal_title":         goalTitle,
		"progress_percent":   progressPercent,
		"notification_type":  "goal_reminder",
	}

	return uc.SendNotification(userID, domain.NotificationTypeGoalReminder, title, message, data, []domain.NotificationChannel{domain.NotificationChannelInApp, domain.NotificationChannelPush}, domain.NotificationPriorityMedium, nil, nil)
}

// SendAchievementNotification sends an achievement notification
func (uc *NotificationUsecase) SendAchievementNotification(userID int, achievementName, description string) (*domain.Notification, error) {
	title := "Achievement Unlocked!"
	message := fmt.Sprintf("Congratulations! You've unlocked: %s", achievementName)
	data := map[string]interface{}{
		"achievement_name":  achievementName,
		"description":       description,
		"notification_type": "achievement",
	}

	return uc.SendNotification(userID, domain.NotificationTypeAchievement, title, message, data, []domain.NotificationChannel{domain.NotificationChannelInApp, domain.NotificationChannelPush, domain.NotificationChannelEmail}, domain.NotificationPriorityHigh, nil, nil)
}

// SendMealReminder sends a meal reminder notification
func (uc *NotificationUsecase) SendMealReminder(userID int, mealType string) (*domain.Notification, error) {
	title := "Meal Reminder"
	message := fmt.Sprintf("Time to log your %s! Remember to track your nutrition.", mealType)
	data := map[string]interface{}{
		"meal_type":         mealType,
		"notification_type": "meal_reminder",
	}

	return uc.SendNotification(userID, domain.NotificationTypeMealReminder, title, message, data, []domain.NotificationChannel{domain.NotificationChannelInApp, domain.NotificationChannelPush}, domain.NotificationPriorityLow, nil, nil)
}

// SendWeeklyReport sends a weekly report notification
func (uc *NotificationUsecase) SendWeeklyReport(userID int, summary interface{}) (*domain.Notification, error) {
	title := "Your Weekly Progress Report"
	message := "Check out your weekly progress summary to see how you're doing!"
	data := map[string]interface{}{
		"summary":           summary,
		"notification_type": "weekly_report",
	}

	return uc.SendNotification(userID, domain.NotificationTypeWeeklyReport, title, message, data, []domain.NotificationChannel{domain.NotificationChannelInApp, domain.NotificationChannelEmail}, domain.NotificationPriorityMedium, nil, nil)
}

// Helper functions
func (uc *NotificationUsecase) filterChannels(channels []domain.NotificationChannel, pref *domain.NotificationPreference, notificationType domain.NotificationType) []domain.NotificationChannel {
	var filtered []domain.NotificationChannel

	for _, channel := range channels {
		switch channel {
		case domain.NotificationChannelEmail:
			if pref.EmailEnabled {
				filtered = append(filtered, channel)
			}
		case domain.NotificationChannelPush:
			if pref.PushEnabled {
				filtered = append(filtered, channel)
			}
		case domain.NotificationChannelSMS:
			if pref.SMSEnabled {
				filtered = append(filtered, channel)
			}
		case domain.NotificationChannelInApp:
			if pref.InAppEnabled {
				filtered = append(filtered, channel)
			}
		}
	}

	// Check type-specific preferences
	switch notificationType {
	case domain.NotificationTypeGoalReminder:
		if !pref.GoalRemindersEnabled {
			return []domain.NotificationChannel{}
		}
	case domain.NotificationTypeMealReminder:
		if !pref.MealRemindersEnabled {
			return []domain.NotificationChannel{}
		}
	case domain.NotificationTypeWeeklyReport:
		if !pref.WeeklyReportsEnabled {
			return []domain.NotificationChannel{}
		}
	}

	return filtered
}

func (uc *NotificationUsecase) sendThroughChannel(notification *domain.Notification, channel domain.NotificationChannel) {
	// Placeholder for actual notification sending logic
	// This would integrate with email service, push notification service, SMS service, etc.
	switch channel {
	case domain.NotificationChannelEmail:
		// Send email notification
		fmt.Printf("Sending email notification to user %d: %s\n", notification.UserID, notification.Title)
	case domain.NotificationChannelPush:
		// Send push notification
		fmt.Printf("Sending push notification to user %d: %s\n", notification.UserID, notification.Title)
	case domain.NotificationChannelSMS:
		// Send SMS notification
		fmt.Printf("Sending SMS notification to user %d: %s\n", notification.UserID, notification.Title)
	case domain.NotificationChannelInApp:
		// In-app notification is handled by the frontend
		fmt.Printf("Creating in-app notification for user %d: %s\n", notification.UserID, notification.Title)
	}
}