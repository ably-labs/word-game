package model

import "time"

type Message struct {
	ID        *uint32 `gorm:"primarykey" json:"id"`
	AuthorID  *uint32
	LobbyID   *int64
	Message   string
	Timestamp time.Time
	Author    User
}
