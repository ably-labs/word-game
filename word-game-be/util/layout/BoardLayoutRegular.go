package layout

import (
	"github.com/ably-labs/word-game/word-game-be/entity"
)

type BoardLayoutRegular struct {
}

const boardLayout = "~--@---~---@--~" +
	"-#---!---!---#-" +
	"--#---@-@---#--" +
	"@--#---@---#--@" +
	"----#-----#----" +
	"-!---!---!---!-" +
	"--@---@-@---@--" +
	"~--@---*---@--~" +
	"--@---@-@---@--" +
	"-!---!---!---!-" +
	"----#-----#----" +
	"@--#---@---#--@" +
	"--#---@-@---#--" +
	"-#---!---!---#-" +
	"~--#---~---#--~"

var boardMap = map[rune]entity.Square{
	'-': {},
	'~': {Bonus: &entity.Bonus{WordMultiplier: 3}},
	'@': {Bonus: &entity.Bonus{LetterMultiplier: 2}},
	'#': {Bonus: &entity.Bonus{WordMultiplier: 2}},
	'!': {Bonus: &entity.Bonus{LetterMultiplier: 3}},
	'*': {Bonus: &entity.Bonus{WordMultiplier: 2, Start: true}},
}

func (b BoardLayoutRegular) PlaceBonus(width int, height int, index int) *entity.Bonus {
	// TODO, this but maths
	return boardMap[rune(boardLayout[index])].Bonus
}
