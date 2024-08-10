package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (base *Base) BeforeCreate(scope *gorm.DB) error {
	// scope.SetColumn("ID", uuid.NewV4().String())
	base.ID = uuid.NewV4()
	return nil
}
