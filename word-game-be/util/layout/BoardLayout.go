package layout

import "github.com/ably-labs/word-game/word-game-be/entity"

type BoardLayout interface {
	PlaceBonus(width int, height int, index int) *entity.Bonus
}
