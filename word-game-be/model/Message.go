package model

import "time"

type Message struct {
	ID        *uint32 `gorm:"primarykey" json:"id"`
	AuthorID  *uint32
	LobbyID   *uint32
	Message   string
	Timestamp time.Time
	Author    User
}
