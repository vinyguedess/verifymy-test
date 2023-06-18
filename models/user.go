package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID   `json:"id" gorm:"primarykey;type:varchar(36)"`
	Name        string      `json:"name" gorm:"type:varchar(255)"`
	DateOfBirth Date        `json:"date_of_birth" gorm:"type:date"`
	Email       string      `json:"email" gorm:"type:varchar(255)"`
	Password    SecretValue `json:"password" gorm:"type:varchar(255)"`
	Address     string      `json:"address" gorm:"type:varchar(255)"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New()

	return nil
}
