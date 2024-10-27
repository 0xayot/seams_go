package models

import (
	uuid "github.com/satori/go.uuid"
)

type Wish struct {
	Base
	UserID      uuid.UUID `gorm:"not null"` // Foreign key reference to User
	Name        string    `gorm:"not null"`
	Description *string
	Price       *string
	Url         *string
	Image       *string
}
