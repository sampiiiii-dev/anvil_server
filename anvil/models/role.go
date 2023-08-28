package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;"`
	Name        string    `gorm:"unique;not null"`
	Description string
	Policies    []Policy `gorm:"many2many:role_policies;"`
}
