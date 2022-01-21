package layout

import (
	"github.com/ably-labs/word-game/word-game-be/entity"
	"math"
)

type BoardLayoutRegular struct {
}

func (b BoardLayoutRegular) PlaceBonus(width int, height int, index int) *entity.Bonus {

	// Corners and border center has triple words
	colCenter := int(math.Ceil(float64(height/2))) * width
	rowCenter := int(math.Ceil(float64(width / 2)))
	if index == width*height || index == colCenter-1 || index == colCenter-width || index == rowCenter || index == width*height-rowCenter-1 {
		return &entity.Bonus{
			WordMultiplier: 3,
		}
	}

	posX := index % width
	posY := index / height

	if posX == posY || index%(width-1) == 0 {
		return &entity.Bonus{
			WordMultiplier: 2,
		}
	}

	return nil
}
