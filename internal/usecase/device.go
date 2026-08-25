package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrDeviceNotFound    = errors.New("device not found")
	ErrDeviceNotConnected = errors.New("device not connected")
)

// DeviceUsecase handles wearable device business logic
type DeviceUsecase struct {
	deviceRepo   domain.DeviceRepository
	userRepo     domain.UserRepository
	nutritionRepo domain.NutritionMeasurementRepository
	fitnessRepo  domain.FitnessMeasurementRepository
	wellbeingRepo domain.WellbeingMeasurementRepository
}

// NewDeviceUsecase creates a new device use case
func NewDeviceUsecase(deviceRepo domain.DeviceRepository, userRepo domain.UserRepository, nutritionRepo domain.NutritionMeasurementRepository, fitnessRepo domain.FitnessMeasurementRepository, wellbeingRepo domain.WellbeingMeasurementRepository) *DeviceUsecase {
	return &DeviceUsecase{
		deviceRepo:   deviceRepo,
		userRepo:     userRepo,
		nutritionRepo: nutritionRepo,
		fitnessRepo:  fitnessRepo,
		wellbeingRepo: wellbeingRepo,
	}
}

// ConnectDevice connects a new wearable device
func (uc *DeviceUsecase) ConnectDevice(userID int, deviceType domain.DeviceType, name, serialNumber, model, firmwareVersion, accessToken, refreshToken string, tokenExpiresAt *time.Time) (*domain.Device, error) {
	// Check if device of this type already exists
	existing, err := uc.deviceRepo.GetByUserAndType(userID, deviceType)
	if err == nil && existing != nil {
		return nil, errors.New("device of this type already connected")
	}

	device := &domain.Device{
		UserID:          userID,
		Type:            deviceType,
		Name:            name,
		SerialNumber:    serialNumber,
		Model:           model,
		FirmwareVersion: firmwareVersion,
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		TokenExpiresAt:  tokenExpiresAt,
		Status:          domain.DeviceStatusConnected,
		SyncFrequency:   3600, // Default 1 hour
		IsEnabled:       true,
	}

	if err := uc.deviceRepo.CreateDevice(device); err != nil {
		return nil, fmt.Errorf("failed to connect device: %w", err)
	}

	return device, nil
}

// GetDevice retrieves a device by ID
func (uc *DeviceUsecase) GetDevice(id int) (*domain.Device, error) {
	device, err := uc.deviceRepo.GetByID(id)
	if err != nil {
		return nil, ErrDeviceNotFound
	}
	return device, nil
}

// GetUserDevices retrieves all devices for a user
func (uc *DeviceUsecase) GetUserDevices(userID int) ([]domain.Device, error) {
	devices, err := uc.deviceRepo.GetByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user devices: %w", err)
	}
	return devices, nil
}

// DisconnectDevice disconnects a wearable device
func (uc *DeviceUsecase) DisconnectDevice(id int) error {
	if err := uc.deviceRepo.UpdateDeviceStatus(id, domain.DeviceStatusDisconnected); err != nil {
		return fmt.Errorf("failed to disconnect device: %w", err)
	}
	return nil
}

// UpdateDevice updates device information
func (uc *DeviceUsecase) UpdateDevice(id int, name *string, accessToken, refreshToken *string, tokenExpiresAt *time.Time, isEnabled *bool) (*domain.Device, error) {
	device, err := uc.deviceRepo.GetByID(id)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	if name != nil {
		device.Name = *name
	}
	if accessToken != nil {
		device.AccessToken = *accessToken
	}
	if refreshToken != nil {
		device.RefreshToken = *refreshToken
	}
	if tokenExpiresAt != nil {
		device.TokenExpiresAt = tokenExpiresAt
	}
	if isEnabled != nil {
		device.IsEnabled = *isEnabled
	}

	if err := uc.deviceRepo.UpdateDevice(device); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return device, nil
}

