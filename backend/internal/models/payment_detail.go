package models

import (
    "time"
    "github.com/google/uuid"
)

type PaymentDetails struct {
    ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    PayerID       uuid.UUID `gorm:"type:uuid;not null"`
    Payer         User      `gorm:"foreignKey:PayerID;references:ID"`
    PayeeID       uuid.UUID `gorm:"type:uuid;not null"`
    Payee         User      `gorm:"foreignKey:PayeeID;references:ID"`
    Amount        int       `gorm:"type:int;not null"`
    Details       string    `gorm:"type:text"`
    TimeSubmitted time.Time `gorm:"type:timestamptz;default:current_timestamp"`
    TimeCompleted time.Time `gorm:"type:timestamptz"`
    Status        string    `gorm:"type:varchar(50);not null"`
}

func NewPaymentDetails(payerID, payeeID uuid.UUID, amount int, details, status string) *PaymentDetails {
    return &PaymentDetails{
        ID:            uuid.New(),
        PayerID:       payerID,
        PayeeID:       payeeID,
        Amount:        amount,
        Details:       details,
        TimeSubmitted: time.Now(),
        Status:        status,
    }
}
