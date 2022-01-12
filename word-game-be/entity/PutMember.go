package entity

import "github.com/ably-labs/word-game/word-game-be/constant"

type PutMember struct {
	// Join codes, for later
	Code string              `json:"code"`
	Type constant.MemberType `json:"type"`
}
