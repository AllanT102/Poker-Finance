package models

import (
	"github.com/google/uuid"
)

type PlayedGames struct {
	GameID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Game      Game      `gorm:"foreignKey:GameID;references:ID"`
	PlayerID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	Player    User      `gorm:"foreignKey:PlayerID;references:ID"`
	BuyIn     int       `gorm:"type:int;not null;default:0"`
	EndAmount int       `gorm:"type:int;not null;default:0"`
}

func NewPlayedGames(gameID uuid.UUID, playerID uuid.UUID, buyIn, endAmount int) *PlayedGames {
	return &PlayedGames{
		GameID:    gameID,
		PlayerID:  playerID,
		BuyIn:     buyIn,
		EndAmount: endAmount,
	}
}
