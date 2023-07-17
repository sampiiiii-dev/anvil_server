package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SignInRecord struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key;"`
	IsIn bool
}
