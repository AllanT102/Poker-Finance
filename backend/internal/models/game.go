package models

import (
    "time"
    "github.com/google/uuid"
)

type Game struct {
    ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    Date time.Time `gorm:"type:timestamptz;not null"`
}

func NewGame(date time.Time) *Game {
    return &Game{
        ID:   uuid.New(),
        Date: date,
    }
}
