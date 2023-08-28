package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserInterface interface {
	Create(db *gorm.DB) error
	Update(db *gorm.DB) error
	FindByEmail(db *gorm.DB, email string) (*User, error)
	Delete(db *gorm.DB) error
}

type User struct {
	gorm.Model
	ID                       uuid.UUID `gorm:"type:uuid;primary_key;"`
	UCardNumber              *string
	UID                      string  `gorm:"index"`
	Email                    *string `validate:"required,email" gorm:"index"`
	Forename                 string  `validate:"required,alpha"`
	Surname                  string  `validate:"required,alpha"`
	PreferredName            *string
	RepStatus                bool
	UserAgreementVersion     int32         `validate:"gte=0"`
	MailingListSubscriptions []MailingList `gorm:"many2many:user_mailinglists;"`
	Roles                    []Role        `gorm:"many2many:user_roles;"`
}

// BeforeCreate will set a UUID rather than a numeric ID.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewSHA1(uuid.NameSpaceURL, []byte(*u.Email))
	return
}

func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) Update(db *gorm.DB) error {
	return db.Save(u).Error
}

func (u *User) FindByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Delete(db *gorm.DB) error {
	return db.Delete(u).Error
}

func (u *User) HasPermission(permissionName string) bool {
	// Logic to check for the highest-level permission among all roles
	return true // Placeholder, you'll implement the actual logic
}
