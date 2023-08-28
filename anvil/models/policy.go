package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Policy struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;"`
	Name        string    `gorm:"unique;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:policy_permissions;"`
}
