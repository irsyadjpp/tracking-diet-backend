package domain

import "time"

type LabTest struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null;index" json:"user_id"`
	TestName    string    `gorm:"size:100;not null" json:"test_name"`
	ResultValue float64   `json:"result_value"`
	Unit        string    `gorm:"size:20" json:"unit"`
	NormalRange string    `gorm:"size:50" json:"normal_range"`
	Notes       string    `gorm:"type:text" json:"notes"`
	MeasuredAt  time.Time `gorm:"not null" json:"measured_at"`
}

type LabTestRepository interface {
	Create(lt *LabTest) error
	GetByID(id int) (*LabTest, error)
	ListByUser(userID int) ([]LabTest, error)
	Update(lt *LabTest) error
	Delete(id int) error
}
