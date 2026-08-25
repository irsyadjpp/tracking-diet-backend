package domain

import "time"

type LabTest struct {
	ID             int       `gorm:"primaryKey;autoIncrement;column:lab_id" json:"id"`
	UserID         int       `gorm:"not null;index;column:user_id" json:"user_id"`
	TestName       string    `gorm:"size:100;not null;column:test_name" json:"test_name"`
	ResultValue    *string   `gorm:"column:result_value" json:"result_value"`
	Unit           *string   `gorm:"size:20;column:unit" json:"unit"`
	ReferenceRange *string  `gorm:"size:50;column:reference_range" json:"reference_range"`
	MeasuredAt     time.Time `gorm:"not null;column:measured_at" json:"measured_at"`
}

// TableName specifies the table name for LabTest model
func (LabTest) TableName() string {
	return "lab_tests"
}

type LabTestRepository interface {
	Create(lt *LabTest) error
	GetByID(id int) (*LabTest, error)
	ListByUser(userID int) ([]LabTest, error)
	ListByUserWithTestName(userID int, testName string) ([]LabTest, error)
	ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]LabTest, error)
	Update(lt *LabTest) error
	Delete(id int) error
}