// SyncDevice manually triggers a device sync
func (uc *DeviceUsecase) SyncDevice(deviceID int) error {
	device, err := uc.deviceRepo.GetByID(deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	if device.Status != domain.DeviceStatusConnected {
		return ErrDeviceNotConnected
	}

	// Create sync log
	syncLog := &domain.SyncLog{
		DeviceID:      deviceID,
		Status:        "syncing",
		RecordsSynced: 0,
	}
	uc.deviceRepo.CreateSyncLog(syncLog)

	// Update device status to syncing
	uc.deviceRepo.UpdateDeviceStatus(deviceID, domain.DeviceStatusSyncing)

	// In a real implementation, this would call the device's API to fetch data
	// For now, we'll simulate a successful sync
	recordsSynced := 0

	// Simulate fetching and processing data
	if err := uc.processDeviceData(device); err != nil {
		// Update sync log with error
		now := time.Now()
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
		syncLog.CompletedAt = &now
		uc.deviceRepo.CreateSyncLog(syncLog)

		// Update device status to error
		uc.deviceRepo.UpdateDeviceStatus(deviceID, domain.DeviceStatusError)
		return fmt.Errorf("failed to sync device: %w", err)
	}

	// Update sync log with success
	now := time.Now()
	syncLog.Status = "success"
	syncLog.RecordsSynced = recordsSynced
	syncLog.CompletedAt = &now
	uc.deviceRepo.CreateSyncLog(syncLog)

	// Update device status and sync times
	nextSyncAt := time.Now().Add(time.Duration(device.SyncFrequency) * time.Second)
	uc.deviceRepo.UpdateSyncTime(deviceID, now, nextSyncAt)
	uc.deviceRepo.UpdateDeviceStatus(deviceID, domain.DeviceStatusConnected)

	return nil
}

// GetDeviceData retrieves data from a device
func (uc *DeviceUsecase) GetDeviceData(deviceID int, dataType domain.DataType, startDate, endDate time.Time) ([]domain.DeviceData, error) {
	data, err := uc.deviceRepo.GetDeviceData(deviceID, dataType, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get device data: %w", err)
	}
	return data, nil
}

// GetSyncLogs retrieves sync logs for a device
func (uc *DeviceUsecase) GetSyncLogs(deviceID int, limit int) ([]domain.SyncLog, error) {
	logs, err := uc.deviceRepo.GetSyncLogs(deviceID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get sync logs: %w", err)
	}
	return logs, nil
}

// processDeviceData processes data from a wearable device and stores it in the appropriate measurements
func (uc *DeviceUsecase) processDeviceData(device *domain.Device) error {
	// Get unprocessed data
	deviceData, err := uc.deviceRepo.GetUnprocessedData(device.ID)
	if err != nil {
		return fmt.Errorf("failed to get unprocessed data: %w", err)
	}

	// Process each data point
	for _, data := range deviceData {
		switch data.DataType {
		case domain.DataTypeSteps:
			// Would store in nutrition measurement as step count
			// This is a simplified implementation
		case domain.DataTypeHeartRate:
			// Would store in metabolic measurement
		case domain.DataTypeSleep:
			// Would store in wellbeing measurement
		case domain.DataTypeCalories:
			// Would store in nutrition measurement
		case domain.DataTypeActiveMinutes:
			// Would store in fitness measurement
		}

		// Mark as processed
		uc.deviceRepo.MarkDataAsProcessed(data.ID)
	}

	return nil
}

// EnableDevice enables a device for automatic syncing
func (uc *DeviceUsecase) EnableDevice(deviceID int) error {
	device, err := uc.deviceRepo.GetByID(deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	device.IsEnabled = true
	if err := uc.deviceRepo.UpdateDevice(device); err != nil {
		return fmt.Errorf("failed to enable device: %w", err)
	}

	return nil
}

// DisableDevice disables a device from automatic syncing
func (uc *DeviceUsecase) DisableDevice(deviceID int) error {
	device, err := uc.deviceRepo.GetByID(deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	device.IsEnabled = false
	if err := uc.deviceRepo.UpdateDevice(device); err != nil {
		return fmt.Errorf("failed to disable device: %w", err)
	}

	return nil
}

// SetSyncFrequency sets the sync frequency for a device
func (uc *DeviceUsecase) SetSyncFrequency(deviceID int, frequency int) error {
	device, err := uc.deviceRepo.GetByID(deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	device.SyncFrequency = frequency
	if err := uc.deviceRepo.UpdateDevice(device); err != nil {
		return fmt.Errorf("failed to set sync frequency: %w", err)
	}

	return nil
}