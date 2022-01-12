package entity

import (
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably/ably-go/ably"
)

type AuthSuccess struct {
	User         model.User        `json:"user"`
	TokenRequest ably.TokenRequest `json:"token"`
}
