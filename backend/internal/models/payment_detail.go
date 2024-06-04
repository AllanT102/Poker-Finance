package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentDetails struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GameID        uuid.UUID `gorm:"type:uuid;not null"`
	Game          Game      `gorm:"foreignKey:GameID;references:ID"`
	PayerID       uuid.UUID `gorm:"type:uuid;not null"`
	Payer         User      `gorm:"foreignKey:PayerID;references:ID"`
	PayeeID       uuid.UUID `gorm:"type:uuid;not null"`
	Payee         User      `gorm:"foreignKey:PayeeID;references:ID"`
	Amount        float64   `gorm:"type:float;not null"`
	Details       string    `gorm:"type:text"`
	TimeSubmitted time.Time `gorm:"type:timestamptz;default:current_timestamp"`
	TimeCompleted time.Time `gorm:"type:timestamptz"`
	Status        string    `gorm:"type:varchar(50);not null"` // emailSent, complete, pending
}

func NewPaymentDetails(gameID uuid.UUID, payerID uuid.UUID, payeeID uuid.UUID, amount float64, details, status string) *PaymentDetails {
	return &PaymentDetails{
		ID:            uuid.New(),
		GameID:        gameID,
		PayerID:       payerID,
		PayeeID:       payeeID,
		Amount:        amount,
		Details:       details,
		TimeSubmitted: time.Now(),
		Status:        status,
	}
}
