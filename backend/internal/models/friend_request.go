package models

import (
    "time"
    "github.com/google/uuid"
)

type FriendRequest struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
    FriendID  uuid.UUID `json:"friend_id" gorm:"type:uuid;not null"`
    Status    string    `json:"status" gorm:"type:varchar(100);not null"` // e.g., "pending", "accepted", "declined"
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewFriends(userID, friendID uuid.UUID, status string) *FriendRequest {
    return &FriendRequest{
        ID:        uuid.New(),
        UserID:    userID,
        FriendID:  friendID,
        Status:    status,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}
