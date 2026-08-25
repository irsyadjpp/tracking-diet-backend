package domain

import "time"

// DeviceType represents different types of wearable devices
type DeviceType string

const (
	DeviceTypeFitbit      DeviceType = "fitbit"
	DeviceTypeAppleWatch  DeviceType = "apple_watch"
	DeviceTypeGarmin      DeviceType = "garmin"
	DeviceTypeSamsung     DeviceType = "samsung"
	DeviceTypeWithings    DeviceType = "withings"
	DeviceTypePolar       DeviceType = "polar"
	DeviceTypeOura        DeviceType = "oura"
	DeviceTypeWhoop       DeviceType = "whoop"
)

// DeviceStatus represents the status of a device connection
type DeviceStatus string

const (
	DeviceStatusConnected    DeviceStatus = "connected"
	DeviceStatusDisconnected DeviceStatus = "disconnected"
	DeviceStatusSyncing      DeviceStatus = "syncing"
	DeviceStatusError        DeviceStatus = "error"
)

// Device represents a user's connected wearable device
type Device struct {
	ID              int         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int         `gorm:"not null;index" json:"user_id"`
	Type            DeviceType  `gorm:"size:50;not null" json:"type"`
	Name            string      `gorm:"size:100" json:"name"`
	SerialNumber    string      `gorm:"size:100" json:"serial_number"`
	Model           string      `gorm:"size:100" json:"model"`
	FirmwareVersion string      `gorm:"size:50" json:"firmware_version"`
	AccessToken     string      `gorm:"size:500" json:"access_token"`
	RefreshToken    string      `gorm:"size:500" json:"refresh_token"`
	TokenExpiresAt  *time.Time  `json:"token_expires_at"`
	Status          DeviceStatus `gorm:"size:20;default:disconnected" json:"status"`
	LastSyncAt      *time.Time  `json:"last_sync_at"`
	NextSyncAt      *time.Time  `json:"next_sync_at"`
	SyncFrequency   int         `gorm:"default:3600" json:"sync_frequency"` // in seconds
	IsEnabled       bool        `gorm:"default:true" json:"is_enabled"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Device model
func (Device) TableName() string {
	return "devices"
}

// DataType represents different types of data that can be synced
type DataType string

const (
	DataTypeSteps         DataType = "steps"
	DataTypeHeartRate     DataType = "heart_rate"
	DataTypeSleep         DataType = "sleep"
	DataTypeCalories      DataType = "calories"
	DataTypeDistance      DataType = "distance"
	DataTypeActiveMinutes DataType = "active_minutes"
	DataTypeFloors        DataType = "floors"
	DataTypeSpO2          DataType = "sp_o2"
	DataTypeStress        DataType = "stress"
	DataTypeBodyBattery   DataType = "body_battery"
)

// DeviceData represents synced data from a wearable device
type DeviceData struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID    int       `gorm:"not null;index" json:"device_id"`
	DataType    DataType  `gorm:"size:50;not null;index" json:"data_type"`
	Value       float64   `gorm:"not null" json:"value"`
	Unit        string    `gorm:"size:20" json:"unit"`
	RecordedAt  time.Time `gorm:"not null;index" json:"recorded_at"`
	SyncedAt    time.Time `gorm:"autoCreateTime" json:"synced_at"`
	IsProcessed bool      `gorm:"default:false" json:"is_processed"`
}

// TableName specifies the table name for DeviceData model
func (DeviceData) TableName() string {
	return "device_data"
}

// SyncLog represents a log of device sync operations
type SyncLog struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID    int       `gorm:"not null;index" json:"device_id"`
	Status      string    `gorm:"size:20;not null" json:"status"` // success, failed, partial
	RecordsSynced int     `gorm:"default:0" json:"records_synced"`
	ErrorMessage string   `gorm:"type:text" json:"error_message"`
	StartedAt   time.Time `gorm:"autoCreateTime" json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// TableName specifies the table name for SyncLog model
func (SyncLog) TableName() string {
	return "sync_logs"
}

// DeviceRepository defines operations for managing wearable devices
type DeviceRepository interface {
	// Device management
	CreateDevice(device *Device) error
	GetByID(id int) (*Device, error)
	GetByUser(userID int) ([]Device, error)
	GetByUserAndType(userID int, deviceType DeviceType) (*Device, error)
	UpdateDevice(device *Device) error
	DeleteDevice(id int) error
	UpdateDeviceStatus(id int, status DeviceStatus) error
	UpdateSyncTime(id int, lastSyncAt, nextSyncAt time.Time) error

	// Data management
	StoreDeviceData(data *DeviceData) error
	GetDeviceData(deviceID int, dataType DataType, startDate, endDate time.Time) ([]DeviceData, error)
	GetUnprocessedData(deviceID int) ([]DeviceData, error)
	MarkDataAsProcessed(id int) error
	DeleteOldData(olderThan time.Time) error

	// Sync logs
	CreateSyncLog(log *SyncLog) error
	GetSyncLogs(deviceID int, limit int) ([]SyncLog, error)
	GetLatestSyncLog(deviceID int) (*SyncLog, error)
}