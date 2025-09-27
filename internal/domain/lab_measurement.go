package domain

import "time"

type LabMeasurement struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Result    string    `json:"result"`
	Unit      string    `json:"unit"`
	RefRange  string    `json:"ref_range"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LabMeasurementRepository interface {
	CreateLabMeasurement(lab *LabMeasurement) error
	GetLabMeasurementByID(id int64) (*LabMeasurement, error)
	GetLabMeasurementsByUserID(userID int64) ([]*LabMeasurement, error)
	UpdateLabMeasurement(lab *LabMeasurement) error
	DeleteLabMeasurement(id int64) error
}
