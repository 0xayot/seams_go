package models

import (
	"gorm.io/datatypes"
)

type Measurement struct {
	Base
	UserID       uint   `gorm:"not null"` // Foreign key reference to User
	Name         string `gorm:"not null"`
	MeasuredBy   *string
	Measurements datatypes.JSON `gorm:"type:jsonb"` // PostgreSQL JSONB type
	ShoeSize     *string
	Active       *bool
}
