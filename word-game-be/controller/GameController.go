package controller

import (
	_ "embed"
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/constant"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably-labs/word-game/word-game-be/util"
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
	g.GET("/boards/deck", gc.GetDeck)
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

func (gc *GameController) GetDeck(c echo.Context) error {
	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)
	if lobbyMember.MemberType == constant.MemberTypeSpectator {
		return c.JSON(403, entity.ErrSpectating)
	}
	return c.JSON(200, lobbyMember.Deck)
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
	if boardUpdate && *lobby.PlayerTurnID != lobbyMember.UserID {
		fmt.Println(lobby.PlayerTurnID, lobbyMember.UserID)
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
		_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "moveTile", map[string]interface{}{
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
	lobby := c.Get("lobby").(*model.Lobby)
	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)

	score := util.ValidateBoard(lobby.Board)

	if score == 0 {
		return c.JSON(400, entity.ErrInvalidPlay)
	}

	lobbyMember.Score += score

	_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "message", entity.ChatSent{
		Message: fmt.Sprintf("%s scored %d points", lobbyMember.User.Name, score),
		Author:  "system",
	})

	_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "scoreUpdate", lobbyMember.Score)

	remainingTiles := lobbyMember.Deck.TileCount()

	if len(*lobby.Bag.Squares) > 0 {
		lobbyMember.Deck.AddTiles(util.TakeFromBag(7-remainingTiles, &lobby.Bag))
	}

	fmt.Println("Deck length", len(*lobbyMember.Deck.Squares))

	indices := util.GetPlacedTileIndices(lobby.Board)
	for _, index := range indices {
		(*lobby.Board.Squares)[index].Tile.Draggable = false
	}

	gc.db.Order("joined_at").Find(&lobby.Members)

	for i, member := range lobby.Members {
		if member.UserID == *lobby.PlayerTurnID {
			if i == len(lobby.Members)-1 {
				*lobby.PlayerTurnID = 0
			} else {
				*lobby.PlayerTurnID = lobby.Members[i+1].UserID
			}
			break
		}
	}

	fmt.Println("It is now this players turn ", *lobby.PlayerTurnID)
	_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "lobbyUpdate", lobby)

	gc.db.Save(&lobbyMember)
	gc.db.Save(&lobby)

	return c.NoContent(204)
}
