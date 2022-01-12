package model

import "github.com/ably-labs/word-game/word-game-be/constant"

type LobbyMember struct {
	UserID     uint32              `gorm:"primaryKey;autoIncrement:false"`
	LobbyID    uint32              `gorm:"primaryKey;autoIncrement:false"`
	MemberType constant.MemberType `json:"type"`
}
