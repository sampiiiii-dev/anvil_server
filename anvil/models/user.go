package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                       uuid.UUID `gorm:"type:uuid;primary_key;"`
	UCardNumber              *string
	UID                      string
	Email                    *string
	Forename                 string
	Surname                  string
	PreferredName            *string
	RepStatus                bool
	UserAgreementVersion     int32
	MailingListSubscriptions []MailingList `gorm:"many2many:user_mailinglists;"`

	// Time
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

/**
 * TableName
 *
 * @return string
 */
func (User) TableName() string {
	return "users"
}

/**
 * isRep
 *
 * @return bool
 */
func (u *User) IsRep() bool {
	return u.RepStatus
}
