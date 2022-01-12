package model

import "time"

type LobbyState string

const (
	LobbyWaitingForPlayers LobbyState = "waiting"
	LobbyInGame            LobbyState = "inGame"
	LobbyRoundOver         LobbyState = "roundOver"
)

type Lobby struct {
	ID        *uint32    `gorm:"primarykey" json:"id"`
	Name      string     `json:"name"`
	CreatorID *uint32    `json:"creatorId"`
	CreatedAt time.Time  `json:"createdAt"`
	State     LobbyState `json:"state"`
	Private   bool       `json:"private"`
}
