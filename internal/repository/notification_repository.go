package repository

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *domain.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) GetByID(id int) (*domain.Notification, error) {
	var notification domain.Notification
	if err := r.db.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) GetByUser(userID int) ([]domain.Notification, error) {
	var notifications []domain.Notification
	if err := r.db.Where("user_id = ?", userID).Order("sent_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) GetUnread(userID int) ([]domain.Notification, error) {
	var notifications []domain.Notification
	if err := r.db.Where("user_id = ? AND read_at IS NULL", userID).Order("sent_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) MarkAsRead(id int) error {
	now := time.Now()
	return r.db.Model(&domain.Notification{}).Where("id = ?", id).Update("read_at", now).Error
}

func (r *notificationRepository) Delete(id int) error {
	return r.db.Delete(&domain.Notification{}, id).Error
}

func (r *notificationRepository) GetPreferences(userID int) (*domain.NotificationPreference, error) {
	var pref domain.NotificationPreference
	if err := r.db.Where("user_id = ?", userID).First(&pref).Error; err != nil {
		return nil, err
	}
	return &pref, nil
}

func (r *notificationRepository) UpdatePreferences(pref *domain.NotificationPreference) error {
	return r.db.Save(pref).Error
}

func (r *notificationRepository) CreateDefaultPreferences(userID int) error {
	pref := &domain.NotificationPreference{
		UserID:                 userID,
		EmailEnabled:           true,
		PushEnabled:            true,
		SMSEnabled:             false,
		InAppEnabled:           true,
		GoalRemindersEnabled:   true,
		MealRemindersEnabled:   true,
		WeeklyReportsEnabled:   true,
	}
	return r.db.Create(pref).Error
}