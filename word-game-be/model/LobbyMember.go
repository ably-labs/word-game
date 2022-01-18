package model

import (
	"github.com/ably-labs/word-game/word-game-be/constant"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"time"
)

type LobbyMember struct {
	UserID     uint32              `gorm:"primaryKey;autoIncrement:false" json:"id"`
	LobbyID    uint32              `gorm:"primaryKey;autoIncrement:false" json:"lobbyId"`
	MemberType constant.MemberType `json:"type"`
	User       DisplayUser         `json:"user"`
	Deck       entity.SquareSet    `json:"-"`
	JoinedAt   time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"joined"`
	Score      int                 `json:"score" gorm:"default:0"`
}
