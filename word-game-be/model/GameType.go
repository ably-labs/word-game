package model

type GameType struct {
	ID   uint32 `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}
