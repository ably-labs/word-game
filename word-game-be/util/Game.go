package util

import (
	_ "embed"
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"math"
	"math/rand"
	"strings"
	"time"
)

//go:embed words.txt
var wordList string
var words = strings.Split(wordList, "\n")

// TODO: All of these should be configurable, but in the interests of scope I'm hardcoding them

var distributions = map[string]int{
	"": 2, "A": 9, "B": 2, "C": 2, "D": 4,
	"E": 12, "F": 2, "G": 3, "H": 2, "I": 9,
	"J": 1, "K": 1, "L": 4, "M": 2, "N": 6,
	"O": 8, "P": 2, "Q": 1, "R": 6, "S": 4,
	"T": 6, "U": 4, "V": 2, "W": 2, "X": 1,
	"Y": 2, "Z": 1,
}

var score = map[string]int{
	"":  0,
	"A": 1, "E": 1, "I": 1, "O": 1, "U": 1, "L": 1, "N": 1, "S": 1, "T": 1, "R": 1,
	"D": 2, "G": 2,
	"B": 3, "C": 3, "M": 3, "P": 3,
	"F": 4, "H": 4, "V": 4, "W": 4, "Y": 4,
	"K": 5,
	"J": 8, "X": 8,
	"Q": 10, "Z": 10,
}

// NewTileBag creates a new standard bag of tiles
func NewTileBag() entity.SquareSet {
	bag := make([]entity.Square, 100)
	cursor := 0
	for letter, amt := range distributions {
		for i := cursor; i < cursor+amt; i++ {
			bag[i] = entity.Square{Tile: &entity.Tile{
				Letter: letter,
				Score:  score[letter],
			}}
		}
		cursor += amt
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(bag), func(i, j int) {
		bag[i], bag[j] = bag[j], bag[i]
	})
	return entity.SquareSet{Squares: &bag}
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
	'~': {Bonus: &entity.Bonus{Text: "TRIPLE WORD", Type: "triple-word"}},
	'@': {Bonus: &entity.Bonus{Text: "DOUBLE LETTER", Type: "double-letter"}},
	'#': {Bonus: &entity.Bonus{Text: "DOUBLE WORD", Type: "double-word"}},
	'!': {Bonus: &entity.Bonus{Text: "TRIPLE LETTER", Type: "triple-letter"}},
	'*': {Bonus: &entity.Bonus{Text: "START", Type: "double-word"}},
}

func NewBoardFromLayout(layout string, width int, height int) entity.SquareSet {
	board := make([]entity.Square, width*height)

	for i, char := range layout {
		square, ok := boardMap[char]
		if ok {
			board[i] = square
		} else {
			letter := string(char)
			board[i] = entity.Square{Tile: &entity.Tile{
				Letter:    letter,
				Score:     score[letter],
				Draggable: false,
			}}
		}
	}

	return entity.SquareSet{
		Squares: &board,
		Width:   width,
		Height:  height,
	}
}

func NewBoard() entity.SquareSet {
	return NewBoardFromLayout(boardLayout, 15, 15)
}

func TakeFromBag(n int, bag *entity.SquareSet) []entity.Square {
	tiles := (*bag.Squares)[:n]
	*bag.Squares = (*bag.Squares)[n+1:]
	for i := range tiles {
		tiles[i].Tile.Draggable = true
	}
	return tiles
}

// GetPlacedTileIndices gets all Draggable tiles in an entity.SquareSet
func GetPlacedTileIndices(squareSet entity.SquareSet) []int {
	indices := make([]int, 0)
	for i, sq := range *squareSet.Squares {
		if sq.Tile != nil && sq.Tile.Draggable {
			indices = append(indices, i)
		}
	}
	return indices
}

// ValidatePlacement validates that the tiles are placed in a valid way (regardless of if it's a valid word)
//func ValidatePlacement(squareSet entity.SquareSet) bool {
//	squares := *squareSet.Squares
//	indices := GetPlacedTileIndices(squareSet)
//	seenCount := 0
//	startIndex := 0
//
//	for i := indices[0]; i > 0; i++ {
//		if squares[i].Tile == nil {
//			startIndex = i + 1
//			break
//		}
//	}
//
//	return true
//}

func GetWordBoundsHoz(squareSet entity.SquareSet, target int) (int, int) {
	squares := *squareSet.Squares

	// Get the start and ends of this row
	rowStart := target - (target % squareSet.Width)
	rowEnd := target + (squareSet.Width - 1 - (target % squareSet.Width))

	start := rowStart
	end := rowEnd

	// walk backwards through the board until we run out of placed tiles on that row
	for i := target; i > rowStart; i-- {
		fmt.Println("walkback", i, squares[i].Tile)
		if squares[i].Tile == nil {
			start = i + 1
			break
		}
	}

	// walk forwards in the same manner
	for i := target; i <= rowEnd; i++ {
		fmt.Println("walk forward", i, squares[i].Tile)
		if squares[i].Tile == nil {
			end = i - 1
			break
		}
	}

	return start, end
}

func GetWordBoundsVert(squareSet entity.SquareSet, target int) (int, int) {
	fmt.Println("Target", target)
	squares := *squareSet.Squares

	// Get the start and ends of this row
	rowNum := int(math.Floor(float64(target / squareSet.Width)))
	fmt.Println("rowNum", rowNum)
	colStart := target - (squareSet.Width * rowNum)
	colEnd := target + squareSet.Width*(squareSet.Height-rowNum-1)

	fmt.Println("colStart colEnd", colStart, colEnd)

	start := colStart
	end := colEnd

	for i := target; i > colStart; i -= squareSet.Width {
		fmt.Println("walkback", i, squares[i].Tile)
		if squares[i].Tile == nil {
			start = i + squareSet.Width
			break
		}
	}

	fmt.Println("Finished walkback with ", start)

	for i := target; i <= colEnd; i += squareSet.Width {
		fmt.Println("walk forward", i, squares[i].Tile)
		if squares[i].Tile == nil {
			end = i - squareSet.Width
			break
		}
	}

	fmt.Println("Finished walk forward with ", end)

	return start, end
}
