package layout

import (
	"github.com/ably-labs/word-game/word-game-be/entity"
)

type BoardLayoutStripe struct {
}

func (b BoardLayoutStripe) PlaceBonus(width int, height int, index int) *entity.Bonus {

	repeat := width / 2

	if index%repeat == 0 {
		return &entity.Bonus{LetterMultiplier: 3}
	}

	if index%repeat == 1 {
		return &entity.Bonus{WordMultiplier: 3}
	}

	if index%repeat == 2 {
		return &entity.Bonus{LetterMultiplier: 2}
	}

	if index%repeat == 3 {
		return &entity.Bonus{WordMultiplier: 2}
	}

	return nil
}
