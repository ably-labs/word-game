package model

import (
	"github.com/ably-labs/word-game/word-game-be/entity"
	"time"
)

type Lobby struct {
	ID *int64 `gorm:"primarykey" json:"-"`
	// Poor JS gets upset with numbers this big, so we need to pass it as a string rather than an int64
	IdStr          string            `gorm:"-" json:"id"`
	Name           string            `json:"name"`
	CreatorID      *uint32           `json:"creatorId"`
	CreatedAt      time.Time         `json:"createdAt"`
	State          entity.LobbyState `json:"state"`
	Private        bool              `json:"private"`
	Joinable       bool              `json:"joinable"`
	CurrentPlayers uint8             `json:"currentPlayers"`
	MaxPlayers     uint8             `json:"maxPlayers"`
	GameTypeID     uint32            `json:"-"`
	PlayerTurnID   *uint32           `json:"playerTurnId"`
	GameType       GameType          `json:"gameType"`
	Members        []LobbyMember     `json:"-"`
	Messages       []Message         `json:"-"`
	Creator        DisplayUser       `json:"creator" gorm:"foreignKey:id;references:creator_id"`
	Board          entity.SquareSet  `json:"-"`
	Bag            entity.SquareSet  `json:"-"`
}
