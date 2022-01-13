package entity

type SendChat struct {
	Message string `json:"message"`
}

type ChatSent struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}
