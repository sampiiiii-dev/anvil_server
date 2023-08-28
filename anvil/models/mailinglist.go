package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MailingList struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name        string
	Description string
	Users       []User `gorm:"many2many:user_mailinglists;"`
}
