package models

import (
	"github.com/google/uuid"
)

type PlayedGames struct {
	GameID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Game      Game      `gorm:"foreignKey:GameID;references:ID"`
	PlayerID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	Player    User      `gorm:"foreignKey:PlayerID;references:ID"`
	BuyIn     float64   `gorm:"type:float;not null;default:0"`
	EndAmount float64   `gorm:"type:float;not null;default:0"`
}

func NewPlayedGames(gameID uuid.UUID, playerID uuid.UUID, buyIn float64, endAmount float64) *PlayedGames {
	return &PlayedGames{
		GameID:    gameID,
		PlayerID:  playerID,
		BuyIn:     buyIn,
		EndAmount: endAmount,
	}
}
