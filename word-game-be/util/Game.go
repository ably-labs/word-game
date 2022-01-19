package util

import (
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"math"
	"math/rand"
	"strings"
	"time"
)

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
			upperLetter := strings.ToUpper(letter)
			board[i] = entity.Square{Tile: &entity.Tile{
				Letter:    upperLetter,
				Score:     score[upperLetter],
				Draggable: letter != upperLetter,
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

func ValidateBoard(squareSet entity.SquareSet) int {

	newWords := GetNewWords(squareSet)

	if len(newWords) == 0 {
		// Invalid placement or no tiles placed
		return 0
	}

	totalScore := 0
	for _, word := range newWords {
		constructedWord := ""
		score := 0
		multiplier := 1
		for _, square := range word {
			constructedWord += square.Tile.Letter
			score += GetSquareScore(*square)
			multiplier += GetSquareWordMultiplier(*square)
		}
		if !IsValidWord(constructedWord) {
			fmt.Println("Invalid word: ", constructedWord)
			return 0
		}
		fmt.Println(constructedWord, score, multiplier)
		totalScore += score * multiplier
	}
	return totalScore
}

func GetNewWords(squareSet entity.SquareSet) [][]*entity.Square {
	indices := GetPlacedTileIndices(squareSet)

	fmt.Println("ind 0 row start", GetRowStart(squareSet, indices[0]))
	fmt.Println("ind 1 row start", GetRowStart(squareSet, indices[1]))
	// If the first and second tiles are on the same row, this must be a horizontal word
	isHoz := GetRowStart(squareSet, indices[0]) == GetRowStart(squareSet, indices[1])

	fmt.Println("isHoz", isHoz)

	originalWord := GetSquaresForWord(squareSet, indices[0], isHoz)

	// Check every single draggable tile is inside the original word
	seenCount := 0
	for i, square := range originalWord {
		if square.Tile == nil {
			fmt.Printf("WARN: Word starting at %d contains invalid square %v at position %d\n", indices[0], square, i)
			continue
		}
		if square.Tile.Draggable {
			seenCount++
		}
	}

	// If there are less draggable tiles inside the original word, the placement is invalid
	if seenCount < len(indices) {
		return [][]*entity.Square{}
	}

	words := [][]*entity.Square{
		originalWord,
	}

	// Collect the new word boundaries for each row
	for _, index := range indices {
		wordSquares := GetSquaresForWord(squareSet, index, !isHoz)
		if len(wordSquares) > 1 {
			words = append(words, wordSquares)
		}
	}

	return words
}

func GetSquaresForWord(squareSet entity.SquareSet, index int, isHoz bool) []*entity.Square {
	start, end, interval := 0, 0, 0
	if isHoz {
		start, end = GetWordBoundsHoz(squareSet, index)
		interval = 1
	} else {
		start, end = GetWordBoundsVert(squareSet, index)
		interval = squareSet.Width
	}
	wordSquares := make([]*entity.Square, 0)

	for i := start; i <= end; i += interval {
		wordSquares = append(wordSquares, &(*squareSet.Squares)[i])
	}
	return wordSquares
}

func GetWordBoundsHoz(squareSet entity.SquareSet, target int) (int, int) {
	squares := *squareSet.Squares

	// Get the start and ends of this row
	rowStart := GetRowStart(squareSet, target)
	rowEnd := GetRowEnd(squareSet, target)

	start := rowStart
	end := rowEnd

	// walk backwards through the board until we run out of placed tiles on that row
	for i := target; i > rowStart; i-- {
		if squares[i].Tile == nil {
			start = i + 1
			break
		}
	}

	// walk forwards in the same manner
	for i := target; i <= rowEnd; i++ {
		if squares[i].Tile == nil {
			end = i - 1
			break
		}
	}

	return start, end
}

func GetWordBoundsVert(squareSet entity.SquareSet, target int) (int, int) {
	squares := *squareSet.Squares

	colStart := GetColStart(squareSet, target)
	colEnd := GetColEnd(squareSet, target)

	start := colStart
	end := colEnd

	for i := target; i > colStart; i -= squareSet.Width {
		fmt.Println("WBV Walk back ", i, squares[i].Tile)
		if squares[i].Tile == nil {
			start = i + squareSet.Width
			break
		}
	}

	for i := target; i <= colEnd; i += squareSet.Width {
		fmt.Println("WBV Walk forward ", i, squares[i].Tile)
		if squares[i].Tile == nil {
			end = i - squareSet.Width
			break
		}
	}

	return start, end
}

func GetSquareScore(square entity.Square) int {
	// There is no letter on this tile, so there is no score
	if square.Tile == nil {
		return 0
	}

	// No bonus tile, so the score is the raw tile score
	if square.Bonus == nil {
		return square.Tile.Score
	}

	if square.Bonus.Type == "double-letter" {
		return square.Tile.Score * 2
	}

	if square.Bonus.Type == "triple-letter" {
		return square.Tile.Score * 3
	}

	return square.Tile.Score
}

func GetSquareWordMultiplier(square entity.Square) int {
	if square.Tile == nil || square.Bonus == nil {
		return 0
	}

	if square.Bonus.Type == "double-word" {
		return 2
	}

	if square.Bonus.Type == "triple-word" {
		return 3
	}

	return 0
}

// GetColStart gets the start index of a column based on the width of the entity.SquareSet
func GetColStart(squareSet entity.SquareSet, index int) int {
	rowNum := int(math.Floor(float64(index / squareSet.Width)))
	return index - (squareSet.Width * rowNum)
}

// GetColEnd gets the end index of a column based on the width of the entity.SquareSet
func GetColEnd(squareSet entity.SquareSet, index int) int {
	rowNum := int(math.Floor(float64(index / squareSet.Width)))
	return index + squareSet.Width*(squareSet.Height-rowNum-1)
}

// GetRowStart gets the start index of a row based on the width of the entity.SquareSet
func GetRowStart(squareSet entity.SquareSet, index int) int {
	return index - (index % squareSet.Width)
}

// GetRowEnd gets the end index of a row based on the width of the entity.SquareSet
func GetRowEnd(squareSet entity.SquareSet, index int) int {
	return index + (squareSet.Width - 1 - (index % squareSet.Width))
}
