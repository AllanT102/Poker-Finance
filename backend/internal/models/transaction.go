package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null"`
	User             User           `gorm:"foreignKey:UserID;references:ID"`
	PaymentDetailsID uuid.UUID      `gorm:"type:uuid;not null"`
	PaymentDetails   PaymentDetails `gorm:"foreignKey:PaymentDetailsID;references:ID"`
	CreatedAt        time.Time      `gorm:"type:timestamptz;default:current_timestamp"`
	Status           string         `gorm:"type:varchar(50);not null"`
}

func NewTransaction(userID, paymentDetailsID uuid.UUID, status string) *Transaction {
	return &Transaction{
		ID:               uuid.New(),
		UserID:           userID,
		PaymentDetailsID: paymentDetailsID,
		CreatedAt:        time.Now(),
		Status:           status,
	}
}
