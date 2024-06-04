package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email       string    `gorm:"type:varchar(255);unique_index;not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	DisplayName string    `gorm:"type:varchar(255)"`
	Balance     float64   `gorm:"type:float;not null"`
	CreatedAt   time.Time `gorm:"type:timestamptz;default:current_timestamp"`
}

func NewUser(email, name, displayName string) *User {
	return &User{
		ID:          uuid.New(), // generateUUID() needs to be defined or imported
		Email:       email,
		Name:        name,
		DisplayName: displayName,
		Balance:     0,
	}
}
