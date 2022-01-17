package controller

import (
	_ "embed"
	"github.com/ably-labs/word-game/word-game-be/constant"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably/ably-go/ably"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GameController struct {
	db   *gorm.DB
	ably *ably.Realtime
}

func NewGameController(e *echo.Group, db *gorm.DB, ably *ably.Realtime) *GameController {
	gc := &GameController{
		db:   db,
		ably: ably,
	}

	g := e.Group("/:id", middleware.RequireAuthorisation, middleware.ValidateLobby, middleware.RequireLobbyMember)

	g.GET("/boards", gc.GetBoards)
	g.PATCH("/boards", gc.PatchBoard)
	g.POST("/boards", gc.EndTurn, middleware.RequireTurn)

	return gc
}

func (gc *GameController) GetBoards(c echo.Context) error {
	lobby := c.Get("lobby").(*model.Lobby)

	boards := make(map[string]entity.SquareSet)
	boards["main"] = lobby.Board

	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)
	if lobbyMember.MemberType != constant.MemberTypeSpectator {
		boards["deck"] = lobbyMember.Deck
	}

	return c.JSON(200, boards)
}

func (gc *GameController) PatchBoard(c echo.Context) error {
	moveTile := entity.MoveTile{}
	err := c.Bind(&moveTile)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobby := c.Get("lobby").(*model.Lobby)
	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)

	var fromBoard entity.SquareSet
	var toBoard entity.SquareSet

	boardUpdate := false

	if moveTile.From == "main" {
		fromBoard = lobby.Board
		boardUpdate = true
	} else {
		fromBoard = lobbyMember.Deck
	}

	if moveTile.To == "main" {
		toBoard = lobby.Board
		boardUpdate = true
	} else {
		toBoard = lobbyMember.Deck
	}

	// Disallow this move if it affects the board when it's not the users turn
	if boardUpdate && lobby.PlayerTurnID != &lobbyMember.UserID {
		return c.JSON(403, entity.ErrNotYourTurn)
	}

	// Make sure the tiles actually exist
	if !validatePos(fromBoard, moveTile.FromIndex) || !validatePos(toBoard, moveTile.ToIndex) {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	if (*toBoard.Squares)[moveTile.ToIndex].Tile != nil {
		return c.JSON(400, entity.ErrTileOccupied)
	}

	(*toBoard.Squares)[moveTile.ToIndex].Tile = (*fromBoard.Squares)[moveTile.FromIndex].Tile
	(*fromBoard.Squares)[moveTile.FromIndex].Tile = nil

	if boardUpdate {
		_ = publishLobbyMessage(gc.ably, lobby.ID, "moveTile", map[string]interface{}{
			"move": moveTile,
			"tile": (*toBoard.Squares)[moveTile.ToIndex].Tile,
		})
	}

	gc.db.Save(&lobby)
	gc.db.Save(&lobbyMember)

	return c.NoContent(204)
}

func validatePos(board entity.SquareSet, position int) bool {
	return position >= 0 && position < len(*board.Squares)
}

func (gc *GameController) EndTurn(c echo.Context) error {
	//lobby := c.Get("lobby").(*model.Lobby)
	//lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)

	return c.NoContent(204)
}
