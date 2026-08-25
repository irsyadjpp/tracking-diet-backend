package repository

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) domain.DeviceRepository {
	return &deviceRepository{db: db}
}

// Device management
func (r *deviceRepository) CreateDevice(device *domain.Device) error {
	return r.db.Create(device).Error
}

func (r *deviceRepository) GetByID(id int) (*domain.Device, error) {
	var device domain.Device
	if err := r.db.First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) GetByUser(userID int) ([]domain.Device, error) {
	var devices []domain.Device
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) GetByUserAndType(userID int, deviceType domain.DeviceType) (*domain.Device, error) {
	var device domain.Device
	if err := r.db.Where("user_id = ? AND type = ?", userID, deviceType).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) UpdateDevice(device *domain.Device) error {
	return r.db.Save(device).Error
}

func (r *deviceRepository) DeleteDevice(id int) error {
	return r.db.Delete(&domain.Device{}, id).Error
}

func (r *deviceRepository) UpdateDeviceStatus(id int, status domain.DeviceStatus) error {
	return r.db.Model(&domain.Device{}).Where("id = ?", id).Update("status", status).Error
}

func (r *deviceRepository) UpdateSyncTime(id int, lastSyncAt, nextSyncAt time.Time) error {
	return r.db.Model(&domain.Device{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_sync_at": lastSyncAt,
		"next_sync_at": nextSyncAt,
	}).Error
}

// Data management
func (r *deviceRepository) StoreDeviceData(data *domain.DeviceData) error {
	return r.db.Create(data).Error
}

func (r *deviceRepository) GetDeviceData(deviceID int, dataType domain.DataType, startDate, endDate time.Time) ([]domain.DeviceData, error) {
	var deviceData []domain.DeviceData
	query := r.db.Where("device_id = ?", deviceID)
	
	if dataType != "" {
		query = query.Where("data_type = ?", dataType)
	}
	
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("recorded_at >= ? AND recorded_at <= ?", startDate, endDate)
	}
	
	if err := query.Order("recorded_at DESC").Find(&deviceData).Error; err != nil {
		return nil, err
	}
	return deviceData, nil
}

func (r *deviceRepository) GetUnprocessedData(deviceID int) ([]domain.DeviceData, error) {
	var deviceData []domain.DeviceData
	if err := r.db.Where("device_id = ? AND is_processed = ?", deviceID, false).Order("recorded_at ASC").Find(&deviceData).Error; err != nil {
		return nil, err
	}
	return deviceData, nil
}

func (r *deviceRepository) MarkDataAsProcessed(id int) error {
	return r.db.Model(&domain.DeviceData{}).Where("id = ?", id).Update("is_processed", true).Error
}

func (r *deviceRepository) DeleteOldData(olderThan time.Time) error {
	return r.db.Where("recorded_at < ?", olderThan).Delete(&domain.DeviceData{}).Error
}

// Sync logs
func (r *deviceRepository) CreateSyncLog(log *domain.SyncLog) error {
	return r.db.Create(log).Error
}

func (r *deviceRepository) GetSyncLogs(deviceID int, limit int) ([]domain.SyncLog, error) {
	var logs []domain.SyncLog
	if err := r.db.Where("device_id = ?", deviceID).Order("started_at DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *deviceRepository) GetLatestSyncLog(deviceID int) (*domain.SyncLog, error) {
	var log domain.SyncLog
	if err := r.db.Where("device_id = ?", deviceID).Order("started_at DESC").First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}