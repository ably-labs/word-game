package entity

type UpdateLobby struct {
	State LobbyState `json:"state,omitempty"`
}
