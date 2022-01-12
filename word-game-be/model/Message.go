package model

type Message struct {
	ID       *uint32 `gorm:"primarykey" json:"id"`
	AuthorID *uint32
}
