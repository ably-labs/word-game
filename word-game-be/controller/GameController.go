package controller

import (
	"github.com/ably-labs/word-game/word-game-be/constant"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably/ably-go/ably"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type GameController struct {
	db   *gorm.DB
	ably *ably.Realtime
}

func NewGameController(e *echo.Group, db *gorm.DB, ably *ably.Realtime) *GameController {
	bc := &GameController{
		db:   db,
		ably: ably,
	}

	g := e.Group("/:id", middleware.RequireAuthorisation, middleware.ValidateLobby, middleware.RequireLobbyMember)

	g.GET("/boards", bc.GetBoards)

	return bc
}

func (bc *GameController) GetBoards(c echo.Context) error {
	lobby := c.Get("lobby").(*model.Lobby)

	boards := make(map[string]entity.SquareSet)
	boards["main"] = lobby.Board

	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)
	if lobbyMember.MemberType != constant.MemberTypeSpectator {
		boards["deck"] = lobbyMember.Deck
	}

	return c.JSON(200, boards)
}

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

// newTileBag creates a new standard bag of tiles
func newTileBag() entity.SquareSet {
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

const boardLayout = "T--D---T---D--T" +
	"-W---L---L---W-" +
	"--W---D-D---W--" +
	"D--W---D---W--D" +
	"----W-----W----" +
	"-L---L---L---L-" +
	"--D---D-D---D--" +
	"T--D---S---D--T" +
	"--D---D-D---D--" +
	"-L---L---L---L-" +
	"----W-----W----" +
	"D--W---D---W--D" +
	"--W---D-D---W--" +
	"-W---L---L---W-" +
	"T--W---T---W--T"

var boardMap = map[rune]entity.Square{
	'-': {},
	'T': {Bonus: &entity.Bonus{Text: "TRIPLE WORD", Type: "triple-word"}},
	'D': {Bonus: &entity.Bonus{Text: "DOUBLE LETTER", Type: "double-letter"}},
	'W': {Bonus: &entity.Bonus{Text: "DOUBLE WORD", Type: "double-word"}},
	'L': {Bonus: &entity.Bonus{Text: "TRIPLE LETTER", Type: "triple-letter"}},
	'S': {Bonus: &entity.Bonus{Text: "START", Type: "double-word"}},
}

func newBoard() entity.SquareSet {
	board := make([]entity.Square, 15*15)

	for i, char := range boardLayout {
		board[i] = boardMap[char]
	}

	return entity.SquareSet{
		Squares: &board,
		Width:   15,
		Height:  15,
	}
}

func takeFromBag(n int, bag *entity.SquareSet) []entity.Square {
	tiles := (*bag.Squares)[:n]
	*bag.Squares = (*bag.Squares)[n+1:]
	for i := range tiles {
		tiles[i].Tile.Draggable = true
	}
	return tiles
}
