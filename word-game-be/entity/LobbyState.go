package entity

type LobbyState string

const (
	LobbyWaitingForPlayers LobbyState = "waiting"
	LobbyInGame            LobbyState = "inGame"
	LobbyRoundOver         LobbyState = "roundOver"
)
