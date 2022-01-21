package entity

type CreateLobby struct {
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	GameType uint32 `json:"gameType"`
}
