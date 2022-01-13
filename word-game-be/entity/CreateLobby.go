package entity

type CreateLobby struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}
