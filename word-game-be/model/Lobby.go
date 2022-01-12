package model

import "time"

type LobbyState string

const (
	LobbyWaitingForPlayers LobbyState = "waiting"
	LobbyInGame            LobbyState = "inGame"
	LobbyRoundOver         LobbyState = "roundOver"
)

type Lobby struct {
	ID             *uint32    `gorm:"primarykey" json:"id"`
	Name           string     `json:"name"`
	CreatorID      *uint32    `json:"creatorId"`
	CreatedAt      time.Time  `json:"createdAt"`
	State          LobbyState `json:"state"`
	Private        bool       `json:"private"`
	Joinable       bool       `json:"joinable"`
	CurrentPlayers uint8      `json:"currentPlayers"`
	MaxPlayers     uint8      `json:"maxPlayers"`
	GameTypeID     uint32     `json:"-"`
	GameType       GameType   `json:"gameType"`
}
